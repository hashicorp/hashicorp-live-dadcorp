package provider

import (
	"context"
	"encoding/json"

	dadcorp "dadcorp.dev/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAccessPolicyNomad() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAccessPolicyNomadRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"submit_jobs": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"read_job_status": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"cancel_jobs": {
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

func dataAccessPolicyNomadRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policy := dadcorp.NomadPolicy{
		ClusterID:     d.Get("cluster_id").(string),
		SubmitJobs:    d.Get("submit_jobs").(bool),
		ReadJobStatus: d.Get("read_job_status").(bool),
		CancelJobs:    d.Get("cancel_jobs").(bool),
	}
	b, err := json.Marshal(policy)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policy.ClusterID)
	d.Set("json", string(b))
	return nil
}
