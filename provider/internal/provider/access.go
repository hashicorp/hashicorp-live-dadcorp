package provider

import (
	"context"
	"encoding/json"
	"errors"

	dadcorp "dadcorp.dev/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAccessPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessPolicyCreate,
		ReadContext:   resourceAccessPolicyRead,
		UpdateContext: resourceAccessPolicyUpdate,
		DeleteContext: resourceAccessPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"policy_data": {
				Type:     schema.TypeString,
				Required: true,
				// TODO: diffsuppress
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				// TODO: validate
			},
		},
	}
}

func resourceAccessPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	policy := dadcorp.AccessPolicy{
		Type: d.Get("type").(string),
	}
	switch policy.Type {
	case "terraform":
		var data dadcorp.TerraformPolicy
		err = json.Unmarshal([]byte(d.Get("policy_data").(string)), &data)
		if err != nil {
			return diag.FromErr(err)
		}
		policy.PolicyData = data
	case "vault":
		var data dadcorp.VaultPolicy
		err = json.Unmarshal([]byte(d.Get("policy_data").(string)), &data)
		if err != nil {
			return diag.FromErr(err)
		}
		policy.PolicyData = data
	case "nomad":
		var data dadcorp.NomadPolicy
		err = json.Unmarshal([]byte(d.Get("policy_data").(string)), &data)
		if err != nil {
			return diag.FromErr(err)
		}
		policy.PolicyData = data
	case "consul":
		var data dadcorp.ConsulPolicy
		err = json.Unmarshal([]byte(d.Get("policy_data").(string)), &data)
		if err != nil {
			return diag.FromErr(err)
		}
		policy.PolicyData = data
	default:
		return diag.FromErr(errors.New("Access policies must be of type consul, nomad, terraform, or vault"))
	}
	resp, err := client.AccessPolicies.Create(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	policyData, err := json.Marshal(resp.PolicyData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.ID)
	d.Set("type", resp.Type)
	d.Set("policy_data", policyData)
	return nil
}

func resourceAccessPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AccessPolicies.Get(ctx, d.Id())
	if err != nil {
		if err == dadcorp.ErrAccessPolicyNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	policyData, err := json.Marshal(resp.PolicyData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("type", resp.Type)
	d.Set("policy_data", policyData)
	return nil
}

func resourceAccessPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	policy := dadcorp.AccessPolicy{
		ID:   d.Id(),
		Type: d.Get("type").(string),
	}
	switch policy.Type {
	case "terraform":
		var data dadcorp.TerraformPolicy
		err = json.Unmarshal([]byte(d.Get("policy_data").(string)), &data)
		if err != nil {
			return diag.FromErr(err)
		}
		policy.PolicyData = data
	case "vault":
		var data dadcorp.VaultPolicy
		err = json.Unmarshal([]byte(d.Get("policy_data").(string)), &data)
		if err != nil {
			return diag.FromErr(err)
		}
		policy.PolicyData = data
	case "nomad":
		var data dadcorp.NomadPolicy
		err = json.Unmarshal([]byte(d.Get("policy_data").(string)), &data)
		if err != nil {
			return diag.FromErr(err)
		}
		policy.PolicyData = data
	case "consul":
		var data dadcorp.ConsulPolicy
		err = json.Unmarshal([]byte(d.Get("policy_data").(string)), &data)
		if err != nil {
			return diag.FromErr(err)
		}
		policy.PolicyData = data
	default:
		return diag.FromErr(errors.New("Access policies must be of type consul, nomad, terraform, or vault"))
	}
	resp, err := client.AccessPolicies.Update(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	policyData, err := json.Marshal(resp.PolicyData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("type", resp.Type)
	d.Set("policy_data", policyData)
	return nil
}

func resourceAccessPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.AccessPolicies.Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
