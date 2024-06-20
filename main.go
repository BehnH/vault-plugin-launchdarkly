package main

import (
	"log"
	"os"

	"github.com/behnh/vault-plugin-launchdarkly/internal/backend"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	apiClientMetadata := &api.PluginAPIClientMeta{}
	flags := apiClientMetadata.FlagSet()
	_ = flags.Parse(os.Args[1:])

	tlsConf := apiClientMetadata.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConf)

	log.Printf("vault-plugin-launchdarkly %s, commit %s, built at %s", version, commit, date)
	if err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: backend.Init,
		TLSProviderFunc:    tlsProviderFunc,
	}); err != nil {
		log.Fatal(err)
	}
}
