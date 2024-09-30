package auth

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type ResolvedUser struct {
	ProviderName    string
	ProviderId      string
	Username        string
	Firstname       string
	Lastname        string
	Email           string
	IsEmailVerified bool
	Avatar          string
	IsAdmin         bool
}

func (r *ResolvedUser) ProviderUUID() uuid.UUID {
	return uuid.MustParse(r.ProviderId)
}

func (r *ResolvedUser) Resolve(ctx huma.Context) []error {
	user := UserFromContext(ctx.Context())
	if user == nil {
		return []error{
			huma.Error401Unauthorized("Invalid credentials"),
		}
	}
	*r = *user
	return nil
}
