package gateway

import (
	"encoding/json"
	"errors"
	"github.com/duke-git/lancet/v2/cryptor"
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
	maxMessageSize = 32768
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 跨域
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Payload struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data,omitempty"`
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
	// 唯一ID
	uniqueId string
}

func NewChatClient() *ChatClient {
	return &ChatClient{
		send:   make(chan []byte, 256),
		closed: true,
	}
}

// Conn 创建新连接
func (r *ChatClient) Conn(ctx *gin.Context) error {
	// 关闭旧连接
	if r.conn != nil {
		_ = r.conn.Close()
	}

	// 生成唯一ID
	s := ctx.Request.RemoteAddr + ctx.Request.Header.Get("User-Agent") + strconv.FormatInt(time.Now().Unix(), 10)

	uniqueId := cryptor.Md5String(s)
	if uniqueId == "" {
		ctx.JSON(http.StatusOK, gin.H{})
		return errors.New("failed to generate unique ID")
	}

	// 升级为WebSocket协议
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return err
	}

	r.conn = conn
	r.uniqueId = uniqueId
	r.closed = false

	// 注册客户端
	Register(uniqueId, r)
	// 在新协程上执行读写操作
	go r.Read()
	go r.Write()

	// 5秒后执行
	time.AfterFunc(5*time.Second, func() {
		_ = r.SendMsg("test", "12312323")
	})

	return nil
}

// Read 读消息
func (r *ChatClient) Read() {
	defer r.Close()

	// 读 消息大小限制
	r.conn.SetReadLimit(maxMessageSize)
	// 读 超时时间
	err := r.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		// 记录错误
		logrus.Errorf("SetReadDeadline error: %v\n", err)
		return
	}
	// 读 超时操作
	r.conn.SetPongHandler(func(string) error { _ = r.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var payload Payload
		// 读取json
		err := r.conn.ReadJSON(&payload)
		if err != nil {
			// 是否json解析错误
			var unmarshalTypeError *json.UnmarshalTypeError
			if errors.As(err, &unmarshalTypeError) {
				logrus.Warnf("json解析错误: %+v", err)
				_ = r.SendMsg("error", map[string]interface{}{"message": "json解析错误", "code": 400})
				continue
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// 记录错误
				logrus.Errorf("Read message error: %v\n", err)
			}
			break
		}
		logrus.Debugf("接收到的内容: %+v", payload)
		// 调用钩子
		CallReadHook(payload.Event, r, payload)
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
func (r *ChatClient) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		r.Close()
	}()

	for {
		select {
		case message, ok := <-r.send:
			if r.closed {
				logrus.Warnf("客户端[%s]已关闭", r.uniqueId)
				return
			}
			err := r.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logrus.Errorf("SetWriteDeadline error: %v", err)
				return
			}
			if !ok {
				// The hub closed the channel.
				_ = r.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// 写消息
			err = r.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logrus.Errorf("WriteMessage error: %v", err)
				return
			}
		case <-ticker.C:
			_ = r.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := r.conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				return
			}
		}
	}
}

func (r *ChatClient) Send(msg []byte) {
	r.send <- msg
}

func (r *ChatClient) SendMsg(event string, data interface{}) error {
	payload := Payload{
		Event: event,
		Data:  data,
	}
	// 解析为JSON
	j, err := json.Marshal(&payload)
	if err != nil {
		return err
	}
	r.send <- j
	return nil
}

func (r *ChatClient) Close() {
	logrus.Infof("正在关闭[%s]客户端,当前连接池中有%d个连接", r.uniqueId, len(pool.clients))

	// 从连接池中释放客户端
	Unregister(r.uniqueId)
	logrus.Infof("已从连接池中释放[%s]客户端,连接池中还有%d个连接", r.uniqueId, len(pool.clients))

	// 关闭连接
	if !r.closed {
		r.closed = true
		err := r.conn.Close()
		if err != nil {
			logrus.Errorf("Close WebSocket error: %v", err)
			return
		}
		logrus.Infof("已关闭[%s]客户端", r.uniqueId)
	}
}

func (r *ChatClient) GetUniqueId() string {
	return r.uniqueId
}
