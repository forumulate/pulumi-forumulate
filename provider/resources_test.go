package forumulate

import (
	"os"
	"strings"
	"testing"

	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim"
)

// The whole point of the bridge: nothing the Terraform provider registers may
// be left without a Pulumi token, whatever it grows into.
func TestEveryTerraformEntityHasAToken(t *testing.T) {
	info := Provider()

	resources := 0
	info.P.ResourcesMap().Range(func(key string, _ shim.Resource) bool {
		resources++
		r, ok := info.Resources[key]
		if !ok || !strings.HasPrefix(r.Tok.String(), "forumulate:index/") {
			t.Errorf("resource %s has no forumulate:index token (got %v)", key, r)
		}
		return true
	})
	if resources == 0 {
		t.Fatal("the Terraform provider registered no resources")
	}

	info.P.DataSourcesMap().Range(func(key string, _ shim.Resource) bool {
		d, ok := info.DataSources[key]
		if !ok || !strings.HasPrefix(d.Tok.String(), "forumulate:index/get") {
			t.Errorf("data source %s has no forumulate:index/get token (got %v)", key, d)
		}
		return true
	})
}

// The committed schema is what the plugin serves and what the SDKs came from:
// every token must appear in it, or someone forgot `make tfgen build_sdks`.
func TestCommittedSchemaIsCurrent(t *testing.T) {
	raw, err := os.ReadFile("../schema.json")
	if err != nil {
		t.Fatal(err)
	}
	schema := string(raw)
	for key, r := range Provider().Resources {
		if !strings.Contains(schema, `"`+r.Tok.String()+`"`) {
			t.Errorf("%s (%s) is missing from schema.json; run `make tfgen build_sdks`", key, r.Tok)
		}
	}
}
