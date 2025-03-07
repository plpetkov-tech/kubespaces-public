package tenant

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a tenant",
	Long:  `Delete a tenant by removing the namespace and the TLSRoute in the istio-system namespace.",
	Run: func(cmd *cobra.Command, args []string) {
		tenant, _ := cmd.Flags().GetString("tenant")
		org, _ := cmd.Flags().GetString("org")

		if tenant == "" || org == "" {
			fmt.Println("Error: tenant and org must be provided")
			return
		}

		deleteNamespaceCmd := fmt.Sprintf("kubectl delete namespace %s-%s", tenant, org)
		fmt.Println("Running command:", deleteNamespaceCmd)
		output, err := exec.Command("sh", "-c", deleteNamespaceCmd).CombinedOutput()
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Output:", string(output))
			return
		}

		deleteTLSRouteCmd := fmt.Sprintf("kubectl delete tlsroute %s-%s -n istio-system", tenant, org)
		fmt.Println("Running command:", deleteTLSRouteCmd)
		output, err = exec.Command("sh", "-c", deleteTLSRouteCmd).CombinedOutput()
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Output:", string(output))
			return
		}

		fmt.Println("Tenant deleted successfully.")
	},
}

func init() {
	tenantCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringP("tenant", "t", "", "Tenant name")
	deleteCmd.Flags().StringP("org", "o", "", "Organization name")
}
