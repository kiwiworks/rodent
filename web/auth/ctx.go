package auth

import "context"

type resolvedUserKey struct{}

func InjectUser(ctx context.Context, user *ResolvedUser) context.Context {
	return context.WithValue(ctx, resolvedUserKey{}, user)
}

func UserFromContext(ctx context.Context) *ResolvedUser {
	if user, ok := ctx.Value(resolvedUserKey{}).(*ResolvedUser); ok {
		return user
	}
	return nil
}
