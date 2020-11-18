package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
)

type API struct {
	Storer *Storer
}

func (a API) Server(baseURL string) http.Handler {
	var router trout.Router
	router.SetPrefix(baseURL)

	// get information about which regions support which products
	router.Endpoint("/regions").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetRegions))

	// create Vault cluster
	router.Endpoint("/vault/clusters").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostVaultCluster))
	// read Vault cluster
	router.Endpoint("/vault/clusters/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetVaultCluster))
	// update Vault cluster
	router.Endpoint("/vault/clusters/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutVaultCluster))
	// delete Vault cluster
	router.Endpoint("/vault/clusters/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteVaultCluster))

	// create Terraform workspace
	router.Endpoint("/terraform/workspaces").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostTerraformWorkspace))
	// read Terraform workspace
	router.Endpoint("/terraform/workspaces/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetTerraformWorkspace))
	// update Terraform workspace
	router.Endpoint("/terraform/workspaces/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutTerraformWorkspace))
	// delete Terraform workspace
	router.Endpoint("/terraform/workspaces/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteTerraformWorkspace))

	// create Consul cluster
	router.Endpoint("/consul/clusters").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostConsulCluster))
	// read Consul cluster
	router.Endpoint("/consul/clusters/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetConsulCluster))
	// update Consul cluster
	router.Endpoint("/consul/clusters/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutConsulCluster))
	// delete Consul cluster
	router.Endpoint("/consul/clusters/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteConsulCluster))

	// create Nomad cluster
	router.Endpoint("/nomad/clusters").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostNomadCluster))
	// read Nomad cluster
	router.Endpoint("/nomad/clusters/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetNomadCluster))
	// update Nomad cluster
	router.Endpoint("/nomad/clusters/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutNomadCluster))
	// delete Nomad cluster
	router.Endpoint("/nomad/clusters/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteNomadCluster))

	// create access policy
	router.Endpoint("/accessPolicies").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostAccessPolicy))
	// read access policy
	router.Endpoint("/accessPolicies/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetAccessPolicy))
	// update access policy
	router.Endpoint("/accessPolicies/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutAccessPolicy))
	// delete access policy
	router.Endpoint("/accessPolicies/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteAccessPolicy))

	return api.NegotiateMiddleware(router)
}

func isAuthenticated(r *http.Request) bool {
	un, pw, ok := r.BasicAuth()
	if !ok || un != "admin" || pw != "hunter2" {
		return false
	}
	return true
}

type Response struct {
	Regions             []Region             `json:"regions,omitempty"`
	VaultClusters       []VaultCluster       `json:"vaultClusters,omitempty"`
	TerraformWorkspaces []TerraformWorkspace `json:"terraformWorkspaces,omitempty"`
	ConsulClusters      []ConsulCluster      `json:"consulClusters,omitempty"`
	NomadClusters       []NomadCluster       `json:"nomadClusters,omitempty"`
	AccessPolicies      []AccessPolicy       `json:"accessPolicies,omitempty"`
	Errors              []api.RequestError   `json:"errors,omitempty"`
	Status              int                  `json:"-"`
}
