package diagnostics

import "context"

type Diagnostics interface {
	DBPing(ctx context.Context) error

	SelfCall(ctx context.Context) error
}
