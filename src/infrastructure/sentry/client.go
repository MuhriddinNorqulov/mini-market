package sentry

import (
	"context"
	"fmt"
	"time"

	"mini-market/src/infrastructure/env"
	"mini-market/src/infrastructure/telemetry"

	"github.com/getsentry/sentry-go"
)

type Client struct {
	Enabled bool

	sendBodies   bool
	flushTimeout time.Duration
	sentry       *sentry.Client
}

// @inject
func NewClient(cfg *env.Env, role telemetry.Role) *Client {
	if !cfg.SentryEnabled || cfg.SentryDSN == "" {
		return &Client{Enabled: false}
	}

	sc, err := sentry.NewClient(sentry.ClientOptions{
		Dsn: cfg.SentryDSN,

		Environment: cfg.OtelEnvironment,
		Release:     cfg.OtelServiceVersion,
		ServerName:  role.ServiceName(cfg),

		SampleRate: 1.0,

		EnableTracing: false,

		AttachStacktrace: false,
		SendDefaultPII:   false,
	})
	if err != nil {

		panic(fmt.Sprintf("[sentry] client: %v", err))
	}

	return NewClientFromSentry(sc, cfg.SentrySendBodies, time.Duration(cfg.SentryFlushSeconds)*time.Second)
}

func Disabled() *Client {
	return &Client{Enabled: false}
}

func NewClientFromSentry(sc *sentry.Client, sendBodies bool, flushTimeout time.Duration) *Client {
	return &Client{
		Enabled:      sc != nil,
		sendBodies:   sendBodies,
		flushTimeout: flushTimeout,
		sentry:       sc,
	}
}

func (this *Client) Shutdown(ctx context.Context) error {
	if !this.Enabled || this.sentry == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, this.flushTimeout)
	defer cancel()

	if ok := this.sentry.FlushWithContext(ctx); !ok {
		return fmt.Errorf("sentry: flush did not complete within %s", this.flushTimeout)
	}
	return nil
}

func (this *Client) capture(ev *sentry.Event) {

	if !this.Enabled || this.sentry == nil {
		return
	}
	this.sentry.CaptureEvent(ev, nil, sentry.NewScope())
}
