// Command pulumi-tfgen-forumulate generates the Pulumi schema and the language
// SDKs from the bridged Terraform provider.
package main

import (
	forumulate "github.com/forumulate/pulumi-forumulate/provider"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfgen"
)

func main() {
	tfgen.Main("forumulate", forumulate.Provider())
}
