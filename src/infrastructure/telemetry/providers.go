package telemetry

import (
	"context"
	"time"

	"mini-market/src/infrastructure/env"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func authHeaders(cfg *env.Env) map[string]string {
	if cfg.OtelIngestionKey == "" {
		return nil
	}
	return map[string]string{"authorization": cfg.OtelIngestionKey}
}

func newTracerProvider(ctx context.Context, cfg *env.Env, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.OtelEndpoint),
		otlptracegrpc.WithHeaders(authHeaders(cfg)),
	}

	if cfg.OtelInsecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	exp, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),

		sdktrace.WithSampler(sdktrace.ParentBased(
			sdktrace.TraceIDRatioBased(cfg.OtelTracesSampleRatio),
		)),
	), nil
}

func newMeterProvider(ctx context.Context, cfg *env.Env, res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.OtelEndpoint),
		otlpmetricgrpc.WithHeaders(authHeaders(cfg)),
	}
	if cfg.OtelInsecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}
	exp, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp,
			sdkmetric.WithInterval(time.Duration(cfg.OtelMetricIntervalSeconds)*time.Second),
		)),
	), nil
}

func newLoggerProvider(ctx context.Context, cfg *env.Env, res *resource.Resource) (*sdklog.LoggerProvider, error) {
	opts := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(cfg.OtelEndpoint),
		otlploggrpc.WithHeaders(authHeaders(cfg)),
	}
	if cfg.OtelInsecure {
		opts = append(opts, otlploggrpc.WithInsecure())
	}
	exp, err := otlploggrpc.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)),
	), nil
}
