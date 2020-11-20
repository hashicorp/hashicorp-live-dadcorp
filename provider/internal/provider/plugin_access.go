package provider

import (
	"context"
	"fmt"

	dadcorp "dadcorp.dev/client"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tftypes"
)

type accessPolicy struct {
	clients clientFactory
}

func PolicyDataToTerraformValue(policyType string, policyData map[string]interface{}) (tftypes.Value, error) {
	switch policyType {
	case "terraform":
		return tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"workspace_id":      tftypes.String,
				"plan":              tftypes.Bool,
				"apply":             tftypes.Bool,
				"override_policies": tftypes.Bool,
			},
		}, map[string]tftypes.Value{
			"workspace_id":      tftypes.NewValue(tftypes.String, policyData["id"]),
			"plan":              tftypes.NewValue(tftypes.Bool, policyData["plan"]),
			"apply":             tftypes.NewValue(tftypes.Bool, policyData["apply"]),
			"override_policies": tftypes.NewValue(tftypes.Bool, policyData["overridePolicies"]),
		}), nil
	case "vault":
		return tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"cluster_id": tftypes.String,
				"key":        tftypes.String,
				"read":       tftypes.Bool,
				"write":      tftypes.Bool,
				"delete":     tftypes.Bool,
			},
		}, map[string]tftypes.Value{
			"cluster_id": tftypes.NewValue(tftypes.String, policyData["id"]),
			"key":        tftypes.NewValue(tftypes.String, policyData["key"]),
			"read":       tftypes.NewValue(tftypes.Bool, policyData["read"]),
			"write":      tftypes.NewValue(tftypes.Bool, policyData["write"]),
			"delete":     tftypes.NewValue(tftypes.Bool, policyData["delete"]),
		}), nil
	case "nomad":
		return tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"cluster_id":      tftypes.String,
				"submit_jobs":     tftypes.Bool,
				"read_job_status": tftypes.Bool,
				"cancel_jobs":     tftypes.Bool,
			},
		}, map[string]tftypes.Value{
			"cluster_id":      tftypes.NewValue(tftypes.String, policyData["id"]),
			"submit_jobs":     tftypes.NewValue(tftypes.Bool, policyData["submitJobs"]),
			"read_job_status": tftypes.NewValue(tftypes.Bool, policyData["readJobStatus"]),
			"cancel_jobs":     tftypes.NewValue(tftypes.Bool, policyData["cancelJobs"]),
		}), nil
	case "consul":
		return tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"cluster_id": tftypes.String,
				"key":        tftypes.String,
				"read":       tftypes.Bool,
				"write":      tftypes.Bool,
				"delete":     tftypes.Bool,
			},
		}, map[string]tftypes.Value{
			"cluster_id": tftypes.NewValue(tftypes.String, policyData["id"]),
			"key":        tftypes.NewValue(tftypes.String, policyData["key"]),
			"read":       tftypes.NewValue(tftypes.Bool, policyData["read"]),
			"write":      tftypes.NewValue(tftypes.Bool, policyData["write"]),
			"delete":     tftypes.NewValue(tftypes.Bool, policyData["delete"]),
		}), nil
	default:
		return tftypes.Value{}, fmt.Errorf("unexpected policy type: %s", policyType)
	}
}

type TerraformPolicy struct {
	WorkspaceID      string
	Plan             bool
	Apply            bool
	OverridePolicies bool
}

func (p *TerraformPolicy) FromTerraform5Value(val tftypes.Value) error {
	v := map[string]tftypes.Value{}
	err := val.As(&v)
	if err != nil {
		return err
	}

	err = v["workspace_id"].As(&p.WorkspaceID)
	if err != nil {
		return err
	}

	err = v["plan"].As(&p.Plan)
	if err != nil {
		return err
	}

	err = v["apply"].As(&p.Apply)
	if err != nil {
		return err
	}

	err = v["override_policies"].As(&p.OverridePolicies)
	if err != nil {
		return err
	}

	return nil
}

type VaultPolicy struct {
	ClusterID string
	Key       string
	Read      bool
	Write     bool
	Delete    bool
}

func (p *VaultPolicy) FromTerraform5Value(val tftypes.Value) error {
	v := map[string]tftypes.Value{}
	err := val.As(&v)
	if err != nil {
		return err
	}

	err = v["cluster_id"].As(&p.ClusterID)
	if err != nil {
		return err
	}

	err = v["key"].As(&p.Key)
	if err != nil {
		return err
	}

	err = v["read"].As(&p.Read)
	if err != nil {
		return err
	}

	err = v["write"].As(&p.Write)
	if err != nil {
		return err
	}

	err = v["delete"].As(&p.Delete)
	if err != nil {
		return err
	}

	return nil
}

type NomadPolicy struct {
	ClusterID     string
	SubmitJobs    bool
	ReadJobStatus bool
	CancelJobs    bool
}

func (p *NomadPolicy) FromTerraform5Value(val tftypes.Value) error {
	v := map[string]tftypes.Value{}
	err := val.As(&v)
	if err != nil {
		return err
	}

	err = v["cluster_id"].As(&p.ClusterID)
	if err != nil {
		return err
	}

	err = v["submit_jobs"].As(&p.SubmitJobs)
	if err != nil {
		return err
	}

	err = v["read_job_status"].As(&p.ReadJobStatus)
	if err != nil {
		return err
	}

	err = v["cancel_jobs"].As(&p.CancelJobs)
	if err != nil {
		return err
	}

	return nil
}

type ConsulPolicy struct {
	ClusterID string
	Key       string
	Read      bool
	Write     bool
	Delete    bool
}

func (p *ConsulPolicy) FromTerraform5Value(val tftypes.Value) error {
	v := map[string]tftypes.Value{}
	err := val.As(&v)
	if err != nil {
		return err
	}

	err = v["cluster_id"].As(&p.ClusterID)
	if err != nil {
		return err
	}

	err = v["key"].As(&p.Key)
	if err != nil {
		return err
	}

	err = v["read"].As(&p.Read)
	if err != nil {
		return err
	}

	err = v["write"].As(&p.Write)
	if err != nil {
		return err
	}

	err = v["delete"].As(&p.Delete)
	if err != nil {
		return err
	}

	return nil
}

func (v *accessPolicy) accessPolicyType() tftypes.Type {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":          tftypes.String,
			"type":        tftypes.String,
			"policy_data": tftypes.DynamicPseudoType,
		},
	}
}

func (v *accessPolicy) schema() *tfprotov5.Schema {
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
					Name:     "type",
					Type:     tftypes.String,
					Required: true,
				},
				{
					Name:     "policy_data",
					Type:     tftypes.DynamicPseudoType,
					Required: true,
				},
			},
		},
	}
}

func (v *accessPolicy) ValidateResourceTypeConfig(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	val, err := req.Config.Unmarshal(v.accessPolicyType())
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
	if !val.Is(v.accessPolicyType()) {
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
	if values["type"].IsKnown() {
		var policyType string
		err = values["type"].As(&policyType)
		if err != nil {
			return &tfprotov5.ValidateResourceTypeConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected configuration format",
						Detail:   "The resource got a configuration that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("type"),
							},
						},
					},
				},
			}, nil
		}
	}
	return &tfprotov5.ValidateResourceTypeConfigResponse{}, nil
}

func (v *accessPolicy) UpgradeResourceState(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	switch req.Version {
	case 1:
		val, err := req.RawState.Unmarshal(v.accessPolicyType())
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
		dv, err := tfprotov5.NewDynamicValue(v.accessPolicyType(), val)
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

func (v *accessPolicy) ReadResource(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	val, err := req.CurrentState.Unmarshal(v.accessPolicyType())
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
	client, err := v.clients.NewClient()
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
	policy, err := client.AccessPolicies.Get(ctx, id)
	if err != nil {
		if err == dadcorp.ErrAccessPolicyNotFound {
			dv, err := tfprotov5.NewDynamicValue(v.accessPolicyType(), tftypes.NewValue(v.accessPolicyType(), nil))
			if err != nil {
				return &tfprotov5.ReadResourceResponse{
					Diagnostics: []*tfprotov5.Diagnostic{
						{
							Severity: tfprotov5.DiagnosticSeverityError,
							Summary:  "Error removing cluster from state",
							Detail:   "An unexpected error was encountered removing the cluster from state. This is an error with the provider.\n\nError: " + err.Error(),
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
					Summary:  "Error retrieving access policies",
					Detail:   "The provider was unable to retrieve the access policies.\n\nError:\n" + err.Error(),
				},
			},
		}, nil
	}
	policyData, err := PolicyDataToTerraformValue(policy.Type, policy.PolicyData.(map[string]interface{}))
	if err != nil {
		return &tfprotov5.ReadResourceResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error parsing policy data",
					Detail:   "An unexpected error was encountered while parsing the policy data. This is an error with the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	dv, err := tfprotov5.NewDynamicValue(v.accessPolicyType(), tftypes.NewValue(v.accessPolicyType(), map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, policy.ID),
		"type":        tftypes.NewValue(tftypes.String, policy.Type),
		"policy_data": policyData,
	}))
	if err != nil {
		return &tfprotov5.ReadResourceResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error updating access policies in state",
					Detail:   "An unexpected error was encountered updating the access policies from state. This is an error with the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	return &tfprotov5.ReadResourceResponse{
		NewState: &dv,
	}, nil
}

func (v *accessPolicy) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	val, err := req.ProposedNewState.Unmarshal(v.accessPolicyType())
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
	if newState["id"].IsNull() {
		newState["id"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	}

	dv, err := tfprotov5.NewDynamicValue(v.accessPolicyType(), tftypes.NewValue(v.accessPolicyType(), newState))
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

func (v *accessPolicy) ApplyResourceChange(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	plannedStateVal, err := req.PlannedState.Unmarshal(v.accessPolicyType())
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
	priorStateVal, err := req.PriorState.Unmarshal(v.accessPolicyType())
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
	client, err := v.clients.NewClient()
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

	// if plannedStateVal is null, we're deleting the policy
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
		err = client.AccessPolicies.Delete(ctx, id)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Error deleting access policy",
						Detail:   "The provider was unable to delete the access policy.\n\nError:\n" + err.Error(),
					},
				},
			}, nil
		}
		dv, err := tfprotov5.NewDynamicValue(v.accessPolicyType(), tftypes.NewValue(v.accessPolicyType(), nil))
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

	// if plannedStateVal is not null, we're creating or updating the policy
	// so let's get to the planned state
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

	var accessPolicy dadcorp.AccessPolicy
	err = plannedState["type"].As(&accessPolicy.Type)
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected planned state format",
					Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					Attribute: &tftypes.AttributePath{
						Steps: []tftypes.AttributePathStep{
							tftypes.AttributeName("type"),
						},
					},
				},
			},
		}, nil
	}

	switch accessPolicy.Type {
	case "terraform":
		var data TerraformPolicy
		err = plannedState["policy_data"].As(&data)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected policy data format",
						Detail:   "\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("policy_data"),
							},
						},
					},
				},
			}, nil
		}
		accessPolicy.PolicyData = dadcorp.TerraformPolicy(data)
	case "vault":
		var data VaultPolicy
		err = plannedState["policy_data"].As(&data)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected policy data format",
						Detail:   "\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("policy_data"),
							},
						},
					},
				},
			}, nil
		}
		accessPolicy.PolicyData = dadcorp.VaultPolicy(data)
	case "nomad":
		var data NomadPolicy
		err = plannedState["policy_data"].As(&data)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected policy data format",
						Detail:   "\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("policy_data"),
							},
						},
					},
				},
			}, nil
		}
		accessPolicy.PolicyData = dadcorp.NomadPolicy(data)
	case "consul":
		var data ConsulPolicy
		err = plannedState["policy_data"].As(&data)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected policy data format",
						Detail:   "\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("policy_data"),
							},
						},
					},
				},
			}, nil
		}
		accessPolicy.PolicyData = dadcorp.ConsulPolicy(data)
	default:
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected policy data format",
					Detail:   "Access policies must be of type consul, nomad, terraform, or vault",
					Attribute: &tftypes.AttributePath{
						Steps: []tftypes.AttributePathStep{
							tftypes.AttributeName("policy_data"),
						},
					},
				},
			},
		}, nil
	}

	// if priorStateVal is not null, we're updating the policy
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
		err = priorState["id"].As(&accessPolicy.ID)
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
		accessPolicy, err = client.AccessPolicies.Update(ctx, accessPolicy)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Error updating the access policy",
						Detail:   "The provider was unable to update the access policy.\n\nError: " + err.Error(),
					},
				},
			}, nil
		}
	} else {
		// if priorStateVal is null, we're creating the cluster
		accessPolicy, err = client.AccessPolicies.Create(ctx, accessPolicy)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Error creating the access policy",
						Detail:   "The provider was unable to create the access policy.\n\nError: " + err.Error(),
					},
				},
			}, nil
		}
	}

	finalState := map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, accessPolicy.ID),
		"type":        plannedState["type"],
		"policy_data": plannedState["policy_data"],
	}
	dv, err := tfprotov5.NewDynamicValue(v.accessPolicyType(), tftypes.NewValue(v.accessPolicyType(), finalState))
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

func (v *accessPolicy) ImportResourceState(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	client, err := v.clients.NewClient()
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
	accessPolicy, err := client.AccessPolicies.Get(ctx, req.ID)
	if err != nil {
		return &tfprotov5.ImportResourceStateResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error retrieving access policy",
					Detail:   "The provider was unable to retrieve the access policy.\n\nError:\n" + err.Error(),
				},
			},
		}, nil
	}
	policyData, err := PolicyDataToTerraformValue(accessPolicy.Type, accessPolicy.PolicyData.(map[string]interface{}))
	if err != nil {
		return &tfprotov5.ImportResourceStateResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error parsing policy data",
					Detail:   "An unexpected error was encountered while parsing the policy data. This is an error with the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	dv, err := tfprotov5.NewDynamicValue(v.accessPolicyType(), tftypes.NewValue(v.accessPolicyType(), map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, accessPolicy.ID),
		"type":        tftypes.NewValue(tftypes.String, accessPolicy.Type),
		"policy_data": policyData,
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
