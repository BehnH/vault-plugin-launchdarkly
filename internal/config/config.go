package config

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
)

const PathPatternConfig = "config"

type ConfigEntry struct {
	BaseUrl         string `json:"base_url" structs:"base_url" mapstructure:"base_url"`
	Token           string `json:"token" structs:"token" mapstructure:"token"`
	AllowAdminToken bool   `json:"allow_admin_token" structs:"allow_admin_token" mapstructure:"allow_admin_token"`
}

func GetConfig(ctx context.Context, s logical.Storage) (*ConfigEntry, error) {
	var config ConfigEntry
	configRaw, err := s.Get(ctx, PathPatternConfig)
	if err != nil {
		return nil, err
	}
	if configRaw == nil {
		return nil, nil
	}

	if err := configRaw.DecodeJSON(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
