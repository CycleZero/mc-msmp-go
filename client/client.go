package client

import (
	"encoding/json"
	"fmt"
	"github.com/CycleZero/mc-msmp-go/dto"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

// MsmpClient WebSocket客户端结构
type MsmpClient struct {
	// WebSocket连接
	conn *websocket.Conn

	// 服务器地址
	url string

	// 连接状态
	connected bool

	// 互斥锁，保护连接状态
	mutex sync.Mutex

	// 消息处理函数
	messageHandler func(dto.MsmpResponse)

	// 连接关闭处理函数
	closeHandler func(int, string) error

	// 是否启用自动重连
	autoReconnect bool

	// 重连间隔
	reconnectInterval time.Duration

	// 请求ID计数器
	requestID int

	// 等待响应的请求映射
	pendingRequests map[int]chan dto.MsmpResponse

	// 退出信号
	done chan struct{}
}

// NewMsmpClient 创建新的MsmpWebSocket客户端实例
func NewMsmpClient(url string) *MsmpClient {
	return &MsmpClient{
		url:               url,
		connected:         false,
		autoReconnect:     true,
		reconnectInterval: 5 * time.Second,
		requestID:         0,
		pendingRequests:   make(map[int]chan dto.MsmpResponse),
		done:              make(chan struct{}),
	}
}

// SetMessageHandler 设置消息处理函数
func (c *MsmpClient) SetMessageHandler(handler func(dto.MsmpResponse)) {
	c.messageHandler = handler
}

// SetCloseHandler 设置连接关闭处理函数
func (c *MsmpClient) SetCloseHandler(handler func(int, string) error) {
	c.closeHandler = handler
}

// SetAutoReconnect 设置自动重连
func (c *MsmpClient) SetAutoReconnect(autoReconnect bool) {
	c.autoReconnect = autoReconnect
}

// SetReconnectInterval 设置重连间隔
func (c *MsmpClient) SetReconnectInterval(interval time.Duration) {
	c.reconnectInterval = interval
}

// Connect 连接到WebSocket服务器
func (c *MsmpClient) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.connected {
		return fmt.Errorf("client already connected")
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	c.conn = conn
	c.connected = true

	// 启动读取消息的goroutine
	go c.readMessages()

	// 启动自动重连的goroutine（如果启用）
	if c.autoReconnect {
		go c.reconnect()
	}

	log.Printf("Connected to %s", c.url)
	return nil
}

// Disconnect 断开WebSocket连接
func (c *MsmpClient) Disconnect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.connected {
		return fmt.Errorf("client not connected")
	}

	close(c.done)
	c.connected = false
	return c.conn.Close()
}

// readMessages 读取来自服务器的消息
func (c *MsmpClient) readMessages() {
	for {
		select {
		case <-c.done:
			return
		default:
			// 设置读取超时
			c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))

			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Error reading message: %v", err)
				}
				// 连接断开，触发重连逻辑
				c.mutex.Lock()
				c.connected = false
				c.mutex.Unlock()
				return
			}

			// 解析响应
			response, err := dto.ParseResponse(message)
			if err != nil {
				log.Printf("Error parsing response: %v", err)
				continue
			}

			// 检查是否有等待此响应的请求
			c.mutex.Lock()
			if ch, exists := c.pendingRequests[response.GetID()]; exists {
				// 发送到等待的通道
				ch <- response
				// 删除已处理的请求
				delete(c.pendingRequests, response.GetID())
			} else {
				// 调用全局消息处理函数
				if c.messageHandler != nil {
					c.messageHandler(response)
				}
			}
			c.mutex.Unlock()
		}
	}
}

// reconnect 自动重连逻辑
func (c *MsmpClient) reconnect() {
	for {
		select {
		case <-c.done:
			return
		default:
			c.mutex.Lock()
			needsReconnect := !c.connected && c.autoReconnect
			c.mutex.Unlock()

			if needsReconnect {
				log.Printf("Attempting to reconnect to %s", c.url)
				err := c.Connect()
				if err != nil {
					log.Printf("Reconnect failed: %v", err)
					time.Sleep(c.reconnectInterval)
					continue
				}
				log.Printf("Reconnected successfully to %s", c.url)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// SendRequest 发送请求并等待响应
func (c *MsmpClient) SendRequest(method string, params interface{}) (dto.MsmpResponse, error) {
	c.mutex.Lock()
	if !c.connected {
		c.mutex.Unlock()
		return nil, fmt.Errorf("not connected to server")
	}

	// 增加请求ID
	c.requestID++
	id := c.requestID

	// 创建等待响应的通道
	responseChan := make(chan dto.MsmpResponse, 1)
	c.pendingRequests[id] = responseChan

	// 构造请求
	request := dto.NewMsmpRequest(id, method, params)
	c.mutex.Unlock()

	// 发送请求
	data, err := json.Marshal(request)
	if err != nil {
		c.mutex.Lock()
		delete(c.pendingRequests, id)
		c.mutex.Unlock()
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	err = c.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		c.mutex.Lock()
		delete(c.pendingRequests, id)
		c.mutex.Unlock()
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	// 等待响应或超时
	select {
	case response := <-responseChan:
		return response, nil
	case <-time.After(30 * time.Second):
		c.mutex.Lock()
		delete(c.pendingRequests, id)
		c.mutex.Unlock()
		return nil, fmt.Errorf("request timeout")
	case <-c.done:
		return nil, fmt.Errorf("client disconnected")
	}
}

// SendNotification 发送通知（不需要响应）
func (c *MsmpClient) SendNotification(method string, params interface{}) error {
	c.mutex.Lock()
	if !c.connected {
		c.mutex.Unlock()
		return fmt.Errorf("not connected to server")
	}

	// 通知的ID为0
	request := dto.NewMsmpRequest(0, method, params)
	c.mutex.Unlock()

	// 发送请求
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %v", err)
	}

	err = c.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return fmt.Errorf("failed to send notification: %v", err)
	}

	return nil
}

// IsConnected 检查是否已连接
func (c *MsmpClient) IsConnected() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.connected
}
