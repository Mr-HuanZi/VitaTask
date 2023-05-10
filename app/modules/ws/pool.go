package ws

import (
	"github.com/duke-git/lancet/v2/cryptor"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
	// 初始化线程池
	pool = NewPool()
)

// Client 一个接口，我也不知道能用来干嘛，先定义着吧
type Client interface {
	Read()
	Write()
	Send([]byte)
}

// Pool 协程池
type Pool struct {
	// clients. 已注册的客户端
	clients map[uint64]*ChatClient
	// 读写锁
	rw sync.RWMutex
}

func NewPool() *Pool {
	return &Pool{
		clients: make(map[uint64]*ChatClient),
	}
}

func StandardPool() *Pool {
	return pool
}

// Register 注册客户端
func (p *Pool) Register(c *ChatClient) {
	p.rw.Lock()
	defer p.rw.Unlock()
	p.clients[c.userId] = c
}

// Unregister 注销 客户端
func (p *Pool) Unregister(uid uint64) {
	p.rw.Lock()
	defer p.rw.Unlock()
	if _, ok := p.clients[uid]; ok {
		delete(p.clients, uid)
	}
}

// Broadcast 广播消息
func (p *Pool) Broadcast(c *ChatClient, messageType int, message []byte) {
	p.rw.RLock()
	defer p.rw.RUnlock()
	// 遍历所有客户端
	for _, client := range p.clients {
		// 排除自己
		if client != c {
			// 向通道发送消息
			client.send <- message
		}
	}
}

func Unregister(uid uint64) {
	pool.Unregister(uid)
}

func Register(c *ChatClient) {
	pool.Register(c)
}

func Broadcast(c *ChatClient, messageType int, message []byte) {
	pool.Broadcast(c, messageType, message)
}

func GenerateToken(payload []string) string {
	// 将当前时间戳加入到 payload
	payload = append(payload, strconv.FormatInt(time.Now().UnixMilli(), 10))
	// 将切片转换为字符串并哈希
	return cryptor.Md5String(strings.Join(payload, ""))
}

func GetClient(uid uint64) *ChatClient {
	client, ok := pool.clients[uid]
	if ok {
		return client
	}
	return nil
}

func Send(to *ChatClient, module string, data interface{}) error {
	return to.Send(module, data)
}

func Online(uid uint64) bool {
	return GetClient(uid) != nil
}

func GetClients() map[uint64]*ChatClient {
	return pool.clients
}
