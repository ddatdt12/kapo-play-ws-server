package ws

import "context"

type ConnectionContext struct {
	Ctx          context.Context
	Cancel       context.CancelFunc
	UserCanceled bool      // approach one
	ID           string    // approach two
}

func NewConnectionContext() *ConnectionContext {
	ctx := new(ConnectionContext)
	ctx.Ctx, ctx.Cancel = context.WithCancel(context.Background())
	return ctx
}
