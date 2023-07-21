package ws

import (
	"VitaTaskGo/internal/pkg/auth"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

const (
	// 写 超时
	writeWait = 10 * time.Second
	// 读 超时
	pongWait = 60 * time.Second
	// 等 超时。必须小于 pongWait.
	pingPeriod = (pongWait * 9) / 10
	// 消息大小限制
	maxMessageSize = 10240
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 跨域
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Payload struct {
	Module string      `json:"module"`
	Data   interface{} `json:"data,omitempty"`
}

type ChatClient struct {
	// websocket 连接
	conn *websocket.Conn
	// 出站消息的缓冲通道
	send chan []byte
	// 用户ID
	userId uint64
	// 是否关闭
	closed bool
}

// Read 读消息
func (c *ChatClient) Read() {
	defer c.Close()

	// 读 消息大小限制
	c.conn.SetReadLimit(maxMessageSize)
	// 读 超时时间
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		// 记录错误
		logrus.Errorf("SetReadDeadline error: %v\n", err)
		return
	}
	// 读 超时操作
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var payload Payload
		// 读取json
		err := c.conn.ReadJSON(&payload)
		if err != nil {
			c.closed = true
			logrus.Errorf("此处记录读消息的任何错误：%+v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// 记录错误
				logrus.Errorf("Read message error: %v\n", err)
			}
			break
		}
		logrus.Debugf("接收到的内容: %+v", payload)
		// 读取字节流
		//_, message, err := c.conn.ReadMessage()
		//if err != nil {
		//	c.closed = true
		//	logrus.Errorf("此处记录读消息的任何错误：%+v", err)
		//	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		//		// 记录错误
		//		logrus.Errorf("Read message error: %v\n", err)
		//	}
		//	break
		//}
		// 将换行符替换为空格
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//pool.Broadcast(c, websocket.TextMessage, []byte(fmt.Sprintf("已收到消息，内容为%s", message)))
	}
}

// Write 写消息
func (c *ChatClient) Write() {
	defer c.Close()

	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logrus.Errorf("SetWriteDeadline error: %v", err)
				return
			}
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logrus.Errorf("obtain writer error: %v", err)
				return
			}
			_, _ = w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				logrus.Errorf("writer close error: %v", err)
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *ChatClient) Send(module string, data interface{}) error {
	payload := Payload{
		Module: module,
		Data:   data,
	}
	// 解析为JSON
	j, err := json.Marshal(&payload)
	if err != nil {
		return err
	}
	c.send <- j
	return nil
}

func (c *ChatClient) Close() {
	// 注销客户端
	Unregister(c.userId)
	// 关闭连接
	logrus.Infof("连接池中有%d个连接", len(pool.clients))
	if !c.closed {
		err := c.conn.Close()
		if err != nil {
			logrus.Errorf("Close WebSocket error: %v", err)
			return
		}
	}
}

func ClientHandle(ctx *gin.Context) {
	// 获取当前用户
	userInfo, err := auth.CurrUser(ctx)
	if err != nil {
		logrus.Errorf("WebSocket get currUser error: %v", err)
		return
	}
	// 升级为WebSocket协议
	upgrader.Subprotocols = []string{ctx.GetHeader("Sec-WebSocket-Protocol")}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Errorf("create WebSocket error: %v", err)
		return
	}
	// 发送消息
	err = conn.WriteMessage(websocket.TextMessage, []byte("已建立连接,当前时间戳是:"+strconv.FormatInt(time.Now().UnixMilli(), 10)))
	if err != nil {
		logrus.Errorf("send hello messages error: %v", err)
		return
	}
	// 该用户是否已注册
	if Online(userInfo.ID) {
		// 已注册的客户端不执行任何操作
		return
	}
	// 创建客户端实例
	client := &ChatClient{
		conn:   conn,
		send:   make(chan []byte, 256),
		userId: userInfo.ID,
		closed: false,
	}
	// 注册客户端
	Register(client)
	// 在新协程上执行读写操作
	go client.Read()
	go client.Write()
}
