package temp

import (
	"container/list"
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Packet 表示需要发送的网络数据包
type Packet struct {
	Addr net.Addr
	Data []byte
}

// LargePacketQueue 支持大容量缓冲区的发包队列
type LargePacketQueue struct {
	mu          sync.Mutex      // 保护队列的互斥锁
	cond        *sync.Cond      // 用于等待/通知机制
	queue       *list.List      // 链表实现的动态队列
	maxSize     int64           // 队列最大容量
	currentSize int64           // 当前队列大小(原子操作)
	workers     int             // 工作协程数量
	wg          sync.WaitGroup  // 用于等待所有工作协程退出
	ctx         context.Context // 用于控制工作协程退出
	cancel      context.CancelFunc
	conn        net.Conn // 网络连接
	batchSize   int      // 批量处理大小
	closed      bool     // 队列是否已关闭
}

// NewLargePacketQueue 创建一个新的大容量发包队列
// maxSize: 队列最大容量
// workers: 工作协程数量
// conn: 网络连接
func NewLargePacketQueue(maxSize int64, workers int, conn net.Conn) *LargePacketQueue {
	ctx, cancel := context.WithCancel(context.Background())
	q := &LargePacketQueue{
		queue:     list.New(),
		maxSize:   maxSize,
		workers:   workers,
		ctx:       ctx,
		cancel:    cancel,
		conn:      conn,
		batchSize: 64, // 批量大小可根据实际情况调整
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Start 启动工作协程
func (q *LargePacketQueue) Start() {
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}
}

// worker 工作协程，处理队列中的数据包
func (q *LargePacketQueue) worker(id int) {
	defer q.wg.Done()

	batch := make([]*Packet, 0, q.batchSize)

	for {
		// 获取一批数据包
		pkts := q.dequeueBatch(q.batchSize)
		if len(pkts) > 0 {
			batch = append(batch, pkts...)

			// 当达到批量大小或队列暂时为空时处理
			if len(batch) >= q.batchSize || q.queue.Len() == 0 {
				q.processBatch(batch)
				batch = batch[:0]
			}
		}

		// 检查是否需要退出
		select {
		case <-q.ctx.Done():
			// 处理剩余数据包
			if len(batch) > 0 {
				q.processBatch(batch)
			}
			// 处理队列中可能剩余的数据包
			for {
				pkts := q.dequeueBatch(q.batchSize)
				if len(pkts) == 0 {
					break
				}
				q.processBatch(pkts)
			}
			return
		default:
			// 继续循环处理
		}
	}
}

// dequeueBatch 从队列中取出一批数据包
func (q *LargePacketQueue) dequeueBatch(max int) []*Packet {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 如果队列为空，等待新元素
	for q.queue.Len() == 0 && !q.closed {
		q.cond.Wait()
	}

	if q.closed || q.queue.Len() == 0 {
		return nil
	}

	// 取出最多max个元素
	count := 0
	result := make([]*Packet, 0, max)
	for e := q.queue.Front(); e != nil && count < max; e = e.Next() {
		result = append(result, e.Value.(*Packet))
		q.queue.Remove(e)
		count++
	}

	// 更新当前大小
	atomic.AddInt64(&q.currentSize, -int64(count))

	// 通知可能等待的生产者
	q.cond.Broadcast()

	return result
}

// processBatch 处理一批数据包
func (q *LargePacketQueue) processBatch(batch []*Packet) {
	if len(batch) == 0 {
		return
	}

	// 设置写入超时
	if err := q.conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond)); err != nil {
		// 处理超时设置错误
		return
	}

	// 发送数据
	for _, pkt := range batch {
		if _, err := q.conn.WriteTo(pkt.Data, pkt.Addr); err != nil {
			// 处理发送错误，可以根据需要添加重试逻辑
		}
	}
}

// Enqueue 添加数据包到队列
// 非阻塞添加，当队列满时返回错误
func (q *LargePacketQueue) Enqueue(pkt *Packet) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return fmt.Errorf("queue is closed")
	}

	// 检查队列是否已满
	if atomic.LoadInt64(&q.currentSize) >= q.maxSize {
		return fmt.Errorf("queue is full")
	}

	// 添加到队列
	q.queue.PushBack(pkt)
	atomic.AddInt64(&q.currentSize, 1)

	// 通知等待的消费者
	q.cond.Signal()
	return nil
}

// EnqueueBlocking 阻塞添加数据包到队列，直到成功或被取消
func (q *LargePacketQueue) EnqueueBlocking(pkt *Packet, timeout time.Duration) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 如果队列已关闭，直接返回
	if q.closed {
		return fmt.Errorf("queue is closed")
	}

	// 等待队列有空间，同时监听上下文取消信号
	expiration := time.Now().Add(timeout)
	for atomic.LoadInt64(&q.currentSize) >= q.maxSize && !q.closed {
		// 计算等待时间
		waitTime := expiration.Sub(time.Now())
		if waitTime <= 0 {
			return fmt.Errorf("enqueue timeout")
		}

		// 等待指定时间或直到被唤醒
		ch := make(chan struct{})
		go func() {
			q.cond.Wait()
			close(ch)
		}()

		select {
		case <-ch:
			// 被唤醒，继续检查队列空间
		case <-time.After(waitTime):
			// 超时
			return fmt.Errorf("enqueue timeout")
		case <-q.ctx.Done():
			// 上下文被取消
			return q.ctx.Err()
		}
	}

	if q.closed {
		return fmt.Errorf("queue is closed")
	}

	// 添加到队列
	q.queue.PushBack(pkt)
	atomic.AddInt64(&q.currentSize, 1)

	// 通知等待的消费者
	q.cond.Signal()
	return nil
}

// Close 关闭队列，等待所有工作协程退出
func (q *LargePacketQueue) Close() {
	q.mu.Lock()
	q.closed = true
	q.mu.Unlock()

	q.cond.Broadcast() // 唤醒所有等待的goroutine
	q.cancel()         // 通知工作协程退出
	q.wg.Wait()        // 等待所有工作协程完成
}

// Status 返回队列当前状态
func (q *LargePacketQueue) Status() (current, max int64) {
	return atomic.LoadInt64(&q.currentSize), q.maxSize
}
