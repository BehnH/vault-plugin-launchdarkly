package launchdarkly

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	prefix  = "launchdarkly"
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type Backend struct {
	*framework.Backend

	credentialMutex sync.RWMutex
	clientMutex     sync.RWMutex

	client client

	system logical.SystemView
	config *config
}

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := NewBackend(conf.System)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	b.config = &config{
		ApiKey:  conf.Config["api_key"],
		BaseUri: conf.Config["base_uri"],
	}

	return b, nil
}

func NewBackend(system logical.SystemView) *Backend {
	var b Backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(fmt.Sprintf(backendHelp, version, commit, date)),
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config",
			},
		},
		Paths: []*framework.Path{
			b.pathConfig(),
			b.pathProject(),
			b.pathProjectReset(),
		},
		BackendType: logical.TypeLogical,
	}

	b.system = system
	return &b
}

const backendHelp = `
LaunchDarkly plugin for Vault (version: %s, commit: %s, build date: %s)

The LaunchDarkly token engine dynamically generates LaunchDarkly service tokens, based on user-defined inline
permissions, or custom roles.

After mounting this secrets engine, you can configure the initial credentials using the "config/" endpoints. You can
generate service tokens using the "token/" endpoints.`
