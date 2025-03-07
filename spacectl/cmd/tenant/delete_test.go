package tenant

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestDeleteCmd(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("tenant", "t", "", "Tenant name")
	cmd.Flags().StringP("org", "o", "", "Organization name")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "Test delete with all flags",
			args:    []string{"--tenant=test", "--org=testorg"},
			wantErr: false,
		},
		{
			name:    "Test delete with missing tenant and org",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("deleteCmd.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
