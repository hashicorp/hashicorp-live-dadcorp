package provider

import (
	"context"

	dadcorp "dadcorp.dev/client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVaultCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVaultClusterCreate,
		ReadContext:   resourceVaultClusterRead,
		UpdateContext: resourceVaultClusterUpdate,
		DeleteContext: resourceVaultClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					str, ok := i.(string)
					if !ok {
						return diag.Diagnostics{
							{
								Severity:      diag.Error,
								Summary:       "Invalid type",
								Detail:        "Value must be a string.",
								AttributePath: path,
							},
						}
					}
					for _, region := range []string{
						"us-va-1", "us-va-2",
					} {
						if region == str {
							return nil
						}
					}
					return diag.Diagnostics{
						{
							Severity:      diag.Error,
							Summary:       "Invalid region",
							Detail:        `Region must be one of "us-va-1" or "us-va-2".`,
							AttributePath: path,
						},
					}
				},
			},
			"default_lease_ttl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_lease_ttl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tcp_listener": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cluster_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceVaultClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	listener := d.Get("tcp_listener").([]interface{})
	cluster := dadcorp.VaultCluster{
		Name:            d.Get("name").(string),
		Region:          d.Get("region").(string),
		DefaultLeaseTTL: d.Get("default_lease_ttl").(string),
		MaxLeaseTTL:     d.Get("max_lease_ttl").(string),
	}
	if len(listener) > 0 {
		cluster.TCPListener = dadcorp.VaultClusterTCPListener{
			Address:        listener[0].(map[string]interface{})["address"].(string),
			ClusterAddress: listener[0].(map[string]interface{})["cluster_address"].(string),
		}
	}
	resp, err := client.Vault.Clusters.Create(ctx, cluster)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.ID)
	d.Set("name", resp.Name)
	d.Set("region", resp.Region)
	d.Set("default_lease_ttl", resp.DefaultLeaseTTL)
	d.Set("max_lease_ttl", resp.MaxLeaseTTL)
	d.Set("tcp_listener", []interface{}{
		map[string]interface{}{
			"address":         resp.TCPListener.Address,
			"cluster_address": resp.TCPListener.ClusterAddress,
		},
	})
	return nil
}

func resourceVaultClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Vault.Clusters.Get(ctx, d.Id())
	if err != nil {
		if err == dadcorp.ErrVaultClusterNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	d.Set("name", resp.Name)
	d.Set("region", resp.Region)
	d.Set("default_lease_ttl", resp.DefaultLeaseTTL)
	d.Set("max_lease_ttl", resp.MaxLeaseTTL)
	d.Set("tcp_listener", []interface{}{
		map[string]interface{}{
			"address":         resp.TCPListener.Address,
			"cluster_address": resp.TCPListener.ClusterAddress,
		},
	})
	return nil
}

func resourceVaultClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	listener := d.Get("tcp_listener").([]interface{})
	cluster := dadcorp.VaultCluster{
		ID:              d.Id(),
		Name:            d.Get("name").(string),
		Region:          d.Get("region").(string),
		DefaultLeaseTTL: d.Get("default_lease_ttl").(string),
		MaxLeaseTTL:     d.Get("max_lease_ttl").(string),
	}
	if len(listener) > 0 {
		cluster.TCPListener = dadcorp.VaultClusterTCPListener{
			Address:        listener[0].(map[string]interface{})["address"].(string),
			ClusterAddress: listener[0].(map[string]interface{})["cluster_address"].(string),
		}
	}
	resp, err := client.Vault.Clusters.Update(ctx, cluster)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", resp.Name)
	d.Set("region", resp.Region)
	d.Set("default_lease_ttl", resp.DefaultLeaseTTL)
	d.Set("max_lease_ttl", resp.MaxLeaseTTL)
	d.Set("tcp_listener", []interface{}{
		map[string]interface{}{
			"address":         resp.TCPListener.Address,
			"cluster_address": resp.TCPListener.ClusterAddress,
		},
	})
	return nil
}

func resourceVaultClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.Vault.Clusters.Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
