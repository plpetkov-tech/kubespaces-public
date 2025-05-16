package tenant

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
		releaseNameForHelm := fmt.Sprintf("%s-%s", tenant, org) // Combine tenant and org name
		// chartPath := "oci://ghcr.io/kubespaces-io/kubespaces-public/tenant" // Use the OCI chart path
		// Or for local development:
		chartPath := "../charts/tenant" // Path to the chart in the repo
		namespace := fmt.Sprintf("%s-%s", tenant, org)
		var exportKubeconfigServer string
		var exportKubeconfigSANs string
		if cloud == "kind" {
			exportKubeconfigServer = "vcluster.exportKubeConfig.server=https://localhost:8443"
			exportKubeconfigSANs = "vcluster.controlPlane.proxy.extraSANs[0]=localhost"
		} else {
			exportKubeconfigServer = fmt.Sprintf("vcluster.exportKubeConfig.server=https://api.%s.%s.%s.%s.%s", tenant, org, locationShort, cloud, domain)
			exportKubeconfigSANs = fmt.Sprintf("vcluster.controlPlane.proxy.extraSANs[0]=api.%s.%s.%s.%s.%s", tenant, org, locationShort, cloud, domain)
		}
		// Run helm upgrade command
		helmCmd := exec.Command(
			"helm", "upgrade", "--install",
			"--create-namespace",
			"-n", namespace,
			releaseNameForHelm, // First argument - release name
			chartPath,
			"--set", fmt.Sprintf("tenant.location_short=%s", locationShort),
			"--set", fmt.Sprintf("tenant.cloud=%s", cloud),
			"--set", fmt.Sprintf("tenant.name=%s", tenant),
			"--set", fmt.Sprintf("tenant.org=%s", org),
			"--set", fmt.Sprintf("tenant.domain=%s", domain),
			"--set", fmt.Sprintf("controlPlane.distro.k8s.version=%s", k8sVersion),
			"--set", exportKubeconfigServer, 
			"--set", exportKubeconfigSANs)
		helmOutput, err := helmCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error running helm upgrade command:", string(helmOutput))
			return
		}

		if wait {
			secretName := fmt.Sprintf("vc-%s-%s", tenant, org)
			namespace := fmt.Sprintf("%s-%s", tenant, org)

			fmt.Printf("Waiting for secret '%s' in namespace '%s'...\n", secretName, namespace)
			for {
				execCmd := exec.Command("kubectl", "get", "secret", secretName, "-n", namespace)
				if err := execCmd.Run(); err == nil {
					// Secret exists, now get its config
					execCmd = exec.Command("kubectl", "get", "secret", secretName, "-n", namespace, "-o", "jsonpath={.data.config}")
					output, err := execCmd.Output()
					if err != nil {
						fmt.Printf("Error getting kubeconfig: %s\n", err)
						return
					}
					if len(output) > 0 {
						// Decode the base64 data
						decodedConfig, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(output)))
						if err != nil {
							fmt.Println("Error decoding kubeconfig:", err)
							return
						}

						kubeconfig := string(decodedConfig)
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
						fmt.Println("Secret found and kubeconfig processed successfully!")
						break
					}
				}
				fmt.Print("Waiting for secret to be available. Press CTRL-C to exit...\r")
				time.Sleep(1 * time.Second)
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
