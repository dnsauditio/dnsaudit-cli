package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Silent  bool
	NoColor bool
)

var rootCmd = &cobra.Command{
	Use:   "dnsaudit",
	Short: "DNSAudit.io CLI - Professional DNS Security Scanning",
	Long:  "A fast and robust command line tool for interacting with the DNSAudit.io API.\nAudit domains, detect misconfigurations and monitor DNS changes from your terminal.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Silent, "silent", "s", false, "Hide the CLI banner and loading text")
	rootCmd.PersistentFlags().BoolVar(&NoColor, "no-color", false, "Disable colored output")
}
