package hooks

import "VitaTaskGo/internal/pkg/gateway"

func init() {
	gateway.RegisterReadHook("auth", "AuthUser", AuthUser)
	gateway.RegisterReadHook("ping", "DefaultPing", PingHandle)
}
