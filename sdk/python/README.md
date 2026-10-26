# Pulumi provider for Forumulate

Manage [Forumulate](https://forumulate.com) as code from Pulumi, in TypeScript, Python, Go or C#.

This provider is **bridged from the Terraform provider**
([`terraform-provider-forumulate`](https://github.com/forumulate/terraform-provider-forumulate))
with `pulumi-terraform-bridge`. Every resource and data source the Terraform provider has, this one
has, with the same arguments and the same behaviour; nothing is written twice, so the two cannot
drift. Resource tokens are mapped automatically: `forumulate_status_component` is
`forumulate:index:StatusComponent`, the data source `forumulate_status_page` is
`forumulate:index:getStatusPage`.

## Install

The provider plugin is published as GitHub release assets (it is not on the Pulumi registry), so
install it with a server URL:

```sh
pulumi plugin install resource forumulate v0.1.0 --server github://api.github.com/forumulate
```

Then add the SDK for your language:

| Language | Package | Install |
| --- | --- | --- |
| Node.js / TypeScript | [`@forumulate/pulumi`](https://www.npmjs.com/package/@forumulate/pulumi) | `npm install @forumulate/pulumi` |
| Python | [`pulumi_forumulate`](https://pypi.org/project/pulumi-forumulate/) | `pip install pulumi_forumulate` |
| Go | `github.com/forumulate/pulumi-forumulate/sdk` | `go get github.com/forumulate/pulumi-forumulate/sdk` |
| .NET | [`Forumulate.Forumulate`](https://www.nuget.org/packages/Forumulate.Forumulate) | `dotnet add package Forumulate.Forumulate` |

Pulumi downloads the plugin on its own when a program uses an SDK whose version it does not have;
the SDKs carry the server URL, so the explicit `plugin install` is only needed on CI images you
pre-warm.

## Authentication

| Setting | Environment variable | Pulumi config | |
| --- | --- | --- | --- |
| Personal access token | `FORUMULATE_TOKEN` | `forumulate:token` (secret) | Needed by every managed resource; the `getStatusPage` data source reads the public page without one. |
| Organization id | `FORUMULATE_ORG_ID` | `forumulate:orgId` | Required. |
| API base | `FORUMULATE_API_URL` | `forumulate:apiUrl` | Optional; `https://forumulate.com/api` by default. |

```sh
export FORUMULATE_TOKEN=...      # a personal access token from your account settings
export FORUMULATE_ORG_ID=...
# or, kept in the stack's encrypted config:
pulumi config set --secret forumulate:token ...
pulumi config set forumulate:orgId ...
```

## Examples

TypeScript:

```ts
import * as forumulate from "@forumulate/pulumi";

const spaceId = "your-status-space-id";

const api = new forumulate.StatusComponent("api", {
    spaceId,
    name: "API",
    url: "https://api.example.com/health",
    group: "Core",
    regions: ["eu-west", "us-east"],
    confirmRegions: 2,
    intervalSeconds: 60,
});

export const apiPing = api.id;
```

Python:

```python
import pulumi
import pulumi_forumulate as forumulate

api = forumulate.StatusComponent(
    "api",
    space_id="your-status-space-id",
    name="API",
    url="https://api.example.com/health",
    group="Core",
    regions=["eu-west", "us-east"],
    confirm_regions=2,
    interval_seconds=60,
)

pulumi.export("api", api.id)
```

Go:

```go
package main

import (
	forumulate "github.com/forumulate/pulumi-forumulate/sdk/go/forumulate"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		api, err := forumulate.NewStatusComponent(ctx, "api", &forumulate.StatusComponentArgs{
			SpaceId:         pulumi.String("your-status-space-id"),
			Name:            pulumi.String("API"),
			Url:             pulumi.String("https://api.example.com/health"),
			Group:           pulumi.String("Core"),
			IntervalSeconds: pulumi.Int(60),
		})
		if err != nil {
			return err
		}
		ctx.Export("api", api.ID())
		return nil
	})
}
```

Runnable versions are in [`examples/`](examples). The Terraform provider's
[README](https://github.com/forumulate/terraform-provider-forumulate) lists each resource, what it
manages and its import id (`pulumi import forumulate:index/statusComponent:StatusComponent api <space_id>/<component_id>`).

## Build and regenerate

Layout:

```
provider/            the bridge: resources.go (token mapping), cmd/pulumi-resource-forumulate (plugin),
                     cmd/pulumi-tfgen-forumulate (schema + SDK generator)
schema.json          generated Pulumi schema (provider/cmd/.../schema-embed.json is the copy the plugin embeds)
sdk/{nodejs,python,go,dotnet}   generated SDKs
examples/            one program per language
ci/cloudbuild.yaml   the same generation on Google Cloud Build (when a laptop should not download the bridge)
```

When the Terraform provider gains a resource, nothing in `provider/` is edited. Regenerate:

```sh
make tfgen build_sdks     # schema.json, then all four SDKs (needs the pulumi CLI, node, python, dotnet)
make provider             # bin/pulumi-resource-forumulate for this machine
make test                 # go vet + go test for provider/
```

`provider/go.mod` has a `replace` pointing at `../../terraform-provider-forumulate` while both live in
the Forumulate monorepo; the public mirror drops it and pins a released Terraform provider version.

## Releasing

Push a `vX.Y.Z` tag. `.github/workflows/release.yml` builds the plugin for darwin, linux and
windows (amd64, arm64) with GoReleaser and attaches them to the GitHub release, then publishes the
SDKs: npm needs the `NPM_TOKEN` repository secret, PyPI `PYPI_API_TOKEN`, NuGet `NUGET_PUBLISH_KEY`.
A publish job whose secret is not set is skipped, so the release still succeeds and the job can be
re-run after the secret is added. The Go SDK needs no publishing: it is the `sdk/vX.Y.Z` tag the
workflow pushes.

## License

MIT.
