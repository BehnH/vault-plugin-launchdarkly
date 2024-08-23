package launchdarkly

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestBackend_PathConfig(t *testing.T) {
	var resp *logical.Response
	var err error

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	t.Run("test write config", func(t *testing.T) {
		b := NewBackend(config.System)
		if err := b.Setup(context.Background(), config); err != nil {
			t.Fatal(err)
		}

		// Test write op
		configData := map[string]interface{}{
			"api_key":  "api_key",
			"base_uri": "https://notreal.launchdarkly.com",
		}
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "config",
			Data:      configData,
			Storage:   config.StorageView,
		})

		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("failed to write config:. resp:%#v err:%v", resp, err)
		}
	})

	t.Run("test read config", func(t *testing.T) {
		b := NewBackend(config.System)
		if err := b.Setup(context.Background(), config); err != nil {
			t.Fatal(err)
		}

		// Test write op
		configData := map[string]interface{}{
			"api_key":  "api_key",
			"base_uri": "https://notreal.launchdarkly.com",
		}
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "config",
			Data:      configData,
			Storage:   config.StorageView,
		})

		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("failed setup:. resp:%#v err:%v", resp, err)
		}

		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.ReadOperation,
			Path:      "config",
			Data:      configData,
			Storage:   config.StorageView,
		})

		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("failed to read config:. resp:%#v err:%v", resp, err)
		}

		expected := map[string]interface{}{
			"base_uri": "https://notreal.launchdarkly.com",
		}

		if diff := deep.Equal(expected, resp.Data); diff != nil {
			t.Fatal(diff)
		}
	})

	t.Run("test create config missing api_key", func(t *testing.T) {
		b := NewBackend(config.System)
		if err := b.Setup(context.Background(), config); err != nil {
			t.Fatal(err)
		}

		// Missing api_key
		configData := map[string]interface{}{
			"base_uri": "https://app.launchdarkly.com",
		}

		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "config",
			Data:      configData,
			Storage:   config.StorageView,
		})

		if err == nil {
			t.Fatalf("expected error but got nil")
		}
	})
}
