package cas

import (
	"context"
	"fmt"
	"net/url"

	casdoor "github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"rodent/system/errors"
)

type CasdoorClient struct {
	client *casdoor.Client
	org    string
	app    string
}

type CasdoorClientConfig struct {
	Endpoint         url.URL
	ClientId         string
	ClientSecret     string
	Certificate      string
	OrganizationName string
	ApplicationName  string
}

func NewCasdoorClient(cfg *CasdoorClientConfig) (*CasdoorClient, error) {
	client := casdoor.NewClient(
		cfg.Endpoint.String(), cfg.ClientId, cfg.ClientSecret,
		cfg.Certificate,
		cfg.OrganizationName, cfg.ApplicationName,
	)
	return &CasdoorClient{
		client: client,
		org:    cfg.OrganizationName,
		app:    cfg.ApplicationName,
	}, nil
}

func (c *CasdoorClient) Domain() Domain {
	return Domain{
		Organisation: c.org,
		Application:  c.app,
	}
}

func (c *CasdoorClient) ApplicationDetails(_ context.Context, domain Domain) (*casdoor.Application, error) {
	if domain.Organisation != c.org {
		return nil, errors.Newf("domain org mismatch '%s' whereas client is valid for '%s'", domain.Organisation, c.org)
	}
	if domain.Application != c.app {
		return nil, errors.Newf("domain app mismatch '%s' whereas client is valid for '%s'", domain.Application, c.app)
	}

	app, err := c.client.GetApplication(domain.Application)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get application '%s' details", domain.Application)
	}
	if app == nil {
		return nil, errors.Wrapf(err, "could not get application '%s' details", domain.Application)
	}
	return app, nil
}

func (c *CasdoorClient) RedirectUrlForDomain(_ context.Context, domain Domain) string {
	return fmt.Sprintf("http://localhost:9000/api/v1/oidc/%s/callback", domain.Organisation)
}

type AuthFlows struct {
	SignIn string
	SignUp string
}

func (c *CasdoorClient) Flows(ctx context.Context, domain Domain, redirectUrl string) (*AuthFlows, error) {
	app, err := c.ApplicationDetails(ctx, domain)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get application '%s' details", domain.Application)
	}

	if redirectUrl == "" {
		redirectUrl = c.RedirectUrlForDomain(ctx, domain)
	}

	var flows AuthFlows
	flows.SignIn = c.client.GetSigninUrl(redirectUrl)
	if app.EnableSignUp {
		flows.SignUp = c.client.GetSignupUrl(true, redirectUrl)
	}

	return &flows, nil
}

func (c *CasdoorClient) ExchangeCode(_ context.Context, code, state string) (*oauth2.Token, error) {
	token, err := c.client.GetOAuthToken(code, state)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get OAuth token")
	}
	return token, nil
}

func (c *CasdoorClient) Refresh(_ context.Context, refreshToken string) (*oauth2.Token, error) {
	token, err := c.client.RefreshOAuthToken(refreshToken)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to refresh OAuth token")
	}
	return token, nil
}

func (c *CasdoorClient) ParseToken(token string) (*casdoor.Claims, error) {
	claims, err := c.client.ParseJwtToken(token)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse JWT token")
	}
	return claims, nil
}

func (c *CasdoorClient) GetUser(name string) (*casdoor.User, error) {
	user, err := c.client.GetUser(name)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get user %s", name)
	}
	return user, nil
}

func (c *CasdoorClient) ResolveUserFromToken(token string) (*ResolvedUser, error) {
	claims, err := c.ParseToken(token)
	if err != nil {
		return nil, err
	}
	user, err := c.GetUser(claims.Name)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.Wrapf(err, "could not get user %s", claims.Name)
	}

	return &ResolvedUser{
		CasId:           uuid.MustParse(claims.Id),
		Email:           claims.Email,
		IsEmailVerified: claims.EmailVerified,
		Name:            user.Name,
		Firstname:       user.FirstName,
		Lastname:        user.LastName,
		Avatar:          user.Avatar,
		IsAdmin:         user.IsAdmin,
	}, nil
}
