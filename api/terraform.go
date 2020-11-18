package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
	"github.com/hashicorp/go-uuid"
)

type TerraformWorkspace struct {
	ID                  string                    `json:"id"`
	Name                string                    `json:"name"`
	AgentPoolID         string                    `json:"agentPoolID"`
	AllowDestroyPlan    *bool                     `json:"allowDestroyPlan"`
	AutoApply           bool                      `json:"autoApply"`
	Description         string                    `json:"description"`
	ExecutionMode       string                    `json:"executionMode"`
	FileTriggersEnabled *bool                     `json:"fileTriggersEnabled"`
	SourceName          string                    `json:"sourceName"`
	SourceURL           string                    `json:"sourceURL"`
	QueueAllRuns        bool                      `json:"queueAllRuns"`
	SpeculativeEnabled  *bool                     `json:"speculativeEnabled"`
	TerraformVersion    string                    `json:"terraformVersion"`
	TriggerPrefixes     []string                  `json:"triggerPrefixes"`
	WorkingDirectory    string                    `json:"workingDirectory"`
	VCSRepo             TerraformWorkspaceVCSRepo `json:"vcsRepo"`
}

type TerraformWorkspaceVCSRepo struct {
	OAuthTokenID      string `json:"oauthTokenID"`
	Branch            string `json:"branch"`
	IngressSubmodules bool   `json:"ingressSubmodules"`
	Identifier        string `json:"identifier"`
}

func (workspace *TerraformWorkspace) FillDefaults() {
	if workspace.AllowDestroyPlan == nil {
		adp := true
		workspace.AllowDestroyPlan = &adp
	}
	if workspace.ExecutionMode == "" {
		workspace.ExecutionMode = "remote"
	}
	if workspace.FileTriggersEnabled == nil {
		fte := true
		workspace.FileTriggersEnabled = &fte
	}
	if workspace.SpeculativeEnabled == nil {
		se := true
		workspace.SpeculativeEnabled = &se
	}
	if workspace.TerraformVersion == "" {
		workspace.TerraformVersion = "0.13.5"
	}
	if workspace.TriggerPrefixes == nil {
		workspace.TriggerPrefixes = []string{}
	}
	if workspace.VCSRepo.Branch == "" {
		workspace.VCSRepo.Branch = "main"
	}
}

func (a API) handleGetTerraformWorkspace(w http.ResponseWriter, r *http.Request) {
	workspace, err := a.Storer.GetTerraformWorkspace(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrTerraformWorkspaceNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{TerraformWorkspaces: []TerraformWorkspace{workspace}})
}

func (a API) handlePostTerraformWorkspace(w http.ResponseWriter, r *http.Request) {
	var workspace TerraformWorkspace
	err := api.Decode(r, &workspace)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	if workspace.ID == "" {
		workspace.ID, err = uuid.GenerateUUID()
		if err != nil {
			api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
			return
		}
	}
	if workspace.Name == "" {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/name", Slug: api.RequestErrMissing}}})
		return
	}
	workspace.FillDefaults()
	err = a.Storer.CreateTerraformWorkspace(workspace)
	if err != nil {
		if err == ErrTerraformWorkspaceAlreadyExists {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusCreated, Response{TerraformWorkspaces: []TerraformWorkspace{workspace}})
}

func (a API) handlePutTerraformWorkspace(w http.ResponseWriter, r *http.Request) {
	var workspace TerraformWorkspace
	err := api.Decode(r, &workspace)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	if workspace.ID != "" && workspace.ID != trout.RequestVars(r).Get("id") {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
		return
	}
	if workspace.Name == "" {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/name", Slug: api.RequestErrMissing}}})
		return
	}
	workspace.FillDefaults()
	err = a.Storer.UpdateTerraformWorkspace(workspace)
	if err != nil {
		if err == ErrTerraformWorkspaceNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{TerraformWorkspaces: []TerraformWorkspace{workspace}})
}

func (a API) handleDeleteTerraformWorkspace(w http.ResponseWriter, r *http.Request) {
	workspace, err := a.Storer.GetTerraformWorkspace(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrTerraformWorkspaceNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	err = a.Storer.DeleteTerraformWorkspace(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrTerraformWorkspaceNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{TerraformWorkspaces: []TerraformWorkspace{workspace}})
}
