package client

import (
	"context"
	"fmt"

	"github.com/behnh/vault-plugin-launchdarkly/internal/config"
	"github.com/behnh/vault-plugin-launchdarkly/internal/token"
	ldapi "github.com/launchdarkly/api-client-go"
)

type Client interface {
	CreateServiceToken(*token.ServiceTokenStorageEntry) (ldapi.Token, error)
}

type ldClient struct {
	client *ldapi.APIClient
	ctx    context.Context
}

var _ Client = &ldClient{}

func NewClient(config *config.ConfigEntry) (Client, error) {
	if config == nil {
		return nil, fmt.Errorf("no configuration provided")
	}

	ldc := &ldClient{}

	ldcConfig := ldapi.NewConfiguration()
	ldcConfig.Host = config.BaseUrl

	auth := make(map[string]ldapi.APIKey)
	auth["ApiKey"] = ldapi.APIKey{
		Key: config.Token,
	}

	c := context.WithValue(context.Background(), ldapi.ContextAPIKey, auth)

	ldc.ctx = c
	ldc.client = ldapi.NewAPIClient(ldcConfig)

	return ldc, nil
}

func (c *ldClient) CreateServiceToken(se *token.ServiceTokenStorageEntry) (ldapi.Token, error) {
	opt := ldapi.TokenBody{
		Name:         se.Name,
		ServiceToken: true,
	}

	if se.CustomRoleKeys != nil {
		opt.CustomRoleIds = se.CustomRoleKeys
	}
	if se.InlinePolicy != nil {
		opt.InlineRole = se.InlinePolicy
	}

	t, _, err := c.client.AccessTokensApi.PostToken(c.ctx, opt)
	if err != nil {
		return ldapi.Token{}, err
	}

	return t, nil
}
