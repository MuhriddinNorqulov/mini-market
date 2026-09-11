package envelope

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type envelope struct {
	TP *string         `json:"tp"`
	D  json.RawMessage `json:"d"`
}

func Wrap(ctx context.Context, payload []byte) []byte {
	if !trace.SpanContextFromContext(ctx).IsValid() {
		return payload
	}

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	tp, ok := carrier["traceparent"]
	if !ok {
		return payload
	}

	data := payload
	if len(data) == 0 {
		data = []byte("null")
	}

	wrapped, err := json.Marshal(envelope{TP: &tp, D: json.RawMessage(data)})
	if err != nil {

		return payload
	}
	return wrapped
}

func Unwrap(ctx context.Context, payload []byte) (context.Context, []byte) {
	var env envelope
	if err := json.Unmarshal(payload, &env); err != nil || env.TP == nil || env.D == nil {
		return ctx, payload
	}

	carrier := propagation.MapCarrier{"traceparent": *env.TP}
	return otel.GetTextMapPropagator().Extract(ctx, carrier), env.D
}
