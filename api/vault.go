package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
	"github.com/hashicorp/go-uuid"
)

type VaultCluster struct {
	ID              string                  `json:"id"`
	Name            string                  `json:"name"`
	Region          string                  `json:"region"`
	DefaultLeaseTTL string                  `json:"defaultLeaseTTL"`
	MaxLeaseTTL     string                  `json:"maxLeaseTTL"`
	TCPListener     VaultClusterTCPListener `json:"tcpListener"`
}

type VaultClusterTCPListener struct {
	Address        string `json:"address"`
	ClusterAddress string `json:"clusterAddress"`
}

func (cluster *VaultCluster) FillDefaults() {
	if cluster.DefaultLeaseTTL == "" {
		cluster.DefaultLeaseTTL = "768h"
	}
	if cluster.MaxLeaseTTL == "" {
		cluster.MaxLeaseTTL = "768h"
	}
	if cluster.TCPListener.Address == "" {
		cluster.TCPListener.Address = "127.0.0.1:8200"
	}
	if cluster.TCPListener.ClusterAddress == "" {
		cluster.TCPListener.ClusterAddress = "127.0.0.1:8201"
	}
}

func (a API) handleGetVaultCluster(w http.ResponseWriter, r *http.Request) {
	cluster, err := a.Storer.GetVaultCluster(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrVaultClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{VaultClusters: []VaultCluster{cluster}})
}

func (a API) handlePostVaultCluster(w http.ResponseWriter, r *http.Request) {
	var cluster VaultCluster
	err := api.Decode(r, &cluster)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	regions := getRegions(isAuthenticated(r))
	var validRegion bool
	for _, region := range regions {
		if region.ID != cluster.Region {
			continue
		}
		validRegion = true
		if !region.Products.Vault {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/region", Slug: api.RequestErrAccessDenied}}})
			return
		}
	}
	if !validRegion {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/region", Slug: api.RequestErrInvalidValue}}})
		return
	}
	if cluster.ID == "" {
		cluster.ID, err = uuid.GenerateUUID()
		if err != nil {
			api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
			return
		}
	}
	if cluster.Name == "" {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/name", Slug: api.RequestErrMissing}}})
		return
	}
	cluster.FillDefaults()
	err = a.Storer.CreateVaultCluster(cluster)
	if err != nil {
		if err == ErrVaultClusterAlreadyExists {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusCreated, Response{VaultClusters: []VaultCluster{cluster}})
}

func (a API) handlePutVaultCluster(w http.ResponseWriter, r *http.Request) {
	var cluster VaultCluster
	err := api.Decode(r, &cluster)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	if cluster.ID != "" && cluster.ID != trout.RequestVars(r).Get("id") {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
		return
	}
	if cluster.Name == "" {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/name", Slug: api.RequestErrMissing}}})
		return
	}
	cluster.FillDefaults()
	err = a.Storer.UpdateVaultCluster(cluster)
	if err != nil {
		if err == ErrVaultClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{VaultClusters: []VaultCluster{cluster}})
}

func (a API) handleDeleteVaultCluster(w http.ResponseWriter, r *http.Request) {
	cluster, err := a.Storer.GetVaultCluster(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrVaultClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	err = a.Storer.DeleteVaultCluster(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrVaultClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{VaultClusters: []VaultCluster{cluster}})
}
