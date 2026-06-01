package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dnsaudit/cli/internal/api"
	"github.com/dnsaudit/cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	domain   string
	jsonOut  bool
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a DNS security scan on a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if domain == "" {
			fmt.Println("Error: Domain is required. Use -d <domain>")
			os.Exit(1)
		}

		apiKey, err := config.Load()
		if err != nil || apiKey == "" {
			fmt.Println("Error: API key not found. Run 'dnsaudit configure' or set DNSAUDIT_API_KEY")
			os.Exit(1)
		}

		client := api.NewClient(apiKey)
		
		if !jsonOut {
			fmt.Printf("[*] Starting scan for %s...\n", domain)
		}

		result, err := client.Scan(domain)
		if err != nil {
			fmt.Printf("Scan failed: %v\n", err)
			os.Exit(1)
		}

		if jsonOut {
			fmt.Println(string(result))
			return
		}

		// Parse JSON to format nice output
		var scanResult map[string]interface{}
		if err := json.Unmarshal(result, &scanResult); err != nil {
			fmt.Printf("Failed to parse response: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Printf("Target: %v\n", scanResult["domain"])
		fmt.Printf("Status: %v\n", scanResult["status"])
		
		var gradeVal string
		var scoreVal float64
		
		if gradeMap, ok := scanResult["grade"].(map[string]interface{}); ok {
			if g, exists := gradeMap["grade"].(string); exists {
				gradeVal = g
			}
			if s, exists := gradeMap["score"].(float64); exists {
				scoreVal = s
			}
			
			fmt.Printf("\n[+] Security Grade: %s (Score: %.0f)\n", gradeVal, scoreVal)
			if desc, exists := gradeMap["description"].(string); exists {
				fmt.Printf("[*] %s\n", desc)
			}
			
			if bd, ok := gradeMap["breakdown"].(map[string]interface{}); ok {
				fmt.Println("\nIssues Breakdown:")
				fmt.Printf("  Critical: %v\n", bd["critical"])
				fmt.Printf("  Warning:  %v\n", bd["warning"])
				fmt.Printf("  Info:     %v\n", bd["info"])
			}
			
			if recs, ok := gradeMap["recommendations"].([]interface{}); ok && len(recs) > 0 {
				fmt.Println("\nTop Recommendations:")
				for _, r := range recs {
					fmt.Printf("  - %v\n", r)
				}
			}
		} else {
			fmt.Printf("Grade:  %v (Score: %v)\n", scanResult["grade"], scanResult["securityScore"])
		}

		if issues, ok := scanResult["issues"].([]interface{}); ok && len(issues) > 0 {
			fmt.Println("\nKey Findings:")
			count := 0
			for _, issueInterface := range issues {
				if count >= 5 {
					fmt.Println("  ... and more (export to JSON/PDF to see all)")
					break
				}
				if issueMap, ok := issueInterface.(map[string]interface{}); ok {
					itype := issueMap["type"]
					rtype := issueMap["recordType"]
					desc := issueMap["description"]
					
					// Remove newlines and truncate
					descStr := fmt.Sprintf("%v", desc)
					
					// Basic newline removal for clean terminal output
					descStrClean := ""
					for _, ch := range descStr {
						if ch != '\n' && ch != '\r' {
							descStrClean += string(ch)
						} else if ch == '\n' {
							descStrClean += " - "
						}
					}
					
					if len(descStrClean) > 80 {
						descStrClean = descStrClean[:77] + "..."
					}
					
					// Only show warning or critical in the summary to keep it clean
					if itype == "warning" || itype == "critical" {
						fmt.Printf("  [%s] %s: %s\n", itype, rtype, descStrClean)
						count++
					}
				}
			}
		}
		fmt.Println("--------------------------------------------------------------------------------")
	},
}

func init() {
	scanCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain to scan (required)")
	scanCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "Output raw JSON instead of formatted text")
	rootCmd.AddCommand(scanCmd)
}
