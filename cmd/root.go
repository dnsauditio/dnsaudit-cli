package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev"

var (
	Silent  bool
	NoColor bool
	NoCache bool
)

var rootCmd = &cobra.Command{
	Use:   "dnsaudit",
	Short: "DNSAudit.io CLI - Professional DNS Security Scanning",
	Long:  "A fast and robust command line tool for interacting with the DNSAudit.io API.\nAudit domains, detect misconfigurations and monitor DNS changes from your terminal.",
}

func GetBanner() string {
	// Check if we are piped or redirecting
	fileInfo, _ := os.Stdout.Stat()
	isPiped := (fileInfo.Mode() & os.ModeCharDevice) == 0

	disableColors := isPiped || NoColor
	hideBanner := isPiped || Silent

	if hideBanner {
		return ""
	}

	if !disableColors {
		return fmt.Sprintf("+---------------------------------------+\n|   \033[1mDNSAudit.io  CLI  %s\033[0m            |\n|   DNS Security from the terminal      |\n+---------------------------------------+\n", Version)
	}
	return fmt.Sprintf("+---------------------------------------+\n|   DNSAudit.io  CLI  %s            |\n|   DNS Security from the terminal      |\n+---------------------------------------+\n", Version)
}

func PrintBanner() {
	fmt.Print(GetBanner())
}

func Execute() {
	cobra.AddTemplateFunc("banner", GetBanner)
	
	// Set global help and usage templates to include the banner
	customHelpTemplate := `{{banner}}
{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

	rootCmd.SetHelpTemplate(customHelpTemplate)
	// Usage is also printed when no args are provided
	rootCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Silent, "silent", "s", false, "Hide the CLI banner and loading text")
	rootCmd.PersistentFlags().BoolVar(&NoColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().BoolVar(&NoCache, "no-cache", false, "Bypass cached results and force a fresh scan")
}
