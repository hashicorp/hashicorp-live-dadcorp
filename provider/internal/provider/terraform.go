package provider

import (
	"context"

	dadcorp "dadcorp.dev/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTerraformWorkspace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTerraformWorkspaceCreate,
		ReadContext:   resourceTerraformWorkspaceRead,
		UpdateContext: resourceTerraformWorkspaceUpdate,
		DeleteContext: resourceTerraformWorkspaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agent_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"allow_destroy_plan": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auto_apply": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execution_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				// TODO: validate
			},
			"file_triggers_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"queue_all_runs": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"speculative_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"terraform_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"trigger_prefixes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"working_directory": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vcs_repo": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oauth_token_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"branch": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ingress_submodules": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"identifier": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceTerraformWorkspaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	allowDestroyPlan := d.Get("allow_destroy_plan").(bool)
	fileTriggersEnabled := d.Get("file_triggers_enabled").(bool)
	speculativeEnabled := d.Get("speculative_enabled").(bool)
	workspace := dadcorp.TerraformWorkspace{
		Name:                d.Get("name").(string),
		AgentPoolID:         d.Get("agent_pool_id").(string),
		AllowDestroyPlan:    &allowDestroyPlan,
		AutoApply:           d.Get("auto_apply").(bool),
		Description:         d.Get("description").(string),
		ExecutionMode:       d.Get("execution_mode").(string),
		FileTriggersEnabled: &fileTriggersEnabled,
		QueueAllRuns:        d.Get("queue_all_runs").(bool),
		SpeculativeEnabled:  &speculativeEnabled,
		TerraformVersion:    d.Get("terraform_version").(string),
		WorkingDirectory:    d.Get("working_directory").(string),
	}
	for _, prefix := range d.Get("trigger_prefixes").([]interface{}) {
		workspace.TriggerPrefixes = append(workspace.TriggerPrefixes, prefix.(string))
	}
	vcs := d.Get("vcs_repo").([]interface{})
	if len(vcs) > 0 {
		workspace.VCSRepo = dadcorp.TerraformWorkspaceVCSRepo{
			OAuthTokenID:      vcs[0].(map[string]interface{})["oauth_token_id"].(string),
			Branch:            vcs[0].(map[string]interface{})["branch"].(string),
			IngressSubmodules: vcs[0].(map[string]interface{})["ingress_submodules"].(bool),
			Identifier:        vcs[0].(map[string]interface{})["identifier"].(string),
		}
	}
	resp, err := client.Terraform.Workspaces.Create(ctx, workspace)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.ID)
	d.Set("name", resp.Name)
	d.Set("agent_pool_id", resp.AgentPoolID)
	d.Set("allow_destroy_plan", resp.AllowDestroyPlan)
	d.Set("auto_apply", resp.AutoApply)
	d.Set("description", resp.Description)
	d.Set("execution_mode", resp.ExecutionMode)
	d.Set("file_triggers_enabled", resp.FileTriggersEnabled)
	d.Set("queue_all_runs", resp.QueueAllRuns)
	d.Set("speculative_enabled", resp.SpeculativeEnabled)
	d.Set("terraform_version", resp.TerraformVersion)
	d.Set("working_directory", resp.WorkingDirectory)
	d.Set("vsc_repo", []map[string]interface{}{
		{
			"oauth_token_id":     resp.VCSRepo.OAuthTokenID,
			"branch":             resp.VCSRepo.Branch,
			"ingress_submodules": resp.VCSRepo.IngressSubmodules,
			"identifier":         resp.VCSRepo.Identifier,
		},
	})
	return nil
}

func resourceTerraformWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Terraform.Workspaces.Get(ctx, d.Id())
	if err != nil {
		if err == dadcorp.ErrTerraformWorkspaceNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	d.Set("name", resp.Name)
	d.Set("agent_pool_id", resp.AgentPoolID)
	d.Set("allow_destroy_plan", resp.AllowDestroyPlan)
	d.Set("auto_apply", resp.AutoApply)
	d.Set("description", resp.Description)
	d.Set("execution_mode", resp.ExecutionMode)
	d.Set("file_triggers_enabled", resp.FileTriggersEnabled)
	d.Set("queue_all_runs", resp.QueueAllRuns)
	d.Set("speculative_enabled", resp.SpeculativeEnabled)
	d.Set("terraform_version", resp.TerraformVersion)
	d.Set("working_directory", resp.WorkingDirectory)
	d.Set("vsc_repo", []map[string]interface{}{
		{
			"oauth_token_id":     resp.VCSRepo.OAuthTokenID,
			"branch":             resp.VCSRepo.Branch,
			"ingress_submodules": resp.VCSRepo.IngressSubmodules,
			"identifier":         resp.VCSRepo.Identifier,
		},
	})
	return nil
}

func resourceTerraformWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	allowDestroyPlan := d.Get("allow_destroy_plan").(bool)
	fileTriggersEnabled := d.Get("file_triggers_enabled").(bool)
	speculativeEnabled := d.Get("speculative_enabled").(bool)
	workspace := dadcorp.TerraformWorkspace{
		Name:                d.Get("name").(string),
		AgentPoolID:         d.Get("agent_pool_id").(string),
		AllowDestroyPlan:    &allowDestroyPlan,
		AutoApply:           d.Get("auto_apply").(bool),
		Description:         d.Get("description").(string),
		ExecutionMode:       d.Get("execution_mode").(string),
		FileTriggersEnabled: &fileTriggersEnabled,
		QueueAllRuns:        d.Get("queue_all_runs").(bool),
		SpeculativeEnabled:  &speculativeEnabled,
		TerraformVersion:    d.Get("terraform_version").(string),
		WorkingDirectory:    d.Get("working_directory").(string),
	}
	for _, prefix := range d.Get("trigger_prefixes").([]interface{}) {
		workspace.TriggerPrefixes = append(workspace.TriggerPrefixes, prefix.(string))
	}
	vcs := d.Get("vcs_repo").([]interface{})
	if len(vcs) > 0 {
		workspace.VCSRepo = dadcorp.TerraformWorkspaceVCSRepo{
			OAuthTokenID:      vcs[0].(map[string]interface{})["oauth_token_id"].(string),
			Branch:            vcs[0].(map[string]interface{})["branch"].(string),
			IngressSubmodules: vcs[0].(map[string]interface{})["ingress_submodules"].(bool),
			Identifier:        vcs[0].(map[string]interface{})["identifier"].(string),
		}
	}
	resp, err := client.Terraform.Workspaces.Update(ctx, workspace)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", resp.Name)
	d.Set("agent_pool_id", resp.AgentPoolID)
	d.Set("allow_destroy_plan", resp.AllowDestroyPlan)
	d.Set("auto_apply", resp.AutoApply)
	d.Set("description", resp.Description)
	d.Set("execution_mode", resp.ExecutionMode)
	d.Set("file_triggers_enabled", resp.FileTriggersEnabled)
	d.Set("queue_all_runs", resp.QueueAllRuns)
	d.Set("speculative_enabled", resp.SpeculativeEnabled)
	d.Set("terraform_version", resp.TerraformVersion)
	d.Set("working_directory", resp.WorkingDirectory)
	d.Set("vsc_repo", []map[string]interface{}{
		{
			"oauth_token_id":     resp.VCSRepo.OAuthTokenID,
			"branch":             resp.VCSRepo.Branch,
			"ingress_submodules": resp.VCSRepo.IngressSubmodules,
			"identifier":         resp.VCSRepo.Identifier,
		},
	})
	return nil
}

func resourceTerraformWorkspaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*clientFactory).NewClient()
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.Terraform.Workspaces.Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
