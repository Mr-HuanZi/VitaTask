package gateway

import (
	"github.com/duke-git/lancet/v2/cryptor"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// 初始化线程池
	pool = NewPool()
)

// Client 一个接口，我也不知道能用来干嘛，先定义着吧
type Client interface {
	Read()
	Write()
	Send([]byte)
}

type ClientList map[string]Client

// Pool 协程池
type Pool struct {
	// clients. 已注册的客户端
	clients ClientList
	// 读写锁
	rw sync.RWMutex
}

func NewPool() *Pool {
	return &Pool{
		clients: make(ClientList),
	}
}

func StandardPool() *Pool {
	return pool
}

// Register 注册客户端
func (p *Pool) Register(uniqueId string, c Client) {
	p.rw.Lock()
	defer p.rw.Unlock()
	p.clients[uniqueId] = c
}

// Unregister 注销 客户端
func (p *Pool) Unregister(uniqueId string) {
	p.rw.Lock()
	defer p.rw.Unlock()
	if _, ok := p.clients[uniqueId]; ok {
		delete(p.clients, uniqueId)
	}
}

func Unregister(uniqueId string) {
	pool.Unregister(uniqueId)
}

func Register(uniqueId string, c Client) {
	pool.Register(uniqueId, c)
}

func GenerateToken(payload []string) string {
	// 将当前时间戳加入到 payload
	payload = append(payload, strconv.FormatInt(time.Now().UnixMilli(), 10))
	// 将切片转换为字符串并哈希
	return cryptor.Md5String(strings.Join(payload, ""))
}

func GetClient(id string) Client {
	client, ok := pool.clients[id]
	if ok {
		return client
	}
	return nil
}

func Send(to Client, msg []byte) {
	to.Send(msg)
}

func Online(id string) bool {
	return GetClient(id) != nil
}

func GetClients() ClientList {
	return pool.clients
}
