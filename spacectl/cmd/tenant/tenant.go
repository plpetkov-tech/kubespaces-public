package tenant

import (
	"github.com/spf13/cobra"
)

var TenantCmd = &cobra.Command{
	Use:   "tenant",
	Short: "Manage tenants",
	Long:  `Manage tenants in the system`,
}

func init() {
	TenantCmd.AddCommand(createCmd)
	TenantCmd.AddCommand(updateCmd)
	TenantCmd.AddCommand(deleteCmd)
}
