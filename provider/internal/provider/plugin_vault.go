package provider

import (
	"context"

	dadcorp "dadcorp.dev/client"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tftypes"
)

type vault struct {
	clients clientFactory
}

func (v *vault) clusterType() tftypes.Type {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":                tftypes.String,
			"name":              tftypes.String,
			"region":            tftypes.String,
			"default_lease_ttl": tftypes.String,
			"max_lease_ttl":     tftypes.String,
			"tcp_listener":      v.tcpListenerType(),
		},
	}
}

func (v *vault) tcpListenerType() tftypes.Type {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"address":         tftypes.String,
			"cluster_address": tftypes.String,
		},
	}
}

func (v *vault) schema() *tfprotov5.Schema {
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
					Name:     "region",
					Type:     tftypes.String,
					Required: true,
				},
				{
					Name:     "default_lease_ttl",
					Type:     tftypes.String,
					Optional: true,
					Computed: true,
				},
				{
					Name:     "max_lease_ttl",
					Type:     tftypes.String,
					Optional: true,
					Computed: true,
				},
			},
			BlockTypes: []*tfprotov5.SchemaNestedBlock{
				{
					TypeName: "tcp_listener",
					Nesting:  tfprotov5.SchemaNestedBlockNestingModeSingle,
					Block: &tfprotov5.SchemaBlock{
						Attributes: []*tfprotov5.SchemaAttribute{
							{
								Name:     "address",
								Type:     tftypes.String,
								Optional: true,
								Computed: true,
							},
							{
								Name:     "cluster_address",
								Type:     tftypes.String,
								Optional: true,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func (v *vault) ValidateResourceTypeConfig(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	val, err := req.Config.Unmarshal(v.clusterType())
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
	if !val.Is(v.clusterType()) {
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
	if values["region"].IsKnown() {
		var region string
		err = values["region"].As(&region)
		if err != nil {
			return &tfprotov5.ValidateResourceTypeConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected configuration format",
						Detail:   "The resource got a configuration that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("region"),
							},
						},
					},
				},
			}, nil
		}
		client, err := v.clients.NewClient()
		if err != nil {
			return &tfprotov5.ValidateResourceTypeConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Error creating client",
						Detail:   "The provider was unable to create a client.\n\nError:\n" + err.Error(),
					},
				},
			}, nil
		}
		regions, err := client.Regions.List(ctx)
		if err != nil {
			return &tfprotov5.ValidateResourceTypeConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Error retrieving regions",
						Detail:   "The provider was unable to retrieve a list of regions from the API.\n\nError:\n" + err.Error(),
					},
				},
			}, nil
		}
		var found bool
		for _, r := range regions {
			if r.ID == region {
				found = true
			}
		}
		if !found {
			return &tfprotov5.ValidateResourceTypeConfigResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unknown region",
						Detail:   "This is not a valid region.",
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("region"),
							},
						},
					},
				},
			}, nil
		}
	}
	return &tfprotov5.ValidateResourceTypeConfigResponse{}, nil
}

func (v *vault) UpgradeResourceState(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	switch req.Version {
	case 1:
		val, err := req.RawState.Unmarshal(v.clusterType())
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
		dv, err := tfprotov5.NewDynamicValue(v.clusterType(), val)
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

func (v *vault) ReadResource(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	val, err := req.CurrentState.Unmarshal(v.clusterType())
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
	cluster, err := client.Vault.Clusters.Get(ctx, id)
	if err != nil {
		if err == dadcorp.ErrVaultClusterNotFound {
			dv, err := tfprotov5.NewDynamicValue(v.clusterType(), tftypes.NewValue(v.clusterType(), nil))
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
					Summary:  "Error retrieving cluster",
					Detail:   "The provider was unable to retrieve the cluster.\n\nError:\n" + err.Error(),
				},
			},
		}, nil
	}
	dv, err := tfprotov5.NewDynamicValue(v.clusterType(), tftypes.NewValue(v.clusterType(), map[string]tftypes.Value{
		"id":                tftypes.NewValue(tftypes.String, cluster.ID),
		"name":              tftypes.NewValue(tftypes.String, cluster.Name),
		"region":            tftypes.NewValue(tftypes.String, cluster.Region),
		"default_lease_ttl": tftypes.NewValue(tftypes.String, cluster.DefaultLeaseTTL),
		"max_lease_ttl":     tftypes.NewValue(tftypes.String, cluster.MaxLeaseTTL),
		"tcp_listener": tftypes.NewValue(v.tcpListenerType(), map[string]tftypes.Value{
			"address":         tftypes.NewValue(tftypes.String, cluster.TCPListener.Address),
			"cluster_address": tftypes.NewValue(tftypes.String, cluster.TCPListener.ClusterAddress),
		}),
	}))
	if err != nil {
		return &tfprotov5.ReadResourceResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error updating cluster in state",
					Detail:   "An unexpected error was encountered updating the cluster from state. This is an error with the provider.\n\nError: " + err.Error(),
				},
			},
		}, nil
	}
	return &tfprotov5.ReadResourceResponse{
		NewState: &dv,
	}, nil
}

func (v *vault) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	val, err := req.ProposedNewState.Unmarshal(v.clusterType())
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
	if newState["default_lease_ttl"].IsNull() {
		newState["default_lease_ttl"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	}
	if newState["max_lease_ttl"].IsNull() {
		newState["max_lease_ttl"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	}
	if !newState["tcp_listener"].IsNull() && newState["tcp_listener"].IsKnown() {
		tcp := map[string]tftypes.Value{}
		err = newState["tcp_listener"].As(&tcp)
		if err != nil {
			return &tfprotov5.PlanResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected state format",
						Detail:   "The resource got a state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("tcp_listener"),
							},
						},
					},
				},
			}, nil
		}
		newTCP := map[string]tftypes.Value{
			"address":         tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"cluster_address": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		}
		if !tcp["address"].IsNull() {
			newTCP["address"] = tcp["address"]
		}
		if !tcp["cluster_address"].IsNull() {
			newTCP["cluster_address"] = tcp["cluster_address"]
		}
		newState["tcp_listener"] = tftypes.NewValue(v.tcpListenerType(), newTCP)
	}
	dv, err := tfprotov5.NewDynamicValue(v.clusterType(), tftypes.NewValue(v.clusterType(), newState))
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

func (v *vault) ApplyResourceChange(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	plannedStateVal, err := req.PlannedState.Unmarshal(v.clusterType())
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
	priorStateVal, err := req.PriorState.Unmarshal(v.clusterType())
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

	// if plannedStateVal is null, we're deleting the cluster
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
		err = client.Vault.Clusters.Delete(ctx, id)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Error deleting cluster",
						Detail:   "The provider was unable to delete the cluster.\n\nError:\n" + err.Error(),
					},
				},
			}, nil
		}
		dv, err := tfprotov5.NewDynamicValue(v.clusterType(), tftypes.NewValue(v.clusterType(), nil))
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

	// if plannedStateVal is not null, we're creating or updating the cluster
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

	var cluster dadcorp.VaultCluster
	err = plannedState["name"].As(&cluster.Name)
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
	err = plannedState["region"].As(&cluster.Region)
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected planned state format",
					Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					Attribute: &tftypes.AttributePath{
						Steps: []tftypes.AttributePathStep{
							tftypes.AttributeName("region"),
						},
					},
				},
			},
		}, nil
	}
	if plannedState["default_lease_ttl"].IsKnown() && !plannedState["default_lease_ttl"].IsNull() {
		err = plannedState["default_lease_ttl"].As(&cluster.DefaultLeaseTTL)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("default_lease_ttl"),
							},
						},
					},
				},
			}, nil
		}
	}
	if plannedState["max_lease_ttl"].IsKnown() && !plannedState["max_lease_ttl"].IsNull() {
		err = plannedState["max_lease_ttl"].As(&cluster.MaxLeaseTTL)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("max_lease_ttl"),
							},
						},
					},
				},
			}, nil
		}
	}
	tcp := map[string]tftypes.Value{}
	err = plannedState["tcp_listener"].As(&tcp)
	if err != nil {
		return &tfprotov5.ApplyResourceChangeResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected planned state format",
					Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
					Attribute: &tftypes.AttributePath{
						Steps: []tftypes.AttributePathStep{
							tftypes.AttributeName("tcp_listener"),
						},
					},
				},
			},
		}, nil
	}
	if tcp["address"].IsKnown() && !tcp["address"].IsNull() {
		err = tcp["address"].As(&cluster.TCPListener.Address)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("tcp_listener"),
								tftypes.AttributeName("address"),
							},
						},
					},
				},
			}, nil
		}
	}
	if tcp["cluster_address"].IsKnown() && !tcp["cluster_address"].IsNull() {
		err = tcp["cluster_address"].As(&cluster.TCPListener.ClusterAddress)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected planned state format",
						Detail:   "The resource got a planned state that did not match its schema. This may indicate an error in the provider.\n\nError: " + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("tcp_listener"),
								tftypes.AttributeName("cluster_address"),
							},
						},
					},
				},
			}, nil
		}
	}

	// if priorStateVal is not null, we're updating the cluster
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
		err = priorState["id"].As(&cluster.ID)
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
		cluster, err = client.Vault.Clusters.Update(ctx, cluster)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Error updating the cluster",
						Detail:   "The provider was unable to update the cluster.\n\nError: " + err.Error(),
					},
				},
			}, nil
		}
	} else {
		// if priorStateVal is null, we're creating the cluster
		cluster, err = client.Vault.Clusters.Create(ctx, cluster)
		if err != nil {
			return &tfprotov5.ApplyResourceChangeResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Error creating the cluster",
						Detail:   "The provider was unable to create the cluster.\n\nError: " + err.Error(),
					},
				},
			}, nil
		}
	}

	finalState := map[string]tftypes.Value{
		"id":     tftypes.NewValue(tftypes.String, cluster.ID),
		"name":   plannedState["name"],
		"region": plannedState["region"],
	}
	if plannedState["default_lease_ttl"].IsKnown() && !plannedState["default_lease_ttl"].IsNull() {
		finalState["default_lease_ttl"] = plannedState["default_lease_ttl"]
	} else {
		finalState["default_lease_ttl"] = tftypes.NewValue(tftypes.String, cluster.DefaultLeaseTTL)
	}
	if plannedState["max_lease_ttl"].IsKnown() && !plannedState["max_lease_ttl"].IsNull() {
		finalState["max_lease_ttl"] = plannedState["max_lease_ttl"]
	} else {
		finalState["max_lease_ttl"] = tftypes.NewValue(tftypes.String, cluster.MaxLeaseTTL)
	}
	finalTCP := map[string]tftypes.Value{}
	if tcp["address"].IsKnown() && !tcp["address"].IsNull() {
		finalTCP["address"] = tcp["address"]
	} else {
		finalTCP["address"] = tftypes.NewValue(tftypes.String, cluster.TCPListener.Address)
	}
	if tcp["cluster_address"].IsKnown() && !tcp["cluster_address"].IsNull() {
		finalTCP["cluster_address"] = tcp["cluster_address"]
	} else {
		finalTCP["cluster_address"] = tftypes.NewValue(tftypes.String, cluster.TCPListener.ClusterAddress)
	}
	finalState["tcp_listener"] = tftypes.NewValue(v.tcpListenerType(), finalTCP)
	dv, err := tfprotov5.NewDynamicValue(v.clusterType(), tftypes.NewValue(v.clusterType(), finalState))
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

func (v *vault) ImportResourceState(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
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
	cluster, err := client.Vault.Clusters.Get(ctx, req.ID)
	if err != nil {
		return &tfprotov5.ImportResourceStateResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error retrieving cluster",
					Detail:   "The provider was unable to retrieve the cluster.\n\nError:\n" + err.Error(),
				},
			},
		}, nil
	}
	dv, err := tfprotov5.NewDynamicValue(v.clusterType(), tftypes.NewValue(v.clusterType(), map[string]tftypes.Value{
		"id":                tftypes.NewValue(tftypes.String, cluster.ID),
		"name":              tftypes.NewValue(tftypes.String, cluster.Name),
		"region":            tftypes.NewValue(tftypes.String, cluster.Region),
		"default_lease_ttl": tftypes.NewValue(tftypes.String, cluster.DefaultLeaseTTL),
		"max_lease_ttl":     tftypes.NewValue(tftypes.String, cluster.MaxLeaseTTL),
		"tcp_listener": tftypes.NewValue(v.tcpListenerType(), map[string]tftypes.Value{
			"address":         tftypes.NewValue(tftypes.String, cluster.TCPListener.Address),
			"cluster_address": tftypes.NewValue(tftypes.String, cluster.TCPListener.ClusterAddress),
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
