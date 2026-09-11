package telemetry

import "mini-market/src/infrastructure/env"

type Role string

const (
	RoleHTTP  Role = "http"
	RoleAsync Role = "async"
)

func (this Role) ServiceName(cfg *env.Env) string {
	if this == RoleAsync {
		return cfg.OtelServiceNameAsync
	}
	return cfg.OtelServiceNameHTTP
}
