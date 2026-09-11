package env

import (
	"mini-market/src/core/domain/ports/config"
)

type ConfigAdapter struct {
	env *Env
}

// @inject
func NewConfigAdapter(env *Env) config.ConfigProvider {
	return &ConfigAdapter{env: env}
}

func (this *ConfigAdapter) GetAccessTokenExpireMinutes() int64 {
	return this.env.AccessTimeExpireMinutes
}

func (this *ConfigAdapter) GetRefreshTokenExpireMinutes() int64 {
	return this.env.RefreshTokenExpireMinutes
}
