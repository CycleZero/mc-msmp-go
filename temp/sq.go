package temp

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// Packet 表示需要发送的网络数据包

// PacketQueue 高性能发包队列
type PacketQueue struct {
	queue     chan *Packet    // 带缓冲的数据包队列
	workers   int             // 工作协程数量
	wg        sync.WaitGroup  // 用于等待所有工作协程退出
	ctx       context.Context // 用于控制工作协程退出
	cancel    context.CancelFunc
	conn      net.Conn // 网络连接
	batchSize int      // 批量处理大小
}

// NewPacketQueue 创建一个新的发包队列
// bufferSize: 队列缓冲区大小
// workers: 工作协程数量
// conn: 网络连接
func NewPacketQueue(bufferSize, workers int, conn net.Conn) *PacketQueue {
	ctx, cancel := context.WithCancel(context.Background())

	return &PacketQueue{
		queue:     make(chan *Packet, bufferSize),
		workers:   workers,
		ctx:       ctx,
		cancel:    cancel,
		conn:      conn,
		batchSize: 32, // 可根据实际情况调整
	}
}

// Start 启动工作协程
func (q *PacketQueue) Start() {
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}
}

// worker 工作协程，处理队列中的数据包
func (q *PacketQueue) worker(id int) {
	defer q.wg.Done()

	// 批量处理的缓冲区
	batch := make([]*Packet, 0, q.batchSize)

	for {
		select {
		case <-q.ctx.Done():
			// 退出前处理剩余的数据包
			q.processBatch(batch)
			return

		case pkt, ok := <-q.queue:
			if !ok {
				return
			}

			// 添加到批量处理缓冲区
			batch = append(batch, pkt)

			// 当达到批量大小或队列暂时为空时处理
			if len(batch) >= q.batchSize || len(q.queue) == 0 {
				q.processBatch(batch)
				batch = batch[:0] // 重置批量缓冲区
			}
		}
	}
}

// processBatch 处理一批数据包
func (q *PacketQueue) processBatch(batch []*Packet) {
	if len(batch) == 0 {
		return
	}

	// 这里可以根据需要优化发送策略
	for _, pkt := range batch {
		// 设置写入超时，避免长时间阻塞
		if err := q.conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond)); err != nil {
			// 处理超时设置错误
			continue
		}

		// 发送数据
		if _, err := q.conn.WriteTo(pkt.Data, pkt.Addr); err != nil {
			// 处理发送错误
		}
	}
}

// Enqueue 添加数据包到队列
// 非阻塞添加，当队列满时返回错误
func (q *PacketQueue) Enqueue(pkt *Packet) error {
	select {
	case q.queue <- pkt:
		return nil
	default:
		// 队列已满，可以选择阻塞等待或返回错误
		// 这里选择返回错误，由调用者决定如何处理
		return fmt.Errorf("queue is full")
	}
}

// EnqueueBlocking 阻塞添加数据包到队列
func (q *PacketQueue) EnqueueBlocking(pkt *Packet) error {
	select {
	case q.queue <- pkt:
		return nil
	case <-q.ctx.Done():
		return q.ctx.Err()
	}
}

// Close 关闭队列，等待所有工作协程退出
func (q *PacketQueue) Close() {
	q.cancel()     // 通知工作协程退出
	q.wg.Wait()    // 等待所有工作协程完成
	close(q.queue) // 关闭队列
}

// 队列状态信息
func (q *PacketQueue) Status() (int, int) {
	return len(q.queue), cap(q.queue)
}
