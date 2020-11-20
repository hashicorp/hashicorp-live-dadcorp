package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
	"github.com/hashicorp/go-uuid"
)

type NomadCluster struct {
	ID         string                `json:"id"`
	Name       string                `json:"name"`
	Datacenter string                `json:"datacenter"`
	BindAddr   string                `json:"bindAddr"`
	Advertise  NomadClusterAdvertise `json:"advertise"`
	Ports      NomadClusterPorts     `json:"ports"`
	Server     NomadClusterServer    `json:"server"`
}

type NomadClusterAdvertise struct {
	HTTP string `json:"http"`
	RPC  string `json:"rpc"`
	Serf string `json:"serf"`
}

type NomadClusterPorts struct {
	HTTP int `json:"http"`
	RPC  int `json:"rpc"`
	Serf int `json:"serf"`
}

type NomadClusterServer struct {
	ServerJoin NomadClusterServerServerJoin `json:"serverJoin"`
}

type NomadClusterServerServerJoin struct {
	RetryJoin     []string `json:"retryJoin"`
	RetryInterval string   `json:"retryInterval"`
	RetryMax      int      `json:"retryMax"`
	StartJoin     []string `json:"startJoin"`
}

func (cluster *NomadCluster) FillDefaults() {
	if cluster.Ports.HTTP == 0 {
		cluster.Ports.HTTP = 4646
	}
	if cluster.Ports.RPC == 0 {
		cluster.Ports.RPC = 4647
	}
	if cluster.Ports.Serf == 0 {
		cluster.Ports.Serf = 4648
	}
	if cluster.Server.ServerJoin.RetryJoin == nil {
		cluster.Server.ServerJoin.RetryJoin = []string{}
	}
	if cluster.Server.ServerJoin.StartJoin == nil {
		cluster.Server.ServerJoin.StartJoin = []string{}
	}
	if cluster.Server.ServerJoin.RetryInterval == "" {
		cluster.Server.ServerJoin.RetryInterval = "30s"
	}
}

func (a API) handleGetNomadCluster(w http.ResponseWriter, r *http.Request) {
	cluster, err := a.Storer.GetNomadCluster(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrNomadClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{NomadClusters: []NomadCluster{cluster}})
}

func (a API) handlePostNomadCluster(w http.ResponseWriter, r *http.Request) {
	var cluster NomadCluster
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
	if cluster.Datacenter == "" {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/datacenter", Slug: api.RequestErrMissing}}})
		return
	}
	cluster.FillDefaults()
	err = a.Storer.CreateNomadCluster(cluster)
	if err != nil {
		if err == ErrNomadClusterAlreadyExists {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusCreated, Response{NomadClusters: []NomadCluster{cluster}})
}

func (a API) handlePutNomadCluster(w http.ResponseWriter, r *http.Request) {
	var cluster NomadCluster
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
	if cluster.Datacenter == "" {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/datacenter", Slug: api.RequestErrMissing}}})
		return
	}
	cluster.FillDefaults()
	err = a.Storer.UpdateNomadCluster(cluster)
	if err != nil {
		if err == ErrNomadClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{NomadClusters: []NomadCluster{cluster}})
}

func (a API) handleDeleteNomadCluster(w http.ResponseWriter, r *http.Request) {
	cluster, err := a.Storer.GetNomadCluster(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrNomadClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	err = a.Storer.DeleteNomadCluster(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrNomadClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{NomadClusters: []NomadCluster{cluster}})
}
