package telemetry

import (
	"mini-market/src/infrastructure/env"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

func newResource(cfg *env.Env, role Role) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(role.ServiceName(cfg)),
		semconv.ServiceVersion(cfg.OtelServiceVersion),
		semconv.DeploymentEnvironmentName(cfg.OtelEnvironment),
	)
}
