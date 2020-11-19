package provider

import (
	"context"

	dadcorp "dadcorp.dev/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNomadCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNomadClusterCreate,
		ReadContext:   resourceNomadClusterRead,
		UpdateContext: resourceNomadClusterUpdate,
		DeleteContext: resourceNomadClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datacenter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bind_addr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"advertise": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"rpc": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"serf": {
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
						"http": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"rpc": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"serf": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"server": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_join": {
							Type:     schema.TypeList,
							MaxItems: 1,
							MinItems: 1,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retry_join": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"start_join": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"retry_max": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"retry_interval": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceNomadClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	cluster := dadcorp.NomadCluster{
		Name:       d.Get("name").(string),
		Datacenter: d.Get("datacenter").(string),
		BindAddr:   d.Get("bind_addr").(string),
	}
	advertise := d.Get("advertise").([]interface{})
	if len(advertise) > 0 {
		cluster.Advertise = dadcorp.NomadClusterAdvertise{
			HTTP: advertise[0].(map[string]interface{})["http"].(string),
			RPC:  advertise[0].(map[string]interface{})["rpc"].(string),
			Serf: advertise[0].(map[string]interface{})["serf"].(string),
		}
	}
	ports := d.Get("ports").([]interface{})
	if len(ports) > 0 {
		cluster.Ports = dadcorp.NomadClusterPorts{
			HTTP: ports[0].(map[string]interface{})["http"].(int),
			RPC:  ports[0].(map[string]interface{})["rpc"].(int),
			Serf: ports[0].(map[string]interface{})["serf"].(int),
		}
	}
	server := d.Get("server").([]interface{})
	if len(server) > 0 {
		serverJoin := server[0].(map[string]interface{})["server_join"].([]interface{})
		if len(serverJoin) > 0 {
			retryJoin := make([]string, 0, len(serverJoin[0].(map[string]interface{})["retry_join"].([]interface{})))
			for _, rj := range serverJoin[0].(map[string]interface{})["retry_join"].([]interface{}) {
				retryJoin = append(retryJoin, rj.(string))
			}
			startJoin := make([]string, 0, len(serverJoin[0].(map[string]interface{})["start_join"].([]interface{})))
			for _, sj := range serverJoin[0].(map[string]interface{})["start_join"].([]interface{}) {
				startJoin = append(startJoin, sj.(string))
			}
			cluster.Server = dadcorp.NomadClusterServer{
				ServerJoin: dadcorp.NomadClusterServerServerJoin{
					RetryJoin:     retryJoin,
					StartJoin:     startJoin,
					RetryMax:      serverJoin[0].(map[string]interface{})["retry_max"].(int),
					RetryInterval: serverJoin[0].(map[string]interface{})["retry_interval"].(string),
				},
			}
		}
	}
	resp, err := client.Nomad.Clusters.Create(ctx, cluster)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.ID)
	d.Set("name", resp.Name)
	d.Set("bind_addr", resp.BindAddr)
	d.Set("datacenter", resp.Datacenter)
	d.Set("advertise", []map[string]interface{}{
		{
			"http": resp.Advertise.HTTP,
			"rpc":  resp.Advertise.RPC,
			"serf": resp.Advertise.Serf,
		},
	})
	d.Set("ports", []map[string]interface{}{
		{
			"http": resp.Ports.HTTP,
			"rpc":  resp.Ports.RPC,
			"serf": resp.Ports.Serf,
		},
	})
	d.Set("server", []map[string]interface{}{
		{
			"server_join": []map[string]interface{}{
				{
					"retry_join":     resp.Server.ServerJoin.RetryJoin,
					"start_join":     resp.Server.ServerJoin.StartJoin,
					"retry_max":      resp.Server.ServerJoin.RetryMax,
					"retry_interval": resp.Server.ServerJoin.RetryInterval,
				},
			},
		},
	})
	return nil
}

func resourceNomadClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Nomad.Clusters.Get(ctx, d.Id())
	if err != nil {
		if err == dadcorp.ErrNomadClusterNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	d.Set("name", resp.Name)
	d.Set("bind_addr", resp.BindAddr)
	d.Set("datacenter", resp.Datacenter)
	d.Set("advertise", []map[string]interface{}{
		{
			"http": resp.Advertise.HTTP,
			"rpc":  resp.Advertise.RPC,
			"serf": resp.Advertise.Serf,
		},
	})
	d.Set("ports", []map[string]interface{}{
		{
			"http": resp.Ports.HTTP,
			"rpc":  resp.Ports.RPC,
			"serf": resp.Ports.Serf,
		},
	})
	d.Set("server", []map[string]interface{}{
		{
			"server_join": []map[string]interface{}{
				{
					"retry_join":     resp.Server.ServerJoin.RetryJoin,
					"start_join":     resp.Server.ServerJoin.StartJoin,
					"retry_max":      resp.Server.ServerJoin.RetryMax,
					"retry_interval": resp.Server.ServerJoin.RetryInterval,
				},
			},
		},
	})
	return nil
}

func resourceNomadClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	cluster := dadcorp.NomadCluster{
		ID:         d.Id(),
		Name:       d.Get("name").(string),
		Datacenter: d.Get("datacenter").(string),
		BindAddr:   d.Get("bind_addr").(string),
	}
	advertise := d.Get("advertise").([]interface{})
	if len(advertise) > 0 {
		cluster.Advertise = dadcorp.NomadClusterAdvertise{
			HTTP: advertise[0].(map[string]interface{})["http"].(string),
			RPC:  advertise[0].(map[string]interface{})["rpc"].(string),
			Serf: advertise[0].(map[string]interface{})["serf"].(string),
		}
	}
	ports := d.Get("ports").([]interface{})
	if len(ports) > 0 {
		cluster.Ports = dadcorp.NomadClusterPorts{
			HTTP: ports[0].(map[string]interface{})["http"].(int),
			RPC:  ports[0].(map[string]interface{})["rpc"].(int),
			Serf: ports[0].(map[string]interface{})["serf"].(int),
		}
	}
	server := d.Get("server").([]interface{})
	if len(server) > 0 {
		serverJoin := server[0].(map[string]interface{})["server_join"].([]interface{})
		if len(serverJoin) > 0 {
			retryJoin := make([]string, 0, len(serverJoin[0].(map[string]interface{})["retry_join"].([]interface{})))
			for _, rj := range serverJoin[0].(map[string]interface{})["retry_join"].([]interface{}) {
				retryJoin = append(retryJoin, rj.(string))
			}
			startJoin := make([]string, 0, len(serverJoin[0].(map[string]interface{})["start_join"].([]interface{})))
			for _, sj := range serverJoin[0].(map[string]interface{})["start_join"].([]interface{}) {
				startJoin = append(startJoin, sj.(string))
			}
			cluster.Server = dadcorp.NomadClusterServer{
				ServerJoin: dadcorp.NomadClusterServerServerJoin{
					RetryJoin:     retryJoin,
					StartJoin:     startJoin,
					RetryMax:      serverJoin[0].(map[string]interface{})["retry_max"].(int),
					RetryInterval: serverJoin[0].(map[string]interface{})["retry_interval"].(string),
				},
			}
		}
	}
	resp, err := client.Nomad.Clusters.Update(ctx, cluster)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", resp.Name)
	d.Set("bind_addr", resp.BindAddr)
	d.Set("datacenter", resp.Datacenter)
	d.Set("advertise", []map[string]interface{}{
		{
			"http": resp.Advertise.HTTP,
			"rpc":  resp.Advertise.RPC,
			"serf": resp.Advertise.Serf,
		},
	})
	d.Set("ports", []map[string]interface{}{
		{
			"http": resp.Ports.HTTP,
			"rpc":  resp.Ports.RPC,
			"serf": resp.Ports.Serf,
		},
	})
	d.Set("server", []map[string]interface{}{
		{
			"server_join": []map[string]interface{}{
				{
					"retry_join":     resp.Server.ServerJoin.RetryJoin,
					"start_join":     resp.Server.ServerJoin.StartJoin,
					"retry_max":      resp.Server.ServerJoin.RetryMax,
					"retry_interval": resp.Server.ServerJoin.RetryInterval,
				},
			},
		},
	})
	return nil
}

func resourceNomadClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.Nomad.Clusters.Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
