package env

import "github.com/caarlos0/env/v10"

func Check() error {
	loadEnv()
	return env.Parse(&Env{})
}
