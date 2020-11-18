package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
	"github.com/hashicorp/go-uuid"
)

type AccessPolicy struct {
	ID         string      `json:"id"`
	PolicyData interface{} `json:"policyData"`
}

type TerraformPolicy struct {
	WorkspaceID      string `json:"id"`
	Plan             bool   `json:"plan"`
	Apply            bool   `json:"apply"`
	OverridePolicies bool   `json:"overridePolicies"`
}

type VaultPolicy struct {
	ClusterID string `json:"id"`
	Key       string `json:"key,omitempty"`
	Read      bool   `json:"read"`
	Write     bool   `json:"write"`
	Delete    bool   `json:"delete"`
}

type NomadPolicy struct {
	ClusterID     string `json:"id"`
	SubmitJobs    bool   `json:"submitJobs"`
	ReadJobStatus bool   `json:"readJobStatus"`
	CancelJobs    bool   `json:"cancelJobs"`
}

type ConsulPolicy struct {
	ClusterID string `json:"id"`
	Key       string `json:"key,omitempty"`
	Read      bool   `json:"read"`
	Write     bool   `json:"write"`
	Delete    bool   `json:"delete"`
}

func (a API) handleGetAccessPolicy(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetAccessPolicy(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrAccessPolicyNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{AccessPolicies: []AccessPolicy{ap}})
}

func (a API) handlePostAccessPolicy(w http.ResponseWriter, r *http.Request) {
	var ap AccessPolicy
	err := api.Decode(r, &ap)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	if ap.ID == "" {
		ap.ID, err = uuid.GenerateUUID()
		if err != nil {
			api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
			return
		}
	}
	if ap.PolicyData == nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData", Slug: api.RequestErrMissing}}})
		return
	}
	switch pd := ap.PolicyData.(type) {
	case TerraformPolicy:
		if pd.WorkspaceID == "" {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
	case VaultPolicy:
		if pd.ClusterID == "" {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
	case NomadPolicy:
		if pd.ClusterID == "" {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
	case ConsulPolicy:
		if pd.ClusterID == "" {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
	default:
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData", Slug: api.RequestErrInvalidValue}}})
		return
	}
	err = a.Storer.CreateAccessPolicy(ap)
	if err != nil {
		if err == ErrAccessPolicyAlreadyExists {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusCreated, Response{AccessPolicies: []AccessPolicy{ap}})
}

func (a API) handlePutAccessPolicy(w http.ResponseWriter, r *http.Request) {
	var ap AccessPolicy
	err := api.Decode(r, &ap)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	if ap.ID != "" && ap.ID != trout.RequestVars(r).Get("id") {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
		return
	}
	if ap.PolicyData == nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData", Slug: api.RequestErrMissing}}})
		return
	}
	switch pd := ap.PolicyData.(type) {
	case TerraformPolicy:
		if pd.WorkspaceID == "" {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
	case VaultPolicy:
		if pd.ClusterID == "" {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
	case NomadPolicy:
		if pd.ClusterID == "" {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
	case ConsulPolicy:
		if pd.ClusterID == "" {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
	default:
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData", Slug: api.RequestErrInvalidValue}}})
		return
	}
	err = a.Storer.UpdateAccessPolicy(ap)
	if err != nil {
		if err == ErrAccessPolicyNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{AccessPolicies: []AccessPolicy{ap}})
}

func (a API) handleDeleteAccessPolicy(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetAccessPolicy(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrAccessPolicyNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	err = a.Storer.DeleteAccessPolicy(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrAccessPolicyNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{AccessPolicies: []AccessPolicy{ap}})
}
