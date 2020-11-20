package provider

import (
	"context"

	dadcorp "dadcorp.dev/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConsulCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConsulClusterCreate,
		ReadContext:   resourceConsulClusterRead,
		UpdateContext: resourceConsulClusterUpdate,
		DeleteContext: resourceConsulClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bind_addr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"addresses": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"http": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"https": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"grpc": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"ports": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"http": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"https": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"grpc": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"serf_lan": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"serf_wan": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"server": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"sidecar_min_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"sidecar_max_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"expose_min_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"expose_max_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceConsulClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	cluster := dadcorp.ConsulCluster{
		Name:     d.Get("name").(string),
		BindAddr: d.Get("bind_addr").(string),
	}
	addresses := d.Get("addresses").([]interface{})
	if len(addresses) > 0 {
		cluster.Addresses = dadcorp.ConsulClusterAddresses{
			DNS:   addresses[0].(map[string]interface{})["dns"].(string),
			HTTP:  addresses[0].(map[string]interface{})["http"].(string),
			HTTPS: addresses[0].(map[string]interface{})["https"].(string),
			GRPC:  addresses[0].(map[string]interface{})["grpc"].(string),
		}
	}
	ports := d.Get("ports").([]interface{})
	if len(ports) > 0 {
		sidecarMin := ports[0].(map[string]interface{})["sidecar_min_port"].(int)
		sidecarMax := ports[0].(map[string]interface{})["sidecar_max_port"].(int)
		exposeMin := ports[0].(map[string]interface{})["expose_min_port"].(int)
		exposeMax := ports[0].(map[string]interface{})["expose_max_port"].(int)
		cluster.Ports = dadcorp.ConsulClusterPorts{
			DNS:            ports[0].(map[string]interface{})["dns"].(int),
			HTTP:           ports[0].(map[string]interface{})["http"].(int),
			HTTPS:          ports[0].(map[string]interface{})["https"].(int),
			GRPC:           ports[0].(map[string]interface{})["grpc"].(int),
			SerfLAN:        ports[0].(map[string]interface{})["serf_lan"].(int),
			SerfWAN:        ports[0].(map[string]interface{})["serf_wan"].(int),
			Server:         ports[0].(map[string]interface{})["server"].(int),
			SidecarMinPort: &sidecarMin,
			SidecarMaxPort: &sidecarMax,
			ExposeMinPort:  &exposeMin,
			ExposeMaxPort:  &exposeMax,
		}
	}
	resp, err := client.Consul.Clusters.Create(ctx, cluster)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.ID)
	d.Set("name", resp.Name)
	d.Set("bind_addr", resp.BindAddr)
	d.Set("addresses", []map[string]interface{}{
		{
			"dns":   resp.Addresses.DNS,
			"http":  resp.Addresses.HTTP,
			"https": resp.Addresses.HTTPS,
			"grpc":  resp.Addresses.GRPC,
		},
	})
	d.Set("ports", []map[string]interface{}{
		{
			"dns":              resp.Ports.DNS,
			"http":             resp.Ports.HTTP,
			"https":            resp.Ports.HTTPS,
			"grpc":             resp.Ports.GRPC,
			"serf_lan":         resp.Ports.SerfLAN,
			"serf_wan":         resp.Ports.SerfWAN,
			"server":           resp.Ports.Server,
			"sidecar_min_port": resp.Ports.SidecarMinPort,
			"sidecar_max_port": resp.Ports.SidecarMaxPort,
			"expose_min_port":  resp.Ports.ExposeMinPort,
			"expose_max_port":  resp.Ports.ExposeMaxPort,
		},
	})
	return nil
}

func resourceConsulClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Consul.Clusters.Get(ctx, d.Id())
	if err != nil {
		if err == dadcorp.ErrConsulClusterNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	d.Set("name", resp.Name)
	d.Set("bind_addr", resp.BindAddr)
	d.Set("addresses", []map[string]interface{}{
		{
			"dns":   resp.Addresses.DNS,
			"http":  resp.Addresses.HTTP,
			"https": resp.Addresses.HTTPS,
			"grpc":  resp.Addresses.GRPC,
		},
	})
	d.Set("ports", []map[string]interface{}{
		{
			"dns":              resp.Ports.DNS,
			"http":             resp.Ports.HTTP,
			"https":            resp.Ports.HTTPS,
			"grpc":             resp.Ports.GRPC,
			"serf_lan":         resp.Ports.SerfLAN,
			"serf_wan":         resp.Ports.SerfWAN,
			"server":           resp.Ports.Server,
			"sidecar_min_port": resp.Ports.SidecarMinPort,
			"sidecar_max_port": resp.Ports.SidecarMaxPort,
			"expose_min_port":  resp.Ports.ExposeMinPort,
			"expose_max_port":  resp.Ports.ExposeMaxPort,
		},
	})
	return nil
}

func resourceConsulClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	cluster := dadcorp.ConsulCluster{
		ID:       d.Id(),
		Name:     d.Get("name").(string),
		BindAddr: d.Get("bind_addr").(string),
	}
	addresses := d.Get("addresses").([]interface{})
	if len(addresses) > 0 {
		cluster.Addresses = dadcorp.ConsulClusterAddresses{
			DNS:   addresses[0].(map[string]interface{})["dns"].(string),
			HTTP:  addresses[0].(map[string]interface{})["http"].(string),
			HTTPS: addresses[0].(map[string]interface{})["https"].(string),
			GRPC:  addresses[0].(map[string]interface{})["grpc"].(string),
		}
	}
	ports := d.Get("ports").([]interface{})
	if len(ports) > 0 {
		sidecarMin := ports[0].(map[string]interface{})["sidecar_min_port"].(int)
		sidecarMax := ports[0].(map[string]interface{})["sidecar_max_port"].(int)
		exposeMin := ports[0].(map[string]interface{})["expose_min_port"].(int)
		exposeMax := ports[0].(map[string]interface{})["expose_max_port"].(int)
		cluster.Ports = dadcorp.ConsulClusterPorts{
			DNS:            ports[0].(map[string]interface{})["dns"].(int),
			HTTP:           ports[0].(map[string]interface{})["http"].(int),
			HTTPS:          ports[0].(map[string]interface{})["https"].(int),
			GRPC:           ports[0].(map[string]interface{})["grpc"].(int),
			SerfLAN:        ports[0].(map[string]interface{})["serf_lan"].(int),
			SerfWAN:        ports[0].(map[string]interface{})["serf_wan"].(int),
			Server:         ports[0].(map[string]interface{})["server"].(int),
			SidecarMinPort: &sidecarMin,
			SidecarMaxPort: &sidecarMax,
			ExposeMinPort:  &exposeMin,
			ExposeMaxPort:  &exposeMax,
		}
	}
	resp, err := client.Consul.Clusters.Update(ctx, cluster)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", resp.Name)
	d.Set("bind_addr", resp.BindAddr)
	d.Set("addresses", []map[string]interface{}{
		{
			"dns":   resp.Addresses.DNS,
			"http":  resp.Addresses.HTTP,
			"https": resp.Addresses.HTTPS,
			"grpc":  resp.Addresses.GRPC,
		},
	})
	d.Set("ports", []map[string]interface{}{
		{
			"dns":              resp.Ports.DNS,
			"http":             resp.Ports.HTTP,
			"https":            resp.Ports.HTTPS,
			"grpc":             resp.Ports.GRPC,
			"serf_lan":         resp.Ports.SerfLAN,
			"serf_wan":         resp.Ports.SerfWAN,
			"server":           resp.Ports.Server,
			"sidecar_min_port": resp.Ports.SidecarMinPort,
			"sidecar_max_port": resp.Ports.SidecarMaxPort,
			"expose_min_port":  resp.Ports.ExposeMinPort,
			"expose_max_port":  resp.Ports.ExposeMaxPort,
		},
	})
	return nil
}

func resourceConsulClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.Consul.Clusters.Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
