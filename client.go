package mcmsmpgo

import (
	"encoding/json"
	"fmt"
	"github.com/CycleZero/mc-msmp-go/container"
	"github.com/CycleZero/mc-msmp-go/dto"
	"github.com/CycleZero/mc-msmp-go/handler"
	"github.com/CycleZero/mc-msmp-go/iface"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type NewClientConfig struct {
	handler       func(*dto.MsmpRequest, dto.MsmpResponse)
	container     iface.MessageContainer
	AutoReconnect bool
}

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
	container iface.MessageContainer

	// 退出信号
	done chan struct{}

	Handler    func(*dto.MsmpRequest, dto.MsmpResponse)
	AuthSecret string
}

// NewMsmpClient 创建新的MsmpWebSocket客户端实例
func NewMsmpClient(url, secret string, config *NewClientConfig) *MsmpClient {
	c := &NewClientConfig{
		handler:       handler.DefaultHandler,
		container:     container.NewMapMessageContainer(),
		AutoReconnect: true,
	}
	if config != nil {
		if config.handler != nil {
			c.handler = config.handler
		}
		if config.container != nil {
			c.container = config.container
		}
		c.AutoReconnect = config.AutoReconnect
	}

	return &MsmpClient{
		url:               url,
		connected:         false,
		autoReconnect:     c.AutoReconnect,
		reconnectInterval: 5 * time.Second,
		requestID:         0,
		container:         c.container,
		done:              make(chan struct{}),
		Handler:           c.handler,
		AuthSecret:        secret,
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
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+c.AuthSecret)
	if c.connected {
		return fmt.Errorf("client already connected")
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.url, headers)
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
				fmt.Println("Error reading message: %v", err)
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Error reading message: %v", err)
				}
				// 连接断开，触发重连逻辑
				c.mutex.Lock()
				c.connected = false
				c.mutex.Unlock()
				return
			}
			fmt.Println("Received message:", string(message))
			// 解析响应
			response, err := dto.ParseResponse(message)
			if err != nil {
				log.Printf("Error parsing response: %v", err)
				continue
			}

			// 检查是否有等待此响应的请求
			err = c.container.AddResponse(response)
			if err != nil {
				log.Printf("Error adding response: %v", err)
				continue
			}
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
func (c *MsmpClient) SendRequest(method string, params interface{}) error {
	c.mutex.Lock()
	if !c.connected {
		c.mutex.Unlock()
		return fmt.Errorf("not connected to server")
	}

	// 增加请求ID
	c.requestID++
	id := c.requestID

	// 构造请求
	request := dto.NewMsmpRequest(id, method, params)
	c.mutex.Unlock()

	// 发送请求
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}
	fmt.Println("Send Request:" + string(data))
	err = c.container.AddRequest(&request)
	if err != nil {
		return err
	}
	err = c.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		err := c.container.CancelRequest(id)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to send request: %v", err)
	}

	return nil

}

func (c *MsmpClient) SendRequestWithCallback(method string, params interface{}, callback func(*dto.MsmpRequest, dto.MsmpResponse)) error {
	c.mutex.Lock()
	if !c.connected {
		c.mutex.Unlock()
		return fmt.Errorf("not connected to server")
	}

	// 增加请求ID
	c.requestID++
	id := c.requestID

	// 构造请求
	request := dto.NewMsmpRequest(id, method, params)
	c.mutex.Unlock()

	// 发送请求
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}
	err = c.container.AddRequestWithCallback(&request, callback)
	if err != nil {
		return err
	}
	err = c.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		err := c.container.CancelRequest(id)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to send request: %v", err)
	}

	return nil

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
