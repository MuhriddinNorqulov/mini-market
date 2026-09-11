package notification

import "context"

type AlertNotifier interface {
	Send(ctx context.Context, message string)
}
