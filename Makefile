# Pulumi provider for Forumulate, bridged from ../terraform-provider-forumulate.
#
# After the Terraform provider gains a resource, regenerate everything with:
#
#     make tfgen build_sdks
#
# (tfgen rewrites schema.json from whatever the Terraform provider registers;
# build_sdks regenerates the Node.js, Python, Go and .NET SDKs from it.) Then
# `make provider` builds the plugin binary and `make test` runs the Go tests.
PACK            := forumulate
PROJECT         := github.com/forumulate/pulumi-forumulate
PROVIDER_PATH   := provider
VERSION_PATH    := $(PROJECT)/$(PROVIDER_PATH)/pkg/version.Version
TFGEN           := pulumi-tfgen-$(PACK)
PROVIDER        := pulumi-resource-$(PACK)
SCHEMA_DIR      := $(PROVIDER_PATH)/cmd/$(PROVIDER)
WORKING_DIR     := $(shell pwd)

# A development version unless the release passes VERSION=X.Y.Z; keeping the
# default fixed keeps the committed schema and SDKs byte-stable between runs.
VERSION         ?= 0.0.1-dev
# PEP 440: 1.2.3-rc1 -> 1.2.3rc1, 0.0.1-dev -> 0.0.1.dev0
PYPI_VERSION    := $(shell echo '$(VERSION)' | sed -E 's/-dev$$/.dev0/; s/-(alpha|a)\.?([0-9]*)/a\2/; s/-(beta|b)\.?([0-9]*)/b\2/; s/-rc\.?([0-9]*)/rc\1/')
LDFLAGS         := -X $(VERSION_PATH)=$(VERSION)

.PHONY: all tfgen provider build_sdks nodejs_sdk python_sdk go_sdk dotnet_sdk test lint_provider clean

all: tfgen provider build_sdks

# Schema: run the bridge over the Terraform provider. schema-embed.json is the
# copy the plugin binary embeds so it can serve its schema offline.
tfgen:
	cd $(PROVIDER_PATH) && go mod tidy
	cd $(PROVIDER_PATH) && go build -o $(WORKING_DIR)/bin/$(TFGEN) -ldflags "$(LDFLAGS)" ./cmd/$(TFGEN)
	./bin/$(TFGEN) schema --out $(SCHEMA_DIR)
	cp $(SCHEMA_DIR)/schema.json schema.json
	cp $(SCHEMA_DIR)/schema.json $(SCHEMA_DIR)/schema-embed.json

# Plugin binary for this machine; releases cross-build through goreleaser.
provider:
	cd $(PROVIDER_PATH) && go build -o $(WORKING_DIR)/bin/$(PROVIDER) -ldflags "$(LDFLAGS)" ./cmd/$(PROVIDER)

# SDKs, all from schema.json via the Pulumi CLI (pulumi package gen-sdk).
build_sdks: nodejs_sdk python_sdk go_sdk dotnet_sdk

nodejs_sdk:
	rm -rf sdk/nodejs
	pulumi package gen-sdk schema.json --language nodejs --out sdk
	cd sdk/nodejs && npm install --no-audit --no-fund && npm run build
	# bin/ is what gets published: tsc output plus the manifest, stamped.
	cp README.md LICENSE sdk/nodejs/package.json sdk/nodejs/bin/
	sed -i.bak 's/\$${VERSION}/$(VERSION)/g' sdk/nodejs/bin/package.json && rm sdk/nodejs/bin/package.json.bak

python_sdk:
	rm -rf sdk/python
	pulumi package gen-sdk schema.json --language python --out sdk
	cp README.md sdk/python/README.md
	sed -i.bak -E 's/^( *)version = ".*"/\1version = "$(PYPI_VERSION)"/' sdk/python/pyproject.toml && rm sdk/python/pyproject.toml.bak

# The Go SDK is its own module (sdk/go.mod) so consumers do not download the
# bridge; its go.sum is produced by `go mod tidy`.
go_sdk:
	rm -rf sdk/go
	pulumi package gen-sdk schema.json --language go --out sdk
	cd sdk && go mod tidy

dotnet_sdk:
	rm -rf sdk/dotnet
	pulumi package gen-sdk schema.json --language dotnet --out sdk
	echo "$(VERSION)" > sdk/dotnet/version.txt

test:
	cd $(PROVIDER_PATH) && go vet ./... && go test ./...

clean:
	rm -rf bin
