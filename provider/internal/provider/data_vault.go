package provider

import (
	"context"
	"encoding/json"

	dadcorp "dadcorp.dev/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAccessPolicyVault() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAccessPolicyVaultRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"read": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"write": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"delete": {
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

func dataAccessPolicyVaultRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policy := dadcorp.VaultPolicy{
		ClusterID: d.Get("cluster_id").(string),
		Key:       d.Get("key").(string),
		Read:      d.Get("read").(bool),
		Write:     d.Get("write").(bool),
		Delete:    d.Get("delete").(bool),
	}
	b, err := json.Marshal(policy)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policy.ClusterID + "." + policy.Key)
	d.Set("json", string(b))
	return nil
}
