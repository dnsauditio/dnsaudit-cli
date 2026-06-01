package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/dnsaudit/cli/internal/api"
	"github.com/dnsaudit/cli/internal/config"
	"github.com/spf13/cobra"
)

var limit int

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View your recent DNS scan history",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, err := config.Load()
		if err != nil || apiKey == "" {
			fmt.Println("Error: API key not found. Run 'dnsaudit configure' or set DNSAUDIT_API_KEY")
			os.Exit(1)
		}

		client := api.NewClient(apiKey)
		
		if !jsonOut {
			fmt.Println("[*] Fetching scan history...")
		}

		result, err := client.History(limit)
		if err != nil {
			fmt.Printf("Failed to fetch history: %v\n", err)
			os.Exit(1)
		}

		if jsonOut {
			fmt.Println(string(result))
			return
		}

		var historyData struct {
			Scans []struct {
				Domain        string `json:"domain"`
				ScanDate      string `json:"scanDate"`
				SecurityScore int         `json:"securityScore"`
				Grade         interface{} `json:"grade"`
			} `json:"scans"`
			Total int `json:"total"`
		}

		if err := json.Unmarshal(result, &historyData); err != nil {
			fmt.Printf("Failed to parse response: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nTotal Scans in History: %d\n", historyData.Total)
		fmt.Println("--------------------------------------------------------------------------------")
		
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "DOMAIN\tDATE\tSCORE\tGRADE")
		fmt.Fprintln(w, "------\t----\t-----\t-----")
		
		for _, s := range historyData.Scans {
			var gradeStr string
			if gMap, ok := s.Grade.(map[string]interface{}); ok {
				if g, exists := gMap["grade"].(string); exists {
					gradeStr = g
				}
			} else if g, ok := s.Grade.(string); ok {
				gradeStr = g
			}
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", s.Domain, s.ScanDate, s.SecurityScore, gradeStr)
		}
		w.Flush()
		fmt.Println("--------------------------------------------------------------------------------")
	},
}

func init() {
	historyCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Number of results to return (max 100)")
	historyCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "Output raw JSON instead of formatted text")
	rootCmd.AddCommand(historyCmd)
}
