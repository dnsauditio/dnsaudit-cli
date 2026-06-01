package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dnsaudit/cli/internal/config"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure your API key",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your DNSAudit.io API Key: ")
		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)

		if apiKey == "" {
			fmt.Println("API key cannot be empty.")
			os.Exit(1)
		}

		err := config.Save(apiKey)
		if err != nil {
			fmt.Printf("Failed to save configuration: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Configuration saved successfully! You can now use the CLI.")
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
