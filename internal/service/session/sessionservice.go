package session

import (
	"context"
	"goth/internal/store"
)

type Session struct {
	IsDemo             bool
	CurrentUser        *store.GoogleUser
	CurrentWorkspace   *any
	AllWorkspaces      []*any
	StripeSubscription *any
}

func (s *Session) GetUserName() string {
	if s != nil && s.CurrentUser != nil {
		return s.CurrentUser.Name
	}
	return "Demo user"
}

func (s *Session) GetProfilePictureUrl() string {
	if s != nil && s.CurrentUser != nil {
		return s.CurrentUser.Picture
	}
	return "https://i.pinimg.com/564x/ab/8f/d1/ab8fd1cd803e62a5120708944874ee49.jpg"
}

func GetSessionFromCtx(ctx context.Context) *Session {
	var value = ctx.Value("session")
	if value == nil {
		return nil
	}
	return value.(*Session)
}

func GetUserNameFromCtx(ctx context.Context) string {
	var session = GetSessionFromCtx(ctx)
	if session != nil && session.CurrentUser != nil {
		return session.CurrentUser.Name
	}
	return "Demo user"
}

func GetProfilePicFromCtx(ctx context.Context) string {
	var session = GetSessionFromCtx(ctx)
	if session != nil && session.CurrentUser != nil {
		return session.CurrentUser.Picture
	}
	return "Demo user"
}
