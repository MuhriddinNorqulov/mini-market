package notification

import (
	"context"

	"mini-market/src/core/domain/ports/notification"
)

type NoopAlertNotifier struct{}

// @inject
func NewNoopAlertNotifier() notification.AlertNotifier {
	return &NoopAlertNotifier{}
}

func (this *NoopAlertNotifier) Send(ctx context.Context, message string) {}
