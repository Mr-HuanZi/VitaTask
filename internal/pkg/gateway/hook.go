package gateway

type ReadHookFunc func(*ChatClient, Payload)

type EventHooks map[string]ReadHookFunc

var ReadHooks = make(map[string]EventHooks)

func RegisterReadHook(event, hook string, f ReadHookFunc) {
	if _, ok := ReadHooks[event]; !ok {
		ReadHooks[event] = make(EventHooks)
	}

	ReadHooks[event][hook] = f
}

func UnregisterReadHook(event, hook string) {
	if _, ok := ReadHooks[event]; ok {
		if _, ok := ReadHooks[event][hook]; ok {
			delete(ReadHooks[event], hook)
		}
	}
}

func CallReadHook(event string, c *ChatClient, p Payload) {
	if _, ok := ReadHooks[event]; ok {
		for _, hook := range ReadHooks[event] {
			hook(c, p)
		}
	}
}
