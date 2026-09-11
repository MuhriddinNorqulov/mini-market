package telemetry

import (
	"context"

	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type QueueMetrics struct {
	inspector *asynq.Inspector
	tel       *Telemetry
}

// @inject
func NewQueueMetrics(inspector *asynq.Inspector, tel *Telemetry) *QueueMetrics {
	return &QueueMetrics{inspector: inspector, tel: tel}
}

func (this *QueueMetrics) Register() error {
	if !this.tel.Enabled {
		return nil
	}

	meter := this.tel.MeterProvider().Meter("asynq")

	pending, err := meter.Int64ObservableGauge("asynq.queue.pending")
	if err != nil {
		return err
	}
	active, err := meter.Int64ObservableGauge("asynq.queue.active")
	if err != nil {
		return err
	}
	retry, err := meter.Int64ObservableGauge("asynq.queue.retry")
	if err != nil {
		return err
	}
	archived, err := meter.Int64ObservableGauge("asynq.queue.archived")
	if err != nil {
		return err
	}

	_, err = meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {

		queues, err := this.inspector.Queues()
		if err != nil {
			return err
		}

		for _, name := range queues {
			info, err := this.inspector.GetQueueInfo(name)
			if err != nil {

				continue
			}
			attrs := metric.WithAttributes(attribute.String("queue", name))
			o.ObserveInt64(pending, int64(info.Pending), attrs)
			o.ObserveInt64(active, int64(info.Active), attrs)
			o.ObserveInt64(retry, int64(info.Retry), attrs)
			o.ObserveInt64(archived, int64(info.Archived), attrs)
		}
		return nil
	}, pending, active, retry, archived)

	return err
}
