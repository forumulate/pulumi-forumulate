// Package forumulate is the Pulumi provider for Forumulate, bridged from the
// Terraform provider so the two can never drift: every resource and data
// source the Terraform provider registers appears here with no edit.
package forumulate

import (
	_ "embed" // bridge-metadata.json is compiled into the binary
	"fmt"
	"path/filepath"

	"github.com/forumulate/pulumi-forumulate/provider/pkg/version"
	tfprovider "github.com/forumulate/terraform-provider-forumulate/pkg/forumulate"
	pf "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge/tokens"
)

const (
	// mainPkg is the Pulumi package name; tokens read forumulate:index:Space.
	mainPkg = "forumulate"
	// mainMod is the one module every token lands in: the Terraform provider
	// has a flat namespace, so the Pulumi one does too.
	mainMod = "index"
)

//go:embed cmd/pulumi-resource-forumulate/bridge-metadata.json
var metadata []byte

// Provider returns the bridged provider description tfgen and the plugin share.
func Provider() tfbridge.ProviderInfo {
	info := tfbridge.ProviderInfo{
		P:            pf.ShimProvider(tfprovider.New(version.Version)()),
		Name:         "forumulate",
		MetadataInfo: tfbridge.NewProviderMetadata(metadata),
		Version:      version.Version,
		DisplayName:  "Forumulate",
		Publisher:    "Forumulate",
		LogoURL:      "https://forumulate.com/favicon.svg",
		Description:  "A Pulumi package for managing Forumulate docs, forums, blogs and status pages as code.",
		Keywords:     []string{"pulumi", "forumulate", "docs", "forum", "blog", "status", "category/cloud"},
		License:      "MIT",
		Homepage:     "https://forumulate.com",
		Repository:   "https://github.com/forumulate/pulumi-forumulate",
		GitHubOrg:    "forumulate",
		// Not on the Pulumi registry: `pulumi plugin install resource forumulate
		// vX --server github://api.github.com/forumulate` fetches the
		// release assets goreleaser publishes on this repo.
		PluginDownloadURL: "github://api.github.com/forumulate",

		// The Terraform provider already reads these from the environment, but
		// declaring them lets `pulumi config` and the program see the same
		// defaults, and marks the token secret.
		Config: map[string]*tfbridge.SchemaInfo{
			"token":   {Secret: tfbridge.True(), Default: &tfbridge.DefaultInfo{EnvVars: []string{"FORUMULATE_TOKEN"}}},
			"org_id":  {Default: &tfbridge.DefaultInfo{EnvVars: []string{"FORUMULATE_ORG_ID"}}},
			"api_url": {Default: &tfbridge.DefaultInfo{EnvVars: []string{"FORUMULATE_API_URL"}}},
		},

		JavaScript: &tfbridge.JavaScriptInfo{
			PackageName: "@forumulate/pulumi",
			Dependencies: map[string]string{
				"@pulumi/pulumi": "^3.0.0",
			},
			DevDependencies: map[string]string{
				"@types/node": "^20.0.0",
			},
		},
		Python: &tfbridge.PythonInfo{
			PackageName: "pulumi_forumulate",
			Requires: map[string]string{
				"pulumi": ">=3.0.0,<4.0.0",
			},
			PyProject: struct{ Enabled bool }{Enabled: true},
		},
		Golang: &tfbridge.GolangInfo{
			ImportBasePath: filepath.Join(
				fmt.Sprintf("github.com/forumulate/pulumi-%[1]s/sdk/", mainPkg),
				tfbridge.GetModuleMajorVersion(version.Version),
				"go",
				mainPkg,
			),
			GenerateResourceContainerTypes: true,
		},
		CSharp: &tfbridge.CSharpInfo{
			RootNamespace: "Forumulate",
			PackageReferences: map[string]string{
				"Pulumi": "3.*",
			},
		},
	}

	// Automatic token mapping: forumulate_status_component ->
	// forumulate:index:StatusComponent, data sources -> getStatusPage. A
	// resource the Terraform provider gains tomorrow maps the same way, so
	// this file does not change when it grows.
	info.MustComputeTokens(tokens.SingleModule("forumulate_", mainMod, tokens.MakeStandard(mainPkg)))
	// Renames in the Terraform provider become Pulumi aliases, recorded in
	// bridge-metadata.json, so existing stacks keep their resources.
	info.MustApplyAutoAliases()
	info.SetAutonaming(255, "-")

	return info
}
