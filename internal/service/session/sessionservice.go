package session

import (
	"context"
)

type Session struct {
	IsDemo             bool
	CurrentUser        *any
	CurrentWorkspace   *any
	AllWorkspaces      []*any
	StripeSubscription *any
}

func GetSessionFromCtx(ctx context.Context) *Session {
	var value = ctx.Value("session")
	if value == nil {
		return nil
	}
	return value.(*Session)
}
