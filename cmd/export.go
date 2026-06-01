package cmd

import (
	"fmt"
	"os"

	"github.com/dnsaudit/cli/internal/api"
	"github.com/dnsaudit/cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	format     string
	outputFile string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export scan results to PDF or JSON",
	Run: func(cmd *cobra.Command, args []string) {
		if domain == "" {
			fmt.Println("Error: Domain is required. Use -d <domain>")
			os.Exit(1)
		}

		if format != "pdf" && format != "json" {
			fmt.Println("Error: Invalid format. Use 'pdf' or 'json'")
			os.Exit(1)
		}

		apiKey, err := config.Load()
		if err != nil || apiKey == "" {
			fmt.Println("Error: API key not found. Run 'dnsaudit configure' or set DNSAUDIT_API_KEY")
			os.Exit(1)
		}

		client := api.NewClient(apiKey)
		
		fmt.Printf("[*] Exporting %s report for %s...\n", format, domain)

		var result []byte
		if format == "json" {
			result, err = client.ExportJSON(domain)
		} else {
			result, err = client.ExportPDF(domain)
		}

		if err != nil {
			fmt.Printf("Export failed: %v\n", err)
			os.Exit(1)
		}

		outPath := outputFile
		if outPath == "" {
			outPath = fmt.Sprintf("%s-report.%s", domain, format)
		}

		err = os.WriteFile(outPath, result, 0644)
		if err != nil {
			fmt.Printf("Failed to write file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("[+] Successfully saved %s report to %s\n", format, outPath)
	},
}

func init() {
	exportCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain to export (required)")
	exportCmd.Flags().StringVarP(&format, "format", "f", "pdf", "Export format: 'pdf' or 'json'")
	exportCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: <domain>-report.<format>)")
	rootCmd.AddCommand(exportCmd)
}
