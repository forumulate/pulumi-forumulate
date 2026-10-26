// Command pulumi-resource-forumulate is the Pulumi provider plugin: it serves
// the schema tfgen generated and bridges every call to the Terraform provider.
package main

import (
	"context"
	_ "embed"

	forumulate "github.com/forumulate/pulumi-forumulate/provider"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
)

//go:embed schema-embed.json
var pulumiSchema []byte

//go:embed bridge-metadata.json
var bridgeMetadata []byte

func main() {
	meta := tfbridge.ProviderMetadata{PackageSchema: pulumiSchema, BridgeMetadata: bridgeMetadata}
	tfbridge.Main(context.Background(), "forumulate", forumulate.Provider(), meta)
}
