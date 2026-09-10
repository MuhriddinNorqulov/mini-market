package telemetry

import (
	"context"
	"errors"
	"fmt"

	"mini-market/src/infrastructure/env"

	runtimemetrics "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
	lognoop "go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Telemetry struct {
	Enabled bool

	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	loggerProvider *sdklog.LoggerProvider
}

// @inject
func NewTelemetry(cfg *env.Env, role Role) *Telemetry {
	if !cfg.OtelEnabled {
		return &Telemetry{Enabled: false}
	}

	ctx := context.Background()
	res := newResource(cfg, role)

	tp, err := newTracerProvider(ctx, cfg, res)
	if err != nil {
		panic(fmt.Sprintf("[OTel] tracer provider: %v", err))
	}
	mp, err := newMeterProvider(ctx, cfg, res)
	if err != nil {
		panic(fmt.Sprintf("[OTel] meter provider: %v", err))
	}
	lp, err := newLoggerProvider(ctx, cfg, res)
	if err != nil {
		panic(fmt.Sprintf("[OTel] logger provider: %v", err))
	}

	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	if err := runtimemetrics.Start(runtimemetrics.WithMeterProvider(mp)); err != nil {
		panic(fmt.Sprintf("[OTel] runtime metrics: %v", err))
	}

	return NewFromProviders(tp, mp, lp)
}

func NewFromProviders(
	tp *sdktrace.TracerProvider,
	mp *sdkmetric.MeterProvider,
	lp *sdklog.LoggerProvider,
) *Telemetry {
	return &Telemetry{Enabled: true, tracerProvider: tp, meterProvider: mp, loggerProvider: lp}
}

func (this *Telemetry) LoggerProvider() log.LoggerProvider {
	if this.loggerProvider == nil {
		return lognoop.NewLoggerProvider()
	}
	return this.loggerProvider
}

func (this *Telemetry) MeterProvider() metric.MeterProvider {
	if this.meterProvider == nil {
		return metricnoop.NewMeterProvider()
	}
	return this.meterProvider
}

func (this *Telemetry) Shutdown(ctx context.Context) error {
	if !this.Enabled {
		return nil
	}

	var errs []error
	if this.tracerProvider != nil {
		errs = append(errs, this.tracerProvider.Shutdown(ctx))
	}
	if this.meterProvider != nil {
		errs = append(errs, this.meterProvider.Shutdown(ctx))
	}
	if this.loggerProvider != nil {
		errs = append(errs, this.loggerProvider.Shutdown(ctx))
	}
	return errors.Join(errs...)
}
