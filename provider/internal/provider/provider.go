package provider

import (
	"context"
	"os"

	dadcorp "dadcorp.dev/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const baseURL = "http://localhost:12345"

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"dadcorp_access_policy":       resourceAccessPolicy(),
			"dadcorp_consul_cluster":      resourceConsulCluster(),
			"dadcorp_ip":                  resourceIP(),
			"dadcorp_nomad_cluster":       resourceNomadCluster(),
			"dadcorp_terraform_workspace": resourceTerraformWorkspace(),
			"dadcorp_vault_cluster":       resourceVaultCluster(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"dadcorp_access_policy_consul":    dataAccessPolicyConsul(),
			"dadcorp_access_policy_nomad":     dataAccessPolicyNomad(),
			"dadcorp_access_policy_terraform": dataAccessPolicyTerraform(),
			"dadcorp_access_policy_vault":     dataAccessPolicyVault(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

type clientFactory struct {
	username, password string
}

func (c *clientFactory) NewClient() (*dadcorp.Client, error) {
	return dadcorp.NewClient(baseURL, c.username, c.password)
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	if os.Getenv("DADCORP_USERNAME") != "" {
		username = os.Getenv("DADCORP_USERNAME")
	}
	password := d.Get("password").(string)
	if os.Getenv("DADCORP_PASSWORD") != "" {
		password = os.Getenv("DADCORP_PASSWORD")
	}
	return &clientFactory{
		username: username,
		password: password,
	}, nil
}
