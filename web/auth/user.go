package auth

import "github.com/danielgtaylor/huma/v2"

type ResolvedUser struct {
	ProviderName    string
	ProviderId      string
	Username        string
	Firstname       string
	Lastname        string
	Email           string
	IsEmailVerified bool
}

func (r *ResolvedUser) Resolve(ctx huma.Context) []error {
	user := userFromContext(ctx.Context())
	if user == nil {
		return []error{
			huma.Error401Unauthorized("Invalid credentials"),
		}
	}
	*r = *user
	return nil
}
