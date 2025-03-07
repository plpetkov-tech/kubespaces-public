package main

import (
	"os"

	"github.com/kubespaces/kubespaces-public/spacectl/cmd/tenant"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spacectl",
	Short: "A CLI tool for managing tenants",
}

func init() {
	rootCmd.AddCommand(tenant.TenantCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
