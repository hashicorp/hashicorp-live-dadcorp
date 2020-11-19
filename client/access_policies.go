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
	ErrAccessPolicyNotFound = errors.New("access policy not found")
)

type AccessPoliciesService struct {
	basePath string
	client   *Client
}

func newAccessPoliciesService(basePath string, client *Client) *AccessPoliciesService {
	return &AccessPoliciesService{
		basePath: basePath,
		client:   client,
	}
}

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

func (a AccessPoliciesService) buildURL(p string) string {
	return path.Join(a.basePath, p)
}

func (a AccessPoliciesService) Create(ctx context.Context, policy AccessPolicy) (AccessPolicy, error) {
	b, err := json.Marshal(policy)
	if err != nil {
		return AccessPolicy{}, fmt.Errorf("error serialising policy: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := a.client.NewRequest(ctx, http.MethodPost, a.buildURL("/"), buf)
	if err != nil {
		return AccessPolicy{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := a.client.Do(req)
	if err != nil {
		return AccessPolicy{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return AccessPolicy{}, err
	}

	if resp.Errors.Contains(serverError) {
		return AccessPolicy{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return AccessPolicy{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrConflict,
		Field: "/id",
	}) {
		return AccessPolicy{}, errors.New("policy already exists")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/policyData",
	}) {
		return AccessPolicy{}, errors.New("policy data must be set")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/policyData/id",
	}) {
		return AccessPolicy{}, errors.New("policy data ID must be set")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrInvalidValue,
		Field: "/policyData",
	}) {
		return AccessPolicy{}, errors.New("invalid policy data")
	}
	if len(resp.Errors) > 0 {
		return AccessPolicy{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.AccessPolicies) < 1 {
		return AccessPolicy{}, errors.New("no Terraform policy returned in response")
	}
	return resp.AccessPolicies[0], nil
}

func (a AccessPoliciesService) Get(ctx context.Context, id string) (AccessPolicy, error) {
	if id == "" {
		return AccessPolicy{}, errors.New("id must be specified")
	}
	req, err := a.client.NewRequest(ctx, http.MethodGet, a.buildURL("/"+id), nil)
	if err != nil {
		return AccessPolicy{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := a.client.Do(req)
	if err != nil {
		return AccessPolicy{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return AccessPolicy{}, err
	}

	if resp.Errors.Contains(serverError) {
		return AccessPolicy{}, errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return AccessPolicy{}, ErrAccessPolicyNotFound
	}
	if len(resp.Errors) > 0 {
		return AccessPolicy{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.AccessPolicies) < 1 {
		return AccessPolicy{}, errors.New("no Terraform policy returned in response")
	}
	return resp.AccessPolicies[0], nil
}

func (a AccessPoliciesService) Update(ctx context.Context, policy AccessPolicy) (AccessPolicy, error) {
	if policy.ID == "" {
		return AccessPolicy{}, errors.New("id must be specified")
	}
	b, err := json.Marshal(policy)
	if err != nil {
		return AccessPolicy{}, fmt.Errorf("error serialising policy: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := a.client.NewRequest(ctx, http.MethodPut, a.buildURL("/"+policy.ID), buf)
	if err != nil {
		return AccessPolicy{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := a.client.Do(req)
	if err != nil {
		return AccessPolicy{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return AccessPolicy{}, err
	}

	if resp.Errors.Contains(serverError) {
		return AccessPolicy{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return AccessPolicy{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return AccessPolicy{}, ErrAccessPolicyNotFound
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/policyData",
	}) {
		return AccessPolicy{}, errors.New("policy data must be set")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/policyData/id",
	}) {
		return AccessPolicy{}, errors.New("policy data ID must be set")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrInvalidValue,
		Field: "/policyData",
	}) {
		return AccessPolicy{}, errors.New("invalid policy data")
	}
	if len(resp.Errors) > 0 {
		return AccessPolicy{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.AccessPolicies) < 1 {
		return AccessPolicy{}, errors.New("no Terraform policy returned in response")
	}
	return resp.AccessPolicies[0], nil
}

func (a AccessPoliciesService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id must be specified")
	}
	req, err := a.client.NewRequest(ctx, http.MethodGet, a.buildURL("/"+id), nil)
	if err != nil {
		return fmt.Errorf("error constructing request: %w", err)
	}
	res, err := a.client.Do(req)
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
		return ErrAccessPolicyNotFound
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	return nil
}
