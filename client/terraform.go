package dadcorp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
)

var (
	ErrTerraformWorkspaceNotFound = errors.New("terraform workspace not found")
)

type TerraformService struct {
	basePath   string
	client     *Client
	Workspaces *TerraformWorkspacesService
}

func newTerraformService(basePath string, client *Client) *TerraformService {
	s := &TerraformService{
		basePath: basePath,
		client:   client,
	}
	s.Workspaces = newTerraformWorkspacesService("workspaces", s)
	return s
}

type TerraformWorkspacesService struct {
	terraformService *TerraformService
	basePath         string
}

func newTerraformWorkspacesService(basePath string, terraform *TerraformService) *TerraformWorkspacesService {
	return &TerraformWorkspacesService{
		basePath:         basePath,
		terraformService: terraform,
	}
}

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

func (t TerraformWorkspacesService) buildURL(p string) string {
	return path.Join(t.terraformService.basePath, t.basePath, p)
}

func (t TerraformWorkspacesService) Create(ctx context.Context, workspace TerraformWorkspace) (TerraformWorkspace, error) {
	b, err := json.Marshal(workspace)
	if err != nil {
		return TerraformWorkspace{}, fmt.Errorf("error serialising workspace: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := t.terraformService.client.NewRequest(ctx, http.MethodPost, t.buildURL("/"), buf)
	if err != nil {
		return TerraformWorkspace{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := t.terraformService.client.Do(req)
	if err != nil {
		return TerraformWorkspace{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return TerraformWorkspace{}, err
	}

	if resp.Errors.Contains(serverError) {
		return TerraformWorkspace{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return TerraformWorkspace{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return TerraformWorkspace{}, errors.New("workspace must have a name")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrConflict,
		Field: "/id",
	}) {
		return TerraformWorkspace{}, errors.New("workspace already exists")
	}
	if len(resp.Errors) > 0 {
		return TerraformWorkspace{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.TerraformWorkspaces) < 1 {
		return TerraformWorkspace{}, errors.New("no Terraform workspace returned in response")
	}
	return resp.TerraformWorkspaces[0], nil
}

func (t TerraformWorkspacesService) Get(ctx context.Context, id string) (TerraformWorkspace, error) {
	if id == "" {
		return TerraformWorkspace{}, errors.New("id must be specified")
	}
	req, err := t.terraformService.client.NewRequest(ctx, http.MethodGet, t.buildURL("/"+id), nil)
	if err != nil {
		return TerraformWorkspace{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := t.terraformService.client.Do(req)
	if err != nil {
		return TerraformWorkspace{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return TerraformWorkspace{}, err
	}

	if resp.Errors.Contains(serverError) {
		return TerraformWorkspace{}, errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return TerraformWorkspace{}, ErrTerraformWorkspaceNotFound
	}
	if len(resp.Errors) > 0 {
		return TerraformWorkspace{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.TerraformWorkspaces) < 1 {
		return TerraformWorkspace{}, errors.New("no Terraform workspace returned in response")
	}
	return resp.TerraformWorkspaces[0], nil
}

func (t TerraformWorkspacesService) Update(ctx context.Context, workspace TerraformWorkspace) (TerraformWorkspace, error) {
	if workspace.ID == "" {
		return TerraformWorkspace{}, errors.New("id must be specified")
	}
	b, err := json.Marshal(workspace)
	if err != nil {
		return TerraformWorkspace{}, fmt.Errorf("error serialising workspace: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := t.terraformService.client.NewRequest(ctx, http.MethodPut, t.buildURL("/"+workspace.ID), buf)
	if err != nil {
		return TerraformWorkspace{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := t.terraformService.client.Do(req)
	if err != nil {
		return TerraformWorkspace{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return TerraformWorkspace{}, err
	}

	if resp.Errors.Contains(serverError) {
		return TerraformWorkspace{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return TerraformWorkspace{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return TerraformWorkspace{}, ErrTerraformWorkspaceNotFound
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return TerraformWorkspace{}, errors.New("workspace must have a name")
	}
	if len(resp.Errors) > 0 {
		return TerraformWorkspace{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.TerraformWorkspaces) < 1 {
		return TerraformWorkspace{}, errors.New("no Terraform workspace returned in response")
	}
	return resp.TerraformWorkspaces[0], nil
}

func (t TerraformWorkspacesService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id must be specified")
	}
	req, err := t.terraformService.client.NewRequest(ctx, http.MethodGet, t.buildURL("/"+id), nil)
	if err != nil {
		return fmt.Errorf("error constructing request: %w", err)
	}
	res, err := t.terraformService.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return err
	}

	if resp.Errors.Contains(serverError) {
		return errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return ErrTerraformWorkspaceNotFound
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	return nil
}
