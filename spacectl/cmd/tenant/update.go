package tenant

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a tenant",
	Run: func(cmd *cobra.Command, args []string) {
		tenant, _ := cmd.Flags().GetString("tenant")
		org, _ := cmd.Flags().GetString("org")
		k8sVersion, _ := cmd.Flags().GetString("k8s-version")
		cloud, _ := cmd.Flags().GetString("cloud")
		locationShort, _ := cmd.Flags().GetString("location-short")
		domain, _ := cmd.Flags().GetString("domain")

		if tenant == "" || org == "" {
			fmt.Println("Tenant and organization must be provided")
			return
		}

		// Run helm upgrade command
		helmCmd := exec.Command("helm", "upgrade", "-i", "-n", fmt.Sprintf("%s-%s", tenant, org), "--create-namespace",
			fmt.Sprintf("--set tenant.location_short=%s", locationShort),
			fmt.Sprintf("--set tenant.cloud=%s", cloud),
			fmt.Sprintf("--set tenant.name=%s", tenant),
			fmt.Sprintf("--set tenant.org=%s", org),
			fmt.Sprintf("--set tenant.domain=%s", domain),
			fmt.Sprintf("--set controlPlane.distro.k8s.version=%s", k8sVersion),
			fmt.Sprintf("--set vcluster.exportKubeConfig.server=https://api.%s.%s.%s.%s.%s", tenant, org, locationShort, cloud, domain),
			fmt.Sprintf("--set vcluster.controlPlane.proxy.extraSANs[0]=api.%s.%s.%s.%s.%s", tenant, org, locationShort, cloud, domain))
		helmOutput, err := helmCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error running helm upgrade command:", string(helmOutput))
			return
		}

		fmt.Println("Tenant updated successfully")
	},
}

func init() {
	updateCmd.Flags().StringP("tenant", "t", "", "Tenant name")
	updateCmd.Flags().StringP("org", "o", "", "Organization name")
	updateCmd.Flags().StringP("k8s-version", "k", "1.31.1", "Kubernetes version")
	updateCmd.Flags().StringP("cloud", "c", "azure", "Cloud provider")
	updateCmd.Flags().StringP("location-short", "l", "ne", "Location short")
	updateCmd.Flags().StringP("domain", "d", "kubespaces.cloud", "Domain")
}
