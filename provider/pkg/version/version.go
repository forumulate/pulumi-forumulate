// Package version carries the release version, stamped at link time:
// -ldflags "-X github.com/forumulate/pulumi-forumulate/provider/pkg/version.Version=v0.1.0".
package version

// Version is the provider version. The default is a valid semver because the
// Plugin Framework bridge refuses to start without one, so an unstamped
// `go build` still yields a working plugin.
var Version = "0.0.1-dev"
