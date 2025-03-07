package tenant

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	"github.com/thanhpk/randstr"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a tenant",
	Run: func(cmd *cobra.Command, args []string) {
		tenant, _ := cmd.Flags().GetString("tenant")
		org, _ := cmd.Flags().GetString("org")
		k8sVersion, _ := cmd.Flags().GetString("k8s-version")
		cloud, _ := cmd.Flags().GetString("cloud")
		locationShort, _ := cmd.Flags().GetString("location-short")
		domain, _ := cmd.Flags().GetString("domain")
		wait, _ := cmd.Flags().GetBool("wait")
		outputFile, _ := cmd.Flags().GetString("output-file")

		if tenant == "" {
			tenant = randstr.String(5)
		}
		if org == "" {
			org = randstr.String(5)
		}

		// Run helm upgrade command
		helmCmd := exec.Command("helm", "upgrade", "-i", "-n", fmt.Sprintf("%s-%s", tenant, org), "--create-namespace",
			"--set", fmt.Sprintf("tenant.location_short=%s", locationShort),
			"--set", fmt.Sprintf("tenant.cloud=%s", cloud),
			"--set", fmt.Sprintf("tenant.name=%s", tenant),
			"--set", fmt.Sprintf("tenant.org=%s", org),
			"--set", fmt.Sprintf("tenant.domain=%s", domain),
			"--set", fmt.Sprintf("controlPlane.distro.k8s.version=%s", k8sVersion),
			"--set", fmt.Sprintf("vcluster.exportKubeConfig.server=https://api.%s.%s.%s.%s.%s", tenant, org, locationShort, cloud, domain),
			"--set", fmt.Sprintf("vcluster.controlPlane.proxy.extraSANs[0]=api.%s.%s.%s.%s.%s", tenant, org, locationShort, cloud, domain))
		helmOutput, err := helmCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error running helm upgrade command:", string(helmOutput))
			return
		}

		if wait {
			secretName := fmt.Sprintf("%s-%s-kubeconfig", tenant, org)
			for {
				execCmd := exec.Command("kubectl", "get", "secret", secretName, "-n", fmt.Sprintf("%s-%s", tenant, org), "-o", "jsonpath='{.data.kubeconfig}'")
				output, err := execCmd.Output()
				if err == nil && len(output) > 0 {
					kubeconfig := string(output)
					if outputFile != "" {
						err := os.WriteFile(outputFile, []byte(kubeconfig), 0644)
						if err != nil {
							fmt.Println("Error writing kubeconfig to file:", err)
							return
						}
					} else {
						kubeconfigPath := os.ExpandEnv("$HOME/.kube/config")
						f, err := os.OpenFile(kubeconfigPath, os.O_APPEND|os.O_WRONLY, 0644)
						if err != nil {
							fmt.Println("Error opening kubeconfig file:", err)
							return
						}
						defer f.Close()
						_, err = f.WriteString(kubeconfig)
						if err != nil {
							fmt.Println("Error writing to kubeconfig file:", err)
							return
						}
					}
					break
				}
				time.Sleep(5 * time.Second)
			}
		}

		fmt.Println("Tenant created successfully")
	},
}

func init() {
	createCmd.Flags().StringP("tenant", "t", "", "Tenant name")
	createCmd.Flags().StringP("org", "o", "", "Organization name")
	createCmd.Flags().StringP("k8s-version", "k", "1.31.1", "Kubernetes version")
	createCmd.Flags().StringP("cloud", "c", "azure", "Cloud provider")
	createCmd.Flags().StringP("location-short", "l", "ne", "Location short")
	createCmd.Flags().StringP("domain", "d", "kubespaces.cloud", "Domain")
	createCmd.Flags().BoolP("wait", "w", false, "Wait for the secret to be created and provide the kubeconfig file")
	createCmd.Flags().StringP("output-file", "f", "", "Output file for the kubeconfig")
}
