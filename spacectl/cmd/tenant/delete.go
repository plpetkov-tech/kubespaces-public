package tenant

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a tenant",
	Run: func(cmd *cobra.Command, args []string) {
		tenant, _ := cmd.Flags().GetString("tenant")
		org, _ := cmd.Flags().GetString("org")

		if tenant == "" || org == "" {
			fmt.Println("Tenant and organization must be provided")
			return
		}

		// Delete the namespace
		deleteNamespaceCmd := exec.Command("kubectl", "delete", "namespace", fmt.Sprintf("%s-%s", tenant, org))
		deleteNamespaceOutput, err := deleteNamespaceCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error deleting namespace:", string(deleteNamespaceOutput))
			return
		}

		// Delete the TLSRoute in the istio-system namespace
		deleteTLSRouteCmd := exec.Command("kubectl", "delete", "tlsroute", fmt.Sprintf("%s-%s-tlsroute", tenant, org), "-n", "istio-system")
		deleteTLSRouteOutput, err := deleteTLSRouteCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error deleting TLSRoute:", string(deleteTLSRouteOutput))
			return
		}

		fmt.Println("Tenant deleted successfully")
	},
}

func init() {
	deleteCmd.Flags().StringP("tenant", "t", "", "Tenant name")
	deleteCmd.Flags().StringP("org", "o", "", "Organization name")
}
