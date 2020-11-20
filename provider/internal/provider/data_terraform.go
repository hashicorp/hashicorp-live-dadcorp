package provider

import (
	"context"
	"encoding/json"

	dadcorp "dadcorp.dev/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAccessPolicyTerraform() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAccessPolicyTerraformRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plan": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"apply": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"override_policies": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataAccessPolicyTerraformRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policy := dadcorp.TerraformPolicy{
		WorkspaceID:      d.Get("workspace_id").(string),
		Plan:             d.Get("plan").(bool),
		Apply:            d.Get("apply").(bool),
		OverridePolicies: d.Get("override_policies").(bool),
	}
	b, err := json.Marshal(policy)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policy.WorkspaceID)
	d.Set("json", string(b))
	return nil
}
