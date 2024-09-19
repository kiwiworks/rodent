package cas

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type ResolvedUser struct {
	CasId           uuid.UUID
	Email           string
	IsEmailVerified bool
	Name            string
	Firstname       string
	Lastname        string
	Avatar          string
	IsAdmin         bool
}

func (r *ResolvedUser) Resolve(ctx huma.Context) []error {
	user, ok := GetPersonFromContext(ctx.Context())
	if !ok {
		return []error{
			huma.Error401Unauthorized("missing mandatory Bearer authentication"),
		}
	}
	*r = *user
	return nil
}
