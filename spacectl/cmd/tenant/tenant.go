package tenant

import (
	"github.com/spf13/cobra"
)

var tenantCmd = &cobra.Command{
	Use:   "tenant",
	Short: "Manage tenants",
	Long:  `Manage tenants in the system`,
}

func init() {
	tenantCmd.AddCommand(updateCmd)
	tenantCmd.AddCommand(deleteCmd)
}
