package config

type ConfigProvider interface {
	GetAccessTokenExpireMinutes() int64
	GetRefreshTokenExpireMinutes() int64
}
