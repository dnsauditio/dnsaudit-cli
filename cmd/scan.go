package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dnsaudit/cli/internal/api"
	"github.com/dnsaudit/cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	domain   string
	jsonOut  bool
)

// Terminal colors
const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorRed    = "\033[31m"
	colorOrange = "\033[33m" // Using yellow/orange ANSI
	colorBlue   = "\033[34m"
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

		if NoCache {
			_ = client.ClearCache(domain) // Silent request to clear cache
		}

		// Check if we are piped or redirecting
		fileInfo, _ := os.Stdout.Stat()
		isPiped := (fileInfo.Mode() & os.ModeCharDevice) == 0

		disableColors := isPiped || NoColor

		if !jsonOut {
			PrintBanner()
			// Even if silent, if not jsonOut and not piped, we might just print starting scan
			// Actually, if it's silent, maybe we shouldn't print "Starting scan"?
			// Let's keep existing logic: hideBanner = isPiped || Silent
			hideBanner := isPiped || Silent
			if !hideBanner {
				fmt.Printf("\n[*] Starting scan for %s...\n", domain)
			} else if !isPiped {
				fmt.Printf("[*] Starting scan for %s...\n", domain)
			}
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

		var actualCritical, actualWarning, actualInfo int
		var issuesList []interface{}
		
		if issues, ok := scanResult["issues"].([]interface{}); ok {
			issuesList = issues
			for _, issueInterface := range issuesList {
				if issueMap, ok := issueInterface.(map[string]interface{}); ok {
					if itype, ok := issueMap["type"].(string); ok {
						if itype == "critical" {
							actualCritical++
						} else if itype == "warning" {
							actualWarning++
						} else if itype == "info" || itype == "informational" {
							actualInfo++
						}
					}
				}
			}
		}

		if gradeMap, ok := scanResult["grade"].(map[string]interface{}); ok {
			if g, exists := gradeMap["grade"].(string); exists {
				gradeVal = g
			}
			if s, exists := gradeMap["score"].(float64); exists {
				scoreVal = s
			}

			if !disableColors {
				fmt.Printf("\n%s[+] Security Grade: %s (Score: %.0f)%s\n", colorBold, gradeVal, scoreVal, colorReset)
			} else {
				fmt.Printf("\n[+] Security Grade: %s (Score: %.0f)\n", gradeVal, scoreVal)
			}

			if desc, exists := gradeMap["description"].(string); exists {
				fmt.Printf("[*] %s\n", desc)
			}

			fmt.Println("\nIssues Breakdown:")
			fmt.Printf("  Critical: %v\n", actualCritical)
			fmt.Printf("  Warning:  %v\n", actualWarning)
			fmt.Printf("  Info:     %v\n", actualInfo)

			if recs, ok := gradeMap["recommendations"].([]interface{}); ok && len(recs) > 0 {
				fmt.Println("\nTop Recommendations:")
				for _, r := range recs {
					fmt.Printf("  - %v\n", r)
				}
			}
		} else {
			fmt.Printf("Grade:  %v (Score: %v)\n", scanResult["grade"], scanResult["securityScore"])
		}

		if len(issuesList) > 0 && (actualCritical > 0 || actualWarning > 0 || actualInfo > 0) {
			if !disableColors {
				fmt.Printf("\n%sKey Findings:%s\n\n", colorBold, colorReset)
			} else {
				fmt.Printf("\nKey Findings:\n\n")
			}
			count := 0
			for _, issueInterface := range issuesList {
				if count >= 8 {
					fmt.Println("  ... and more (export to JSON/PDF to see all)")
					break
				}
				if issueMap, ok := issueInterface.(map[string]interface{}); ok {
					itype := fmt.Sprintf("%v", issueMap["type"])
					rtype := fmt.Sprintf("%v", issueMap["recordType"])
					desc := issueMap["description"]

					// Extract just the title (the first line of the description)
					descStr := fmt.Sprintf("%v", desc)
					title := strings.SplitN(descStr, "\n", 2)[0]
					
					// Basic truncation just in case the title itself is incredibly long
					if len(title) > 80 {
						title = title[:77] + "..."
					}

					// We now show all types of issues, including informational
					if itype == "warning" || itype == "critical" || itype == "info" || itype == "informational" {
						typeStr := fmt.Sprintf("[%s]", itype)
						if !disableColors {
							switch itype {
							case "critical":
								typeStr = fmt.Sprintf("%s[%s]%s", colorRed, itype, colorReset)
							case "warning":
								typeStr = fmt.Sprintf("%s[%s]%s", colorOrange, itype, colorReset)
							case "info", "informational":
								typeStr = fmt.Sprintf("%s[%s]%s", colorBlue, "informational", colorReset) // Always display as informational
							}
						} else if itype == "info" {
							typeStr = "[informational]" // Normalize string output when piped
						}
						fmt.Printf("  %s %s: %s\n", typeStr, rtype, title)
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
