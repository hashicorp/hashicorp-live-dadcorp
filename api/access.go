package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
	"github.com/hashicorp/go-uuid"
)

type AccessPolicy struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
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
	msi, ok := ap.PolicyData.(map[string]interface{})
	if !ok {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData", Slug: api.RequestErrInvalidFormat}}})
		return
	}
	switch ap.Type {
	case "terraform":
		var tf TerraformPolicy
		workspaceID, ok := msi["id"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
		tf.WorkspaceID, ok = workspaceID.(string)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		plan, ok := msi["plan"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/plan", Slug: api.RequestErrMissing}}})
			return
		}
		tf.Plan, ok = plan.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/plan", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		apply, ok := msi["apply"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/apply", Slug: api.RequestErrMissing}}})
			return
		}
		tf.Apply, ok = apply.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/apply", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		overridePolicies, ok := msi["overridePolicies"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/overridePolicies", Slug: api.RequestErrMissing}}})
			return
		}
		tf.OverridePolicies, ok = overridePolicies.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/overridePolicies", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		ap.PolicyData = tf
	case "vault":
		var vault VaultPolicy
		clusterID, ok := msi["id"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
		vault.ClusterID, ok = clusterID.(string)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		key, ok := msi["key"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/key", Slug: api.RequestErrMissing}}})
			return
		}
		vault.Key, ok = key.(string)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/key", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		read, ok := msi["read"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/read", Slug: api.RequestErrMissing}}})
			return
		}
		vault.Read, ok = read.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/read", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		write, ok := msi["write"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/write", Slug: api.RequestErrMissing}}})
			return
		}
		vault.Write, ok = write.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/write", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		del, ok := msi["delete"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/delete", Slug: api.RequestErrMissing}}})
			return
		}
		vault.Delete, ok = del.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/delete", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		ap.PolicyData = vault
	case "nomad":
		var nomad NomadPolicy
		clusterID, ok := msi["id"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
		nomad.ClusterID, ok = clusterID.(string)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		submit, ok := msi["submitJobs"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/submitJobs", Slug: api.RequestErrMissing}}})
			return
		}
		nomad.SubmitJobs, ok = submit.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/submitJobs", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		read, ok := msi["readJobStatus"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/readJobStatus", Slug: api.RequestErrMissing}}})
			return
		}
		nomad.ReadJobStatus, ok = read.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/readJobStatus", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		cancel, ok := msi["cancelJobs"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/cancelJobs", Slug: api.RequestErrMissing}}})
			return
		}
		nomad.CancelJobs, ok = cancel.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/cancelJobs", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		ap.PolicyData = nomad
	case "consul":
		var consul ConsulPolicy
		clusterID, ok := msi["id"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
		consul.ClusterID, ok = clusterID.(string)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		key, ok := msi["key"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/key", Slug: api.RequestErrMissing}}})
			return
		}
		consul.Key, ok = key.(string)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/key", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		read, ok := msi["read"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/read", Slug: api.RequestErrMissing}}})
			return
		}
		consul.Read, ok = read.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/read", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		write, ok := msi["write"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/write", Slug: api.RequestErrMissing}}})
			return
		}
		consul.Write, ok = write.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/write", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		del, ok := msi["delete"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/delete", Slug: api.RequestErrMissing}}})
			return
		}
		consul.Delete, ok = del.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/delete", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		ap.PolicyData = consul
	default:
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/type", Slug: api.RequestErrInvalidValue}}})
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
	msi, ok := ap.PolicyData.(map[string]interface{})
	if !ok {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData", Slug: api.RequestErrInvalidFormat}}})
		return
	}
	switch ap.Type {
	case "terraform":
		var tf TerraformPolicy
		workspaceID, ok := msi["id"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
		tf.WorkspaceID, ok = workspaceID.(string)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		plan, ok := msi["plan"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/plan", Slug: api.RequestErrMissing}}})
			return
		}
		tf.Plan, ok = plan.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/plan", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		apply, ok := msi["apply"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/apply", Slug: api.RequestErrMissing}}})
			return
		}
		tf.Apply, ok = apply.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/apply", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		overridePolicies, ok := msi["overridePolicies"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/overridePolicies", Slug: api.RequestErrMissing}}})
			return
		}
		tf.OverridePolicies, ok = overridePolicies.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/overridePolicies", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		ap.PolicyData = tf
	case "vault":
		var vault VaultPolicy
		clusterID, ok := msi["id"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
		vault.ClusterID, ok = clusterID.(string)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		key, ok := msi["key"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/key", Slug: api.RequestErrMissing}}})
			return
		}
		vault.Key, ok = key.(string)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/key", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		read, ok := msi["read"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/read", Slug: api.RequestErrMissing}}})
			return
		}
		vault.Read, ok = read.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/read", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		write, ok := msi["write"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/write", Slug: api.RequestErrMissing}}})
			return
		}
		vault.Write, ok = write.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/write", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		del, ok := msi["delete"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/delete", Slug: api.RequestErrMissing}}})
			return
		}
		vault.Delete, ok = del.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/delete", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		ap.PolicyData = vault
	case "nomad":
		var nomad NomadPolicy
		clusterID, ok := msi["id"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
		nomad.ClusterID, ok = clusterID.(string)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		submit, ok := msi["submitJobs"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/submitJobs", Slug: api.RequestErrMissing}}})
			return
		}
		nomad.SubmitJobs, ok = submit.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/submitJobs", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		read, ok := msi["readJobStatus"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/readJobStatus", Slug: api.RequestErrMissing}}})
			return
		}
		nomad.ReadJobStatus, ok = read.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/readJobStatus", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		cancel, ok := msi["cancelJobs"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/cancelJobs", Slug: api.RequestErrMissing}}})
			return
		}
		nomad.CancelJobs, ok = cancel.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/cancelJobs", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		ap.PolicyData = nomad
	case "consul":
		var consul ConsulPolicy
		clusterID, ok := msi["id"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrMissing}}})
			return
		}
		consul.ClusterID, ok = clusterID.(string)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/id", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		key, ok := msi["key"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/key", Slug: api.RequestErrMissing}}})
			return
		}
		consul.Key, ok = key.(string)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/key", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		read, ok := msi["read"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/read", Slug: api.RequestErrMissing}}})
			return
		}
		consul.Read, ok = read.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/read", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		write, ok := msi["write"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/write", Slug: api.RequestErrMissing}}})
			return
		}
		consul.Write, ok = write.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/write", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		del, ok := msi["delete"]
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/delete", Slug: api.RequestErrMissing}}})
			return
		}
		consul.Delete, ok = del.(bool)
		if !ok {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/policyData/delete", Slug: api.RequestErrInvalidFormat}}})
			return
		}
		ap.PolicyData = consul
	default:
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/type", Slug: api.RequestErrInvalidValue}}})
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
