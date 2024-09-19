package cas

import "context"

type resolvedUserKey struct{}

func InjectPerson(ctx context.Context, resolvedUser ResolvedUser) context.Context {
	return context.WithValue(ctx, resolvedUserKey{}, resolvedUser)
}

func GetPersonFromContext(ctx context.Context) (*ResolvedUser, bool) {
	value := ctx.Value(resolvedUserKey{})
	if value == nil {
		return nil, false
	}

	person, ok := value.(ResolvedUser)
	if !ok {
		return nil, false
	}
	return &person, true
}
