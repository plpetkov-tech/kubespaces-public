package tenant

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a tenant",
	Long:  `Update a tenant by updating the Helm release.`,
	Run: func(cmd *cobra.Command, args []string) {
		tenant, _ := cmd.Flags().GetString("tenant")
		org, _ := cmd.Flags().GetString("org")

		if tenant == "" || org == "" {
			fmt.Println("Error: tenant and org must be provided")
			return
		}

		updateHelmCmd := fmt.Sprintf("helm upgrade %s-%s .", tenant, org)
		fmt.Println("Running command:", updateHelmCmd)
		output, err := exec.Command("sh", "-c", updateHelmCmd).CombinedOutput()
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Output:", string(output))
			return
		}

		fmt.Println("Helm release updated successfully.")
	},
}

func init() {
	tenantCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringP("tenant", "t", "", "Tenant name")
	updateCmd.Flags().StringP("org", "o", "", "Organization name")
}
