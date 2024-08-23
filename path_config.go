package launchdarkly

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *Backend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: prefix,
		},
		Fields: map[string]*framework.FieldSchema{
			"api_key": {
				Type:        framework.TypeString,
				Description: "LaunchDarkly API key",
				Required:    true,
				DisplayAttrs: &framework.DisplayAttributes{
					Sensitive: true,
				},
			},
			"base_uri": {
				Type:        framework.TypeString,
				Description: "LaunchDarkly base URI",
				Required:    false,
				Default:     "https://app.launchdarkly.com",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
			},
		},
		HelpSynopsis:    pathConfigHelpSynopsis,
		HelpDescription: pathConfigHelpDescription,
	}
}

type config struct {
	ApiKey  string `json:"api_key"`
	BaseUri string `json:"base_uri"`
}

func (b *Backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	apiKey := data.Get("api_key").(string)
	if apiKey == "" {
		return nil, errors.New("api_key is required")
	}

	baseUri := data.Get("base_uri").(string)
	if baseUri == "" {
		baseUri = "https://app.launchdarkly.com"
	}

	entry, err := logical.StorageEntryJSON("config", config{
		ApiKey:  apiKey,
		BaseUri: baseUri,
	})

	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// Reset the client so it can be recreated with the new config
	b.client.ld = nil
	return nil, nil
}

func (b *Backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get(ctx, "config")
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var c config
	if err := entry.DecodeJSON(&c); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"base_uri": c.BaseUri,
		},
	}, nil
}

const pathConfigHelpSynopsis = `
Configure the API key used by the LaunchDarkly plugin to manage API keys and read SDK keys
`

const pathConfigHelpDescription = `
Before being able to use the LaunchDarkly plugin, an API key with access to
manage API keys and read SDK keys must be configured. This path is used to
configure that API key.
`
