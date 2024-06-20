package token

import ldapi "github.com/launchdarkly/api-client-go"

type ServiceTokenStorageEntry struct {
	ID             string            `json:"id" structs:"id" mapstructure:"id"`
	Name           string            `json:"name" structs:"name" mapstructure:"name"`
	CustomRoleKeys []string          `json:"custom_roles" structs:"custom_roles" mapstructure:"custom_roles"`
	InlinePolicy   []ldapi.Statement `json:"inline_policy" structs:"inline_policy" mapstructure:"inline_policy"`
}
