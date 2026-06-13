package main

import "github.com/dnsaudit/cli/cmd"

// This variable is dynamically injected by GoReleaser at build time via ldflags
var version = "dev"

func main() {
	cmd.Version = version
	cmd.Execute()
}
