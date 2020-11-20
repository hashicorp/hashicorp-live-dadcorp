package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tftypes"
)

func NewPlugin() tfprotov5.ProviderServer {
	return &provider{}
}

type provider struct {
	clientFactory clientFactory
}

func (p *provider) GetProviderSchema(ctx context.Context, req *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	return &tfprotov5.GetProviderSchemaResponse{
		Provider: &tfprotov5.Schema{
			Block: &tfprotov5.SchemaBlock{
				Attributes: []*tfprotov5.SchemaAttribute{
					{
						Name:     "password",
						Type:     tftypes.String,
						Optional: true,
					},
					{
						Name:     "username",
						Type:     tftypes.String,
						Optional: true,
					},
				},
			},
		},
		ResourceSchemas: map[string]*tfprotov5.Schema{
			"dadcorp_vault_cluster": (&vault{}).schema(),
		},
		DataSourceSchemas: map[string]*tfprotov5.Schema{},
	}, nil
}

func (p *provider) PrepareProviderConfig(ctx context.Context, req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	return &tfprotov5.PrepareProviderConfigResponse{
		PreparedConfig: nil,
	}, nil
}

func (p *provider) ConfigureProvider(ctx context.Context, req *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	configType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"username": tftypes.String,
			"password": tftypes.String,
		},
	}
	var client clientFactory
	val, err := req.Config.Unmarshal(configType)
	if err != nil {
		return &tfprotov5.ConfigureProviderResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected provider configuration",
					Detail:   "The provider got a configuration that did not match its schema. This may indicate an error in the provider.\n\nError:" + err.Error(),
				},
			},
		}, err
	}
	values := map[string]tftypes.Value{}
	err = val.As(&values)
	if err != nil {
		return &tfprotov5.ConfigureProviderResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Unexpected provider configuration",
					Detail:   "The provider got a configuration that did not match its schema. This may indicate an error in the provider.\n\nError:" + err.Error(),
				},
			},
		}, err
	}
	if values["username"].IsKnown() && !values["username"].IsNull() {
		err = values["username"].As(&client.username)
		if err != nil {
			return &tfprotov5.ConfigureProviderResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected provider configuration",
						Detail:   "The provider got a configuration that did not match its schema. This may indicate an error in the provider.\n\nError:" + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("username"),
							},
						},
					},
				},
			}, err
		}
	}
	if values["password"].IsKnown() && !values["password"].IsNull() {
		err = values["password"].As(&client.password)
		if err != nil {
			return &tfprotov5.ConfigureProviderResponse{
				Diagnostics: []*tfprotov5.Diagnostic{
					{
						Severity: tfprotov5.DiagnosticSeverityError,
						Summary:  "Unexpected provider configuration",
						Detail:   "The provider got a configuration that did not match its schema. This may indicate an error in the provider.\n\nError:" + err.Error(),
						Attribute: &tftypes.AttributePath{
							Steps: []tftypes.AttributePathStep{
								tftypes.AttributeName("password"),
							},
						},
					},
				},
			}, err
		}
	}
	if os.Getenv("DADCORP_USERNAME") != "" {
		client.username = os.Getenv("DADCORP_USERNAME")
	}
	if os.Getenv("DADCORP_PASSWORD") != "" {
		client.password = os.Getenv("DADCORP_PASSWORD")
	}
	p.clientFactory = client
	return &tfprotov5.ConfigureProviderResponse{}, nil
}

func (p *provider) StopProvider(ctx context.Context, req *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	return &tfprotov5.StopProviderResponse{}, nil
}

// resource methods
func (p *provider) ValidateResourceTypeConfig(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	switch req.TypeName {
	case "dadcorp_vault_cluster":
		res := &vault{
			clients: p.clientFactory,
		}
		return res.ValidateResourceTypeConfig(ctx, req)
	}
	return &tfprotov5.ValidateResourceTypeConfigResponse{
		Diagnostics: []*tfprotov5.Diagnostic{
			{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Unknown resource",
				Detail:   fmt.Sprintf("Unknown resource %q can't be validated.", req.TypeName),
			},
		},
	}, nil
}

func (p *provider) UpgradeResourceState(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	switch req.TypeName {
	case "dadcorp_vault_cluster":
		res := &vault{
			clients: p.clientFactory,
		}
		return res.UpgradeResourceState(ctx, req)
	}
	return &tfprotov5.UpgradeResourceStateResponse{
		Diagnostics: []*tfprotov5.Diagnostic{
			{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Unknown resource",
				Detail:   fmt.Sprintf("Unknown resource %q can't be upgraded.", req.TypeName),
			},
		},
	}, nil
}

func (p *provider) ReadResource(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	switch req.TypeName {
	case "dadcorp_vault_cluster":
		res := &vault{
			clients: p.clientFactory,
		}
		return res.ReadResource(ctx, req)
	}
	return &tfprotov5.ReadResourceResponse{
		Diagnostics: []*tfprotov5.Diagnostic{
			{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Unknown resource",
				Detail:   fmt.Sprintf("Unknown resource %q can't be read.", req.TypeName),
			},
		},
	}, nil
}

func (p *provider) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	switch req.TypeName {
	case "dadcorp_vault_cluster":
		res := &vault{
			clients: p.clientFactory,
		}
		return res.PlanResourceChange(ctx, req)
	}
	return &tfprotov5.PlanResourceChangeResponse{
		Diagnostics: []*tfprotov5.Diagnostic{
			{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Unknown resource",
				Detail:   fmt.Sprintf("Unknown resource %q can't be planned.", req.TypeName),
			},
		},
	}, nil
}

func (p *provider) ApplyResourceChange(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	switch req.TypeName {
	case "dadcorp_vault_cluster":
		res := &vault{
			clients: p.clientFactory,
		}
		return res.ApplyResourceChange(ctx, req)
	}
	return &tfprotov5.ApplyResourceChangeResponse{
		Diagnostics: []*tfprotov5.Diagnostic{
			{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Unknown resource",
				Detail:   fmt.Sprintf("Unknown resource %q can't be applied.", req.TypeName),
			},
		},
	}, nil
}

func (p *provider) ImportResourceState(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	switch req.TypeName {
	case "dadcorp_vault_cluster":
		res := &vault{
			clients: p.clientFactory,
		}
		return res.ImportResourceState(ctx, req)
	}
	return &tfprotov5.ImportResourceStateResponse{
		Diagnostics: []*tfprotov5.Diagnostic{
			{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Unknown resource",
				Detail:   fmt.Sprintf("Unknown resource %q can't be imported.", req.TypeName),
			},
		},
	}, nil
}

// data source methods
func (p *provider) ValidateDataSourceConfig(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	switch req.TypeName {
	}
	return &tfprotov5.ValidateDataSourceConfigResponse{
		Diagnostics: []*tfprotov5.Diagnostic{
			{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Unknown data source",
				Detail:   fmt.Sprintf("Unknown data source %q can't be validated.", req.TypeName),
			},
		},
	}, nil
}

func (p *provider) ReadDataSource(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	switch req.TypeName {
	}
	return &tfprotov5.ReadDataSourceResponse{
		Diagnostics: []*tfprotov5.Diagnostic{
			{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Unknown data source",
				Detail:   fmt.Sprintf("Unknown data source %q can't be read.", req.TypeName),
			},
		},
	}, nil
}
