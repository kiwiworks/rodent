package cas

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/logger"
)

type Middleware struct {
	client *CasdoorClient
	api    huma.API
}

type MiddlewareConfig struct {
	fx.In
	Client *CasdoorClient
	Api    huma.API
}

func NewMiddleware(cfg MiddlewareConfig) *Middleware {
	middleware := &Middleware{
		client: cfg.Client,
		api:    cfg.Api,
	}
	cfg.Api.UseMiddleware(middleware.Middleware)
	return middleware
}

func extractBearerToken(authHeader string) (string, bool) {
	if strings.Index(authHeader, "Bearer ") != 0 {
		return "", false
	}
	return authHeader[len("Bearer "):], true
}

func (m *Middleware) protected(ctx huma.Context) bool {
	return len(ctx.Operation().Security) > 0
}

func (m *Middleware) handleErr(ctx huma.Context, status int, message string, err error) {
	log := logger.FromContext(ctx.Context())
	if err := huma.WriteErr(m.api, ctx, status, message); err != nil {
		log.Error("could not write error", zap.String("message", message), zap.Error(err))
	}
}

func (m *Middleware) Middleware(ctx huma.Context, next func(ctx huma.Context)) {
	bearer, ok := extractBearerToken(ctx.Header("Authorization"))
	if !ok {
		if m.protected(ctx) {
			m.handleErr(ctx, http.StatusForbidden, "malformed bearer token", nil)
			return
		}
		next(ctx)
		return
	}
	claims, err := m.client.ParseToken(bearer)
	if err != nil {
		m.handleErr(ctx, http.StatusForbidden, "invalid bearer token", err)
		return
	}
	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		m.handleErr(ctx, http.StatusForbidden, "invalid user id", err)
		return
	}
	user, err := m.client.GetUser(claims.Name)
	if err != nil {
		m.handleErr(ctx, http.StatusForbidden, "invalid user", err)
		return
	}
	newCtx := InjectPerson(ctx.Context(), ResolvedUser{
		CasId:           userId,
		Email:           user.Email,
		IsEmailVerified: user.EmailVerified,
		Name:            user.Name,
		Firstname:       user.FirstName,
		Lastname:        user.LastName,
		Avatar:          user.Avatar,
		IsAdmin:         user.IsAdmin,
	})

	next(huma.WithContext(ctx, newCtx))
}
