package backend

import (
	"context"
	"sync"

	"github.com/behnh/vault-plugin-launchdarkly/internal/client"
	internal_config "github.com/behnh/vault-plugin-launchdarkly/internal/config"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type LaunchDarklyBackend struct {
	*framework.Backend
	view      logical.Storage
	client    client.Client
	lock      sync.RWMutex
	roleLocks []*locksutil.LockEntry
}

func (b *LaunchDarklyBackend) GetClient(ctx context.Context, s logical.Storage) (client.Client, error) {
	b.lock.RLock()
	unlockFunc := b.lock.Unlock
	defer func() { unlockFunc() }()

	if b.client != nil {
		return b.client, nil
	}

	b.lock.RUnlock()
	b.lock.Lock()
	unlockFunc = b.lock.Unlock

	if b.client != nil {
		return b.client, nil
	}

	config, err := internal_config.GetConfig(ctx, s)
	if err != nil {
		return nil, err
	}

	c, err := client.NewClient(config)
	if err != nil {
		return nil, err
	}
	b.client = c

	return c, nil
}

func (b *LaunchDarklyBackend) Reset() {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.client = nil
}

func (b *LaunchDarklyBackend) Invalidate(ctx context.Context, key string) {
	switch key {
	case internal_config.PathPatternConfig:
		b.Reset()
	}
}

func Init(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(c)
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend(conf *logical.BackendConfig) *LaunchDarklyBackend {
	backend := &LaunchDarklyBackend{
		view:      conf.StorageView,
		roleLocks: locksutil.CreateLocks(),
	}

	backend.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help:        backendHelp,
		Paths:       framework.PathAppend(),
		Invalidate:  backend.Invalidate,
	}

	return backend
}

const backendHelp = `
The LaunchDarkly token engine dynamically generates LaunchDarkly service tokens, based on user-defined inline
permissions, or custom roles.

After mounting this secrets engine, you can configure the initial credentials using the "config/" endpoints. You can
generate service tokens using the "token/" endpoints.`
