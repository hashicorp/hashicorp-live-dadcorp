package provider

import (
	"context"

	dadcorp "dadcorp.dev/client"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tftypes"
)

type terraform struct {
	clients clientFactory
}

func (t *terraform) workspaceType() tftypes.Type {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":                    tftypes.String,
			"name":                  tftypes.String,
			"agent_pool_id":         tftypes.String,
			"allow_destroy_plan":    tftypes.Bool,
			"auto_apply":            tftypes.Bool,
			"description":           tftypes.String,
			"execution_mode":        tftypes.String,
			"file_triggers_enabled": tftypes.Bool,
			"queue_all_runs":        tftypes.Bool,
			"speculative_enabled":   tftypes.Bool,
			"terraform_version":     tftypes.String,
			"trigger_prefixes": tftypes.List{
				ElementType: tftypes.String,
			},
			"working_directory": tftypes.String,
			"vcs_repo":          t.vcsRepoType(),
		},
	}
}

func (t *terraform) vcsRepoType() tftypes.Type {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"oauth_token_id":     tftypes.String,
			"branch":             tftypes.String,
			"ingress_submodules": tftypes.Bool,
			"identifier":         tftypes.String,
		},
	}
}

func (t *terraform) schema() *tfprotov5.Schema {
	return &tfprotov5.Schema{
		Version: 1,
		Block: &tfprotov5.SchemaBlock{
			Attributes: []*tfprotov5.SchemaAttribute{
				{
					Name:     "id",
					Type:     tftypes.String,
					Computed: true,
				},
				{
					Name:     "name",
					Type:     tftypes.String,
					Required: true,
				},
				{
					Name:     "agent_pool_id",
					Type:     tftypes.String,
					Optional: true,
					Computed: true,
				},
				{
					Name:     "allow_destroy_plan",
					Type:     tftypes.Bool,
					Optional: true,
					Computed: true,
				},
				{
					Name:     "auto_apply",
					Type:     tftypes.Bool,
					Optional: true,
					Computed: true,
				},
				{
					Name:     "description",
					Type:     tftypes.String,
					Optional: true,
					Computed: true,
				},
				{
					Name:     "execution_mode",
					Type:     tftypes.String,
					Optional: true,
					Computed: true,
				},
				{
					Name:     "file_triggers_enabled",
					Type:     tftypes.Bool,
					Optional: true,
					Computed: true,
				},
				{
					Name:     "queue_all_runs",
					Type:     tftypes.Bool,
					Optional: true,
					Computed: true,
				},
				{
					Name:     "speculative_enabled",
					Type:     tftypes.Bool,
					Optional: true,
					Computed: true,
				},
				{
					Name:     "terraform_version",
					Type:     tftypes.String,
					Optional: true,
					Computed: true,
				},
				{
					Name:     "trigger_prefixes",
					Type:     tftypes.List{ElementType: tftypes.String},
					Optional: true,
					Computed: true,
				},
				{
					Name:     "working_directory",
					Type:     tftypes.String,
					Optional: true,
					Computed: true,
				},
			},
			BlockTypes: []*tfprotov5.SchemaNestedBlock{
				{
					TypeName: "vcs_repo",
					Nesting:  tfprotov5.SchemaNestedBlockNestingModeSingle,
					Block: &tfprotov5.SchemaBlock{
						Attributes: []*tfprotov5.SchemaAttribute{
							{
								Name:     "oauth_token_id",
								Type:     tftypes.String,
								Optional: true,
							},
							{
								Name:     "branch",
								Type:     tftypes.String,
								Optional: true,
								Computed: true,
							},
							{
								Name:     "ingress_submodules",
								Type:     tftypes.Bool,
								Optional: true,
								Computed: true,
							},
							{
								Name:     "identifier",
								Type:     tftypes.String,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
}

func (t *terraform) ValidateResourceTypeConfig(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	val, err := req.Config.Unmarshal(t.workspaceType())
	if err != nil {
		return &tfprotov5.ValidateResourceTypeConfigResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected configuration format",
					Detail:   "The resource got a configuration that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	if !val.Is(t.workspaceType()) {
		return &tfprotov5.ValidateResourceTypeConfigResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected configuration format",
					Detail:   "The resource got a configuration that did not match its schema. This may indicate an error in the provider.",
				},
			},
		}, nil
	}
	values := map[string]tftypes.Value{}
	err = val.As(&values)
	if err != nil {
		return &tfprotov5.ValidateResourceTypeConfigResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected configuration format",
					Detail:   "The resource got a configuration that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	if values["execution_mode"].IsKnown() && !values["execution_mode"].IsNull() && values["agent_pool_id"].IsKnown() && !values["agent_pool_id"].IsNull() {
		var execMode string
		err = values["execution_mode"].As(&execMode)
		if err != nil {
			return &tfprotov5.ValidateResourceTypeConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected configuration format",
						Detail:   "The resource got a configuration that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("execution_mode"),
							},
						},
					},
				},
			}, nil
		}
		var agentPoolID string
		err = values["agent_pool_id"].As(&agentPoolID)
		if err != nil {
			return &tfprotov5.ValidateResourceTypeConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected configuration format",
						Detail:   "The resource got a configuration that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("agent_pool_id"),
							},
						},
					},
				},
			}, nil
		}
		if execMode != "agent" && agentPoolID != "" {
			return &tfprotov5.ValidateResourceTypeConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Invalid execution mode",
						Detail:   "execution_mode cannot be set when agent_pool_id is set.",
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("execution_mode"),
							},
						},
					},
				},
			}, nil
		}
	}
	return &tfprotov5.ValidateResourceTypeConfigResponse{}, nil
}

func (t *terraform) UpgradeResourceState(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	switch req.Version {
	case 1:
		val, err := req.RawState.Unmarshal(t.workspaceType())
		if err != nil {
			return &tfprotov5.UpgradeResourceStateResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected configuration format",
						Detail:   "The resource got a configuration that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					},
				},
			}, nil
		}
		dv, err := tfprotov5.NewDynamicValue(t.workspaceType(), val)
		if err != nil {
			return &tfprotov5.UpgradeResourceStateResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected configuration format",
						Detail:   "The resource got a configuration that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					},
				},
			}, nil
		}
		return &tfprotov5.UpgradeResourceStateResponse{
			UpgradedState: &dv,
		}, nil
	default:
		return &tfprotov5.UpgradeResourceStateResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected state version",
					Detail:   "The provider doesn't know how to upgrade from the current state version. Try an earlier releae of the provider.",
				},
			},
		}, nil
	}
}

func (t *terraform) ReadResource(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	val, err := req.CurrentState.Unmarshal(t.workspaceType())
	if err != nil {
		return &tfprotov5.ReadResourceResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected state format",
					Detail:   "The resource got a state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	state := map[string]tftypes.Value{}
	err = val.As(&state)
	if err != nil {
		return &tfprotov5.ReadResourceResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected state format",
					Detail:   "The resource got a state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	var id string
	err = state["id"].As(&id)
	if err != nil {
		return &tfprotov5.ReadResourceResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected state format",
					Detail:   "The resource got a state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					Attribute: &tftypes.AttributePath{
						Steps: []tftypes.AttributePathStep{
							tftypes.AttributeName("id"),
						},
					},
				},
			},
		}, nil
	}
	client, err := t.clients.NewClient()
	if err != nil {
		return &tfprotov5.ReadResourceResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error creating client",
					Detail:   "The provider was unable to create a client.\n\nError:\n" + err.Error(),
				},
			},
		}, nil
	}
	workspace, err := client.Terraform.Workspaces.Get(ctx, id)
	if err != nil {
		if err == dadcorp.ErrTerraformWorkspaceNotFound {
			dv, err := tfprotov5.NewDynamicValue(t.workspaceType(), tftypes.NewValue(t.workspaceType(), nil))
			if err != nil {
				return &tfprotov5.ReadResourceResponse{
					Diagnostics: []*tfprotov5.Diagnostic{
						{
							Severity: tfprotov5.DiagnosticSeverityError,
							Summary:  "Error removing workspace from state",
							Detail:   "An unexpected error was encountered removing the workspace from state. This is an error with the provider.\n\nError: " + err.Error(),
						},
					},
				}, nil
			}
			return &tfprotov5.ReadResourceResponse{
				NewState: &dv,
			}, nil
		}
		return &tfprotov5.ReadResourceResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error retrieving workspace",
					Detail:   "The provider was unable to retrieve the workspace.\n\nError:\n" + err.Error(),
				},
			},
		}, nil
	}
	triggerPrefixes := make([]tftypes.Value, 0, len(workspace.TriggerPrefixes))
	for _, prefix := range workspace.TriggerPrefixes {
		triggerPrefixes = append(triggerPrefixes, tftypes.NewValue(tftypes.String, prefix))
	}
	dv, err := tfprotov5.NewDynamicValue(t.workspaceType(), tftypes.NewValue(t.workspaceType(), map[string]tftypes.Value{
		"id":                    tftypes.NewValue(tftypes.String, workspace.ID),
		"name":                  tftypes.NewValue(tftypes.String, workspace.Name),
		"agent_pool_id":         tftypes.NewValue(tftypes.String, workspace.AgentPoolID),
		"allow_destroy_plan":    tftypes.NewValue(tftypes.Bool, workspace.AllowDestroyPlan),
		"auto_apply":            tftypes.NewValue(tftypes.Bool, workspace.AutoApply),
		"description":           tftypes.NewValue(tftypes.String, workspace.Description),
		"execution_mode":        tftypes.NewValue(tftypes.String, workspace.ExecutionMode),
		"file_triggers_enabled": tftypes.NewValue(tftypes.Bool, workspace.FileTriggersEnabled),
		"queue_all_runs":        tftypes.NewValue(tftypes.Bool, workspace.QueueAllRuns),
		"speculative_enabled":   tftypes.NewValue(tftypes.Bool, workspace.SpeculativeEnabled),
		"terraform_version":     tftypes.NewValue(tftypes.String, workspace.TerraformVersion),
		"trigger_prefixes":      tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, triggerPrefixes),
		"working_directory":     tftypes.NewValue(tftypes.String, workspace.WorkingDirectory),
		"vcs_repo": tftypes.NewValue(t.vcsRepoType(), map[string]tftypes.Value{
			"oauth_token_id":     tftypes.NewValue(tftypes.String, workspace.VCSRepo.OAuthTokenID),
			"branch":             tftypes.NewValue(tftypes.String, workspace.VCSRepo.Branch),
			"identifier":         tftypes.NewValue(tftypes.String, workspace.VCSRepo.Identifier),
			"ingress_submodules": tftypes.NewValue(tftypes.Bool, workspace.VCSRepo.IngressSubmodules),
		}),
	}))
	if err != nil {
		return &tfprotov5.ReadResourceResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error updating workspace in state",
					Detail:   "An unexpected error was encountered updating the workspace from state. This is an error with the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	return &tfprotov5.ReadResourceResponse{
		NewState: &dv,
	}, nil
}

func (t *terraform) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	val, err := req.ProposedNewState.Unmarshal(t.workspaceType())
	if err != nil {
		return &tfprotov5.PlanResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected state format",
					Detail:   "The resource got a state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	newState := map[string]tftypes.Value{}
	err = val.As(&newState)
	if err != nil {
		return &tfprotov5.PlanResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected state format",
					Detail:   "The resource got a state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	val, err = req.PriorState.Unmarshal(t.workspaceType())
	if err != nil {
		return &tfprotov5.PlanResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected prior state format",
					Detail:   "The resource got a prior state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	oldState := map[string]tftypes.Value{}
	err = val.As(&oldState)
	if err != nil {
		return &tfprotov5.PlanResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected prior state format",
					Detail:   "The resource got a prior state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	if newState["id"].IsNull() {
		newState["id"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	}
	if newState["agent_pool_id"].IsNull() {
		var oldAgentPoolID string
		err = oldState["agent_pool_id"].As(&oldAgentPoolID)
		if err != nil {
			return &tfprotov5.PlanResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected prior state format",
						Detail:   "The resource got a prior state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("agent_pool_id"),
							},
						},
					},
				},
			}, nil
		}
		if oldAgentPoolID == "" {
			newState["agent_pool_id"] = tftypes.NewValue(tftypes.String, "")
		}
	}
	if newState["allow_destroy_plan"].IsNull() {
		newState["allow_destroy_plan"] = tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)
	}
	if newState["auto_apply"].IsNull() {
		newState["auto_apply"] = tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)
	}
	if newState["description"].IsNull() {
		var oldDescription string
		err = oldState["description"].As(&oldDescription)
		if err != nil {
			return &tfprotov5.PlanResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected prior state format",
						Detail:   "The resource got a prior state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("description"),
							},
						},
					},
				},
			}, nil
		}
		if oldDescription == "" {
			newState["description"] = tftypes.NewValue(tftypes.String, "")
		}
	}
	if newState["execution_mode"].IsNull() {
		newState["execution_mode"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	}
	if newState["file_triggers_enabled"].IsNull() {
		newState["file_triggers_enabled"] = tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)
	}
	if newState["speculative_enabled"].IsNull() {
		newState["speculative_enabled"] = tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)
	}
	if newState["terraform_version"].IsNull() {
		newState["terraform_version"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	}
	if newState["trigger_prefixes"].IsNull() {
		newState["trigger_prefixes"] = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, tftypes.UnknownValue)
	}
	if newState["working_directory"].IsNull() {
		newState["working_directory"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	}
	if !newState["vcs_repo"].IsNull() && newState["vcs_repo"].IsKnown() {
		vcs := map[string]tftypes.Value{}
		err = newState["vcs_repo"].As(&vcs)
		if err != nil {
			return &tfprotov5.PlanResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected state format",
						Detail:   "The resource got a state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("vcs_repo"),
							},
						},
					},
				},
			}, nil
		}
		oldVCS := map[string]tftypes.Value{}
		err = oldState["vcs_repo"].As(&oldVCS)
		if err != nil {
			return &tfprotov5.PlanResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected prior state format",
						Detail:   "The resource got a prior state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("vcs_repo"),
							},
						},
					},
				},
			}, nil
		}
		if vcs["ingress_submodules"].IsNull() {
			var ingress bool
			err = oldVCS["ingress_submodules"].As(&ingress)
			if err != nil {
				return &tfprotov5.PlanResourceChangeResponse{
					Diagnostics: []*tfprotov5.Diagnostic{
						{
							Severity: tfprotov5.DiagnosticSeverityError,
							Summary:  "Unexpected prior state format",
							Detail:   "The resource got a prior state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
							Attribute: &tftypes.AttributePath{
								Steps: []tftypes.AttributePathStep{
									tftypes.AttributeName("vcs_repo"),
									tftypes.AttributeName("ingress_submodules"),
								},
							},
						},
					},
				}, nil
			}
			if !ingress {
				vcs["ingress_submodules"] = tftypes.NewValue(tftypes.Bool, false)
			}
		}
		newVCS := map[string]tftypes.Value{
			"oauth_token_id":     vcs["oauth_token_id"],
			"branch":             tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"ingress_submodules": vcs["ingress_submodules"],
			"identifier":         vcs["identifier"],
		}
		if !vcs["branch"].IsNull() {
			newVCS["branch"] = vcs["branch"]
		}
		newState["vcs_repo"] = tftypes.NewValue(t.vcsRepoType(), newVCS)
	}
	dv, err := tfprotov5.NewDynamicValue(t.workspaceType(), tftypes.NewValue(t.workspaceType(), newState))
	if err != nil {
		return &tfprotov5.PlanResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error returning updated plan",
					Detail:   "The resource encountered an unexpected error returning the updated plan. This indicates an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	return &tfprotov5.PlanResourceChangeResponse{
		PlannedState: &dv,
	}, nil
}

func (t *terraform) ApplyResourceChange(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	plannedStateVal, err := req.PlannedState.Unmarshal(t.workspaceType())
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected planned state format",
					Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	priorStateVal, err := req.PriorState.Unmarshal(t.workspaceType())
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected prior state format",
					Detail:   "The resource got a prior state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	client, err := t.clients.NewClient()
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error creating client",
					Detail:   "The provider was unable to create a client.\n\nError:\n" + err.Error(),
				},
			},
		}, nil
	}

	// if plannedStateVal is null, we're deleting the workspace
	if plannedStateVal.IsNull() {
		priorState := map[string]tftypes.Value{}
		err = priorStateVal.As(&priorState)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected prior state format",
						Detail:   "The resource got a prior state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					},
				},
			}, nil
		}
		var id string
		err = priorState["id"].As(&id)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected prior state format",
						Detail:   "The resource got a prior state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("id"),
							},
						},
					},
				},
			}, nil
		}
		err = client.Terraform.Workspaces.Delete(ctx, id)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Error deleting workspace",
						Detail:   "The provider was unable to delete the workspace.\n\nError:\n" + err.Error(),
					},
				},
			}, nil
		}
		dv, err := tfprotov5.NewDynamicValue(t.workspaceType(), tftypes.NewValue(t.workspaceType(), nil))
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Error returning updated state",
						Detail:   "The resource encountered an unexpected error returning the updated state. This indicates an error in the provider.\n\nError: " + err.Error(),
					},
				},
			}, nil
		}
		return &tfprotov5.ApplyResourceChangeResponse{
			NewState: &dv,
		}, nil
	}

	// if plannedStateVal is not null, we're creating or updating the workspace
	// so let's get access to the planned state
	plannedState := map[string]tftypes.Value{}
	err = plannedStateVal.As(&plannedState)
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected planned state format",
					Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}

	var workspace dadcorp.TerraformWorkspace
	err = plannedState["name"].As(&workspace.Name)
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected planned state format",
					Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					Attribute: &tftypes.AttributePath{
						Steps: []tftypes.AttributePathStep{
							tftypes.AttributeName("name"),
						},
					},
				},
			},
		}, nil
	}
	err = plannedState["agent_pool_id"].As(&workspace.AgentPoolID)
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected planned state format",
					Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					Attribute: &tftypes.AttributePath{
						Steps: []tftypes.AttributePathStep{
							tftypes.AttributeName("agent_pool_id"),
						},
					},
				},
			},
		}, nil
	}
	if plannedState["allow_destroy_plan"].IsKnown() && !plannedState["allow_destroy_plan"].IsNull() {
		var adp bool
		err = plannedState["allow_destroy_plan"].As(&adp)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("allow_destroy_plan"),
							},
						},
					},
				},
			}, nil
		}
		workspace.AllowDestroyPlan = &adp
	}
	if plannedState["auto_apply"].IsKnown() && !plannedState["auto_apply"].IsNull() {
		err = plannedState["auto_apply"].As(&workspace.AutoApply)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("auto_apply"),
							},
						},
					},
				},
			}, nil
		}
	}
	if plannedState["description"].IsKnown() && !plannedState["description"].IsNull() {
		err = plannedState["description"].As(&workspace.Description)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("description"),
							},
						},
					},
				},
			}, nil
		}
	}
	if plannedState["execution_mode"].IsKnown() {
		err = plannedState["execution_mode"].As(&workspace.ExecutionMode)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("execution_mode"),
							},
						},
					},
				},
			}, nil
		}
	}
	if plannedState["file_triggers_enabled"].IsKnown() && !plannedState["file_triggers_enabled"].IsNull() {
		var fte bool
		err = plannedState["file_triggers_enabled"].As(&fte)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("file_triggers_enabled"),
							},
						},
					},
				},
			}, nil
		}
		workspace.FileTriggersEnabled = &fte
	}
	err = plannedState["queue_all_runs"].As(&workspace.QueueAllRuns)
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected planned state format",
					Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					Attribute: &tftypes.AttributePath{
						Steps: []tftypes.AttributePathStep{
							tftypes.AttributeName("queue_all_runs"),
						},
					},
				},
			},
		}, nil
	}
	if plannedState["speculative_enabled"].IsKnown() {
		var se bool
		err = plannedState["speculative_enabled"].As(&se)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("speculative_enabled"),
							},
						},
					},
				},
			}, nil
		}
		workspace.SpeculativeEnabled = &se
	}
	if plannedState["terraform_version"].IsKnown() {
		err = plannedState["terraform_version"].As(&workspace.TerraformVersion)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("terraform_version"),
							},
						},
					},
				},
			}, nil
		}
	}
	if plannedState["trigger_prefixes"].IsKnown() {
		var triggerPrefixes []tftypes.Value
		err = plannedState["trigger_prefixes"].As(&triggerPrefixes)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("trigger_prefixes"),
							},
						},
					},
				},
			}, nil
		}
		workspace.TriggerPrefixes = make([]string, 0, len(triggerPrefixes))
		for pos, prefix := range triggerPrefixes {
			var tp string
			err = prefix.As(&tp)
			if err != nil {
				return &tfprotov5.ApplyResourceChangeResponse{
					Diagnostics: []*tfprotov5.Diagnostic{
						{
							Severity: tfprotov5.DiagnosticSeverityError,
							Summary:  "Unexpected planned state format",
							Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
							Attribute: &tftypes.AttributePath{
								Steps: []tftypes.AttributePathStep{
									tftypes.AttributeName("trigger_prefixes"),
									tftypes.ElementKeyInt(pos),
								},
							},
						},
					},
				}, nil
			}
			workspace.TriggerPrefixes = append(workspace.TriggerPrefixes, tp)
		}
	}
	if plannedState["working_directory"].IsKnown() {
		err = plannedState["working_directory"].As(&workspace.WorkingDirectory)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("working_directory"),
							},
						},
					},
				},
			}, nil
		}
	}
	vcs := map[string]tftypes.Value{}
	err = plannedState["vcs_repo"].As(&vcs)
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected planned state format",
					Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					Attribute: &tftypes.AttributePath{
						Steps: []tftypes.AttributePathStep{
							tftypes.AttributeName("vcs_repo"),
						},
					},
				},
			},
		}, nil
	}
	if vcs["branch"].IsKnown() {
		err = vcs["branch"].As(&workspace.VCSRepo.Branch)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("vcs_repo"),
								tftypes.AttributeName("branch"),
							},
						},
					},
				},
			}, nil
		}
	}
	err = vcs["oauth_token_id"].As(&workspace.VCSRepo.OAuthTokenID)
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected planned state format",
					Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					Attribute: &tftypes.AttributePath{
						Steps: []tftypes.AttributePathStep{
							tftypes.AttributeName("vcs_repo"),
							tftypes.AttributeName("oauth_token_id"),
						},
					},
				},
			},
		}, nil
	}
	err = vcs["identifier"].As(&workspace.VCSRepo.Identifier)
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected planned state format",
					Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					Attribute: &tftypes.AttributePath{
						Steps: []tftypes.AttributePathStep{
							tftypes.AttributeName("vcs_repo"),
							tftypes.AttributeName("identifier"),
						},
					},
				},
			},
		}, nil
	}
	err = vcs["ingress_submodules"].As(&workspace.VCSRepo.IngressSubmodules)
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected planned state format",
					Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					Attribute: &tftypes.AttributePath{
						Steps: []tftypes.AttributePathStep{
							tftypes.AttributeName("vcs_repo"),
							tftypes.AttributeName("ingress_submodules"),
						},
					},
				},
			},
		}, nil
	}

	// if priorStateVal is not null, we're updating the workspace
	if !priorStateVal.IsNull() {
		priorState := map[string]tftypes.Value{}
		err = priorStateVal.As(&priorState)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected prior state format",
						Detail:   "The resource got a prior state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					},
				},
			}, nil
		}
		err = priorState["id"].As(&workspace.ID)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected prior state format",
						Detail:   "The resource got a prior state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("id"),
							},
						},
					},
				},
			}, nil
		}
		workspace, err = client.Terraform.Workspaces.Update(ctx, workspace)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Error updating the workspace",
						Detail:   "The provider was unable to update the workspace.\n\nError: " + err.Error(),
					},
				},
			}, nil
		}
	} else {
		// if priorStateVal is null, we're creating the workspace
		workspace, err = client.Terraform.Workspaces.Create(ctx, workspace)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Error creating the workspace",
						Detail:   "The provider was unable to create the workspace.\n\nError: " + err.Error(),
					},
				},
			}, nil
		}
	}

	finalState := map[string]tftypes.Value{
		"id":             tftypes.NewValue(tftypes.String, workspace.ID),
		"name":           plannedState["name"],
		"agent_pool_id":  plannedState["agent_pool_id"],
		"queue_all_runs": plannedState["queue_all_runs"],
	}
	if plannedState["allow_destroy_plan"].IsKnown() && !plannedState["allow_destroy_plan"].IsNull() {
		finalState["allow_destroy_plan"] = plannedState["allow_destroy_plan"]
	} else {
		finalState["allow_destroy_plan"] = tftypes.NewValue(tftypes.Bool, workspace.AllowDestroyPlan)
	}
	if plannedState["auto_apply"].IsKnown() && !plannedState["auto_apply"].IsNull() {
		finalState["auto_apply"] = plannedState["auto_apply"]
	} else {
		finalState["auto_apply"] = tftypes.NewValue(tftypes.Bool, workspace.AutoApply)
	}
	if plannedState["description"].IsKnown() && !plannedState["description"].IsNull() {
		finalState["description"] = plannedState["description"]
	} else {
		finalState["description"] = tftypes.NewValue(tftypes.String, workspace.Description)
	}
	if plannedState["execution_mode"].IsKnown() && !plannedState["execution_mode"].IsNull() {
		finalState["execution_mode"] = plannedState["execution_mode"]
	} else {
		finalState["execution_mode"] = tftypes.NewValue(tftypes.String, workspace.ExecutionMode)
	}
	if plannedState["file_triggers_enabled"].IsKnown() && !plannedState["file_triggers_enabled"].IsNull() {
		finalState["file_triggers_enabled"] = plannedState["file_triggers_enabled"]
	} else {
		finalState["file_triggers_enabled"] = tftypes.NewValue(tftypes.Bool, workspace.FileTriggersEnabled)
	}
	if plannedState["speculative_enabled"].IsKnown() && !plannedState["speculative_enabled"].IsNull() {
		finalState["speculative_enabled"] = plannedState["speculative_enabled"]
	} else {
		finalState["speculative_enabled"] = tftypes.NewValue(tftypes.Bool, workspace.SpeculativeEnabled)
	}
	if plannedState["terraform_version"].IsKnown() && !plannedState["terraform_version"].IsNull() {
		finalState["terraform_version"] = plannedState["terraform_version"]
	} else {
		finalState["terraform_version"] = tftypes.NewValue(tftypes.String, workspace.TerraformVersion)
	}
	if plannedState["trigger_prefixes"].IsKnown() && !plannedState["trigger_prefixes"].IsNull() {
		finalState["trigger_prefixes"] = plannedState["trigger_prefixes"]
	} else {
		prefixes := make([]tftypes.Value, 0, len(workspace.TriggerPrefixes))
		for _, prefix := range workspace.TriggerPrefixes {
			prefixes = append(prefixes, tftypes.NewValue(tftypes.String, prefix))
		}
		finalState["trigger_prefixes"] = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, prefixes)
	}
	if plannedState["working_directory"].IsKnown() && !plannedState["working_directory"].IsNull() {
		finalState["working_directory"] = plannedState["working_directory"]
	} else {
		finalState["working_directory"] = tftypes.NewValue(tftypes.String, workspace.WorkingDirectory)
	}
	finalVCS := map[string]tftypes.Value{
		"oauth_token_id":     vcs["oauth_token_id"],
		"identifier":         vcs["identifier"],
		"ingress_submodules": vcs["ingress_submodules"],
	}
	if vcs["branch"].IsKnown() && !vcs["branch"].IsNull() {
		finalVCS["branch"] = vcs["branch"]
	} else {
		finalVCS["branch"] = tftypes.NewValue(tftypes.String, workspace.VCSRepo.Branch)
	}
	finalState["vcs_repo"] = tftypes.NewValue(t.vcsRepoType(), finalVCS)
	dv, err := tfprotov5.NewDynamicValue(t.workspaceType(), tftypes.NewValue(t.workspaceType(), finalState))
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error returning updated state",
					Detail:   "The resource encountered an unexpected error returning the updated state. This indicates an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	return &tfprotov5.ApplyResourceChangeResponse{
		NewState: &dv,
	}, nil
}

func (t *terraform) ImportResourceState(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	client, err := t.clients.NewClient()
	if err != nil {
		return &tfprotov5.ImportResourceStateResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error creating client",
					Detail:   "The provider was unable to create a client.\n\nError:\n" + err.Error(),
				},
			},
		}, nil
	}
	workspace, err := client.Terraform.Workspaces.Get(ctx, req.ID)
	if err != nil {
		return &tfprotov5.ImportResourceStateResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error retrieving workspace",
					Detail:   "The provider was unable to retrieve the workspace.\n\nError:\n" + err.Error(),
				},
			},
		}, nil
	}
	triggerPrefixes := make([]tftypes.Value, 0, len(workspace.TriggerPrefixes))
	for _, prefix := range workspace.TriggerPrefixes {
		triggerPrefixes = append(triggerPrefixes, tftypes.NewValue(tftypes.String, prefix))
	}
	dv, err := tfprotov5.NewDynamicValue(t.workspaceType(), tftypes.NewValue(t.workspaceType(), map[string]tftypes.Value{
		"id":                    tftypes.NewValue(tftypes.String, workspace.ID),
		"name":                  tftypes.NewValue(tftypes.String, workspace.Name),
		"agent_pool_id":         tftypes.NewValue(tftypes.String, workspace.AgentPoolID),
		"allow_destroy_plan":    tftypes.NewValue(tftypes.Bool, workspace.AllowDestroyPlan),
		"auto_apply":            tftypes.NewValue(tftypes.Bool, workspace.AutoApply),
		"description":           tftypes.NewValue(tftypes.String, workspace.Description),
		"execution_mode":        tftypes.NewValue(tftypes.String, workspace.ExecutionMode),
		"file_triggers_enabled": tftypes.NewValue(tftypes.Bool, workspace.FileTriggersEnabled),
		"queue_all_runs":        tftypes.NewValue(tftypes.Bool, workspace.QueueAllRuns),
		"speculative_enabled":   tftypes.NewValue(tftypes.Bool, workspace.SpeculativeEnabled),
		"terraform_version":     tftypes.NewValue(tftypes.String, workspace.TerraformVersion),
		"trigger_prefixes":      tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, triggerPrefixes),
		"working_directory":     tftypes.NewValue(tftypes.String, workspace.WorkingDirectory),
		"vcs_repo": tftypes.NewValue(t.vcsRepoType(), map[string]tftypes.Value{
			"oauth_token_id":     tftypes.NewValue(tftypes.String, workspace.VCSRepo.OAuthTokenID),
			"branch":             tftypes.NewValue(tftypes.String, workspace.VCSRepo.Branch),
			"identifier":         tftypes.NewValue(tftypes.String, workspace.VCSRepo.Identifier),
			"ingress_submodules": tftypes.NewValue(tftypes.Bool, workspace.VCSRepo.IngressSubmodules),
		}),
	}))
	if err != nil {
		return &tfprotov5.ImportResourceStateResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error returning resource state",
					Detail:   "The resource encountered an unexpected error returning the imported state. This indicates an error in the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	return &tfprotov5.ImportResourceStateResponse{
		ImportedResources: []*tfprotov5.ImportedResource{
			{
				TypeName: req.TypeName,
				State:    &dv,
			},
		},
	}, nil
}
