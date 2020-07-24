// Package version is a convenience utility that provides SDK
// consumers with a ready-to-use version command that
// produces apps versioning information based on flags
// passed at compile time.
//
// Configure the version command
//
// The version command can be just added to your cobra root command.
// At build time, the variables Name, Version, Commit, and BuildTags
// can be passed as build flags as shown in the following example:
//
//  go build -X github.com/netcloth/netcloth-chain/version.Name=nch \
//   -X github.com/netcloth/netcloth-chain/version.ServerName=nchd \
//   -X github.com/netcloth/netcloth-chain/version.ClientName=nchcli \
//   -X github.com/netcloth/netcloth-chain/version.Version=1.0 \
//   -X github.com/netcloth/netcloth-chain/version.Commit=f0f7b7dab7e36c20b757cebce0e8f4fc5b95de60 \
//   -X "github.com/netcloth/netcloth-chain/version.BuildTags=linux darwin amd64"
package version

import (
	"fmt"
	"runtime"
)

var (
	Name       = ""
	ServerName = "<appd>"
	ClientName = "<appcli>"
	Version    = ""
	Commit     = ""
	BuildTags  = ""
	AppVersion = uint64(0)
)

type Info struct {
	Name       string `json:"name" yaml:"name"`
	ServerName string `json:"server_name" yaml:"server_name"`
	ClientName string `json:"client_name" yaml:"client_name"`
	Version    string `json:"version" yaml:"version"`
	AppVersion uint64 `json:"app_version" yaml:"app_version"`
	GitCommit  string `json:"commit" yaml:"commit"`
	BuildTags  string `json:"build_tags" yaml:"build_tags"`
	GoVersion  string `json:"go" yaml:"go"`
}

func NewInfo() Info {
	return Info{
		Name:       Name,
		ServerName: ServerName,
		ClientName: ClientName,
		Version:    Version,
		AppVersion: AppVersion,
		GitCommit:  Commit,
		BuildTags:  BuildTags,
		GoVersion:  fmt.Sprintf("go version %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
	}
}
