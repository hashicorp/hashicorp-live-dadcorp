package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
	"github.com/hashicorp/go-uuid"
)

type ConsulCluster struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	BindAddr  string                 `json:"bindAddr"`
	Addresses ConsulClusterAddresses `json:"addresses"`
	Ports     ConsulClusterPorts     `json:"ports"`
}

type ConsulClusterAddresses struct {
	DNS   string `json:"dns,omitempty"`
	HTTP  string `json:"http,omitempty"`
	HTTPS string `json:"https,omitempty"`
	GRPC  string `json:"grpc,omitempty"`
}

type ConsulClusterPorts struct {
	DNS            int  `json:"dns,omitempty"`
	HTTP           int  `json:"http,omitempty"`
	HTTPS          int  `json:"https,omitempty"`
	GRPC           int  `json:"grpc,omitempty"`
	SerfLAN        int  `json:"serfLan,omitempty"`
	SerfWAN        int  `json:"serfWan,omitempty"`
	Server         int  `json:"server,omitempty"`
	SidecarMinPort *int `json:"sidecarMinPort,omitempty"`
	SidecarMaxPort *int `json:"sidecarMaxPort,omitempty"`
	ExposeMinPort  *int `json:"exposeMinPort,omitempty"`
	ExposeMaxPort  *int `json:"exposeMaxPort,omitempty"`
}

func (cluster *ConsulCluster) FillDefaults() {
	if cluster.Addresses.DNS == "" {
		cluster.Addresses.DNS = "127.0.0.1"
	}
	if cluster.Addresses.HTTP == "" {
		cluster.Addresses.HTTP = "127.0.0.1"
	}
	if cluster.Addresses.HTTPS == "" {
		cluster.Addresses.HTTPS = "127.0.0.1"
	}
	if cluster.Addresses.GRPC == "" {
		cluster.Addresses.GRPC = "127.0.0.1"
	}
	if cluster.Ports.DNS == 0 {
		cluster.Ports.DNS = 8600
	}
	if cluster.Ports.HTTP == 0 {
		cluster.Ports.HTTP = 8500
	}
	if cluster.Ports.HTTPS == 0 {
		cluster.Ports.HTTPS = 8501
	}
	if cluster.Ports.GRPC == 0 {
		cluster.Ports.GRPC = 8502
	}
	if cluster.Ports.SerfLAN == 0 {
		cluster.Ports.SerfLAN = 8301
	}
	if cluster.Ports.SerfWAN == 0 {
		cluster.Ports.SerfWAN = 8302
	}
	if cluster.Ports.Server == 0 {
		cluster.Ports.Server = 8300
	}
}

func (a API) handleGetConsulCluster(w http.ResponseWriter, r *http.Request) {
	cluster, err := a.Storer.GetConsulCluster(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrConsulClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{ConsulClusters: []ConsulCluster{cluster}})
}

func (a API) handlePostConsulCluster(w http.ResponseWriter, r *http.Request) {
	var cluster ConsulCluster
	err := api.Decode(r, &cluster)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
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
	err = a.Storer.CreateConsulCluster(cluster)
	if err != nil {
		if err == ErrConsulClusterAlreadyExists {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusCreated, Response{ConsulClusters: []ConsulCluster{cluster}})
}

func (a API) handlePutConsulCluster(w http.ResponseWriter, r *http.Request) {
	var cluster ConsulCluster
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
	err = a.Storer.UpdateConsulCluster(cluster)
	if err != nil {
		if err == ErrConsulClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{ConsulClusters: []ConsulCluster{cluster}})
}

func (a API) handleDeleteConsulCluster(w http.ResponseWriter, r *http.Request) {
	cluster, err := a.Storer.GetConsulCluster(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrConsulClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	err = a.Storer.DeleteConsulCluster(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrConsulClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{ConsulClusters: []ConsulCluster{cluster}})
}
