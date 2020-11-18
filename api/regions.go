package api

import (
	"net/http"

	"darlinggo.co/api"
)

type Region struct {
	ID       string         `json:"id"`
	Products RegionProducts `json:"products"`
}

type RegionProducts struct {
	Vault     bool `json:"vault"`
	Terraform bool `json:"terraform"`
	Nomad     bool `json:"nomad"`
	Consul    bool `json:"consul"`
}

func (a API) handleGetRegions(w http.ResponseWriter, r *http.Request) {
	regions := getRegions(isAuthenticated(r))
	api.Encode(w, r, http.StatusOK, Response{Regions: regions})
}

func getRegions(authenticated bool) []Region {
	return []Region{
		{
			ID: "us-va-1",
			Products: RegionProducts{
				Vault:     true,
				Consul:    true,
				Nomad:     authenticated,
				Terraform: authenticated,
			},
		},
		{
			ID: "us-va-2",
			Products: RegionProducts{
				Vault:     true,
				Consul:    true,
				Nomad:     authenticated,
				Terraform: authenticated,
			},
		},
		{
			ID: "us-or-1",
			Products: RegionProducts{
				Vault:     authenticated,
				Consul:    true,
				Nomad:     false,
				Terraform: false,
			},
		},
		{
			ID: "us-or-2",
			Products: RegionProducts{
				Vault:     authenticated,
				Consul:    true,
				Nomad:     false,
				Terraform: false,
			},
		},
		{
			ID: "gb-lon-1",
			Products: RegionProducts{
				Vault:     authenticated,
				Consul:    true,
				Nomad:     false,
				Terraform: false,
			},
		},
		{
			ID: "gb-lon-2",
			Products: RegionProducts{
				Vault:     authenticated,
				Consul:    true,
				Nomad:     false,
				Terraform: false,
			},
		},
		{
			ID: "jp-tok-1",
			Products: RegionProducts{
				Vault:     authenticated,
				Consul:    true,
				Nomad:     false,
				Terraform: false,
			},
		},
		{
			ID: "jp-tok-2",
			Products: RegionProducts{
				Vault:     authenticated,
				Consul:    true,
				Nomad:     false,
				Terraform: false,
			},
		},
	}
}
