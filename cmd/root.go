package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dnsaudit",
	Short: "DNSAudit.io CLI - Professional DNS Security Scanning",
	Long:  `A fast and robust command line tool for interacting with the DNSAudit.io API. Built for DevOps teams and penetration testers.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
