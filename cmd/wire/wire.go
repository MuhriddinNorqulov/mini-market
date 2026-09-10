package wire

import (
	"mini-market/src/entrypoint/asynctask"
	"mini-market/src/entrypoint/http"
	"mini-market/src/infrastructure/telemetry"

	"github.com/google/wire"
)

func InitHttpApp() *http.App {
	wire.Build(ProviderSet, wire.Value(telemetry.RoleHTTP))
	return nil
}

func InitAsyncApp() *asynctask.App {
	wire.Build(ProviderSet, wire.Value(telemetry.RoleAsync))
	return nil
}
