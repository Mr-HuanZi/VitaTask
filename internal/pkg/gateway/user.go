package gateway

import "sync"

var (
	UserPool = NewUserPool()
)

// UserClientList 用户池 map的 Key 表示用户ID，值为 ClientList 的 Key
type UserClientList map[string]string

type UserClientPool struct {
	userPool UserClientList
	rw       sync.RWMutex
}

func NewUserPool() *UserClientPool {
	return &UserClientPool{userPool: make(UserClientList)}
}

// BingUserToClient 绑定 UserId 到 ClientID
func (receiver *UserClientPool) BingUserToClient(userId, ClientId string) {
	receiver.rw.Lock()
	defer receiver.rw.Unlock()
	receiver.userPool[userId] = ClientId
}

// Unbind 解绑
func (receiver *UserClientPool) Unbind(userId string) {
	receiver.rw.Lock()
	defer receiver.rw.Unlock()
	if _, ok := receiver.userPool[userId]; ok {
		delete(receiver.userPool, userId)
	}
}

// GetClientID 获取ClientID
func (receiver *UserClientPool) GetClientID(userId string) string {
	receiver.rw.RLock()
	defer receiver.rw.RUnlock()
	if _, ok := receiver.userPool[userId]; ok {
		return receiver.userPool[userId]
	}
	return ""
}

func BingUserToClient(userId, ClientId string) {
	UserPool.BingUserToClient(userId, ClientId)
}

func GetClientID(userId string) string {
	return UserPool.GetClientID(userId)
}
