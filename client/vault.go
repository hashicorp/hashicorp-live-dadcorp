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
	ErrVaultClusterNotFound           = errors.New("vault cluster not found")
	ErrVaultClusterRegionNotFound     = errors.New("vault cluster region not found")
	ErrVaultClusterRegionAccessDenied = errors.New("authenticated user doesn't have the ability to provision Vault clusters in that region")
)

type VaultService struct {
	basePath string
	client   *Client
	Clusters *VaultClustersService
}

func newVaultService(basePath string, client *Client) *VaultService {
	s := &VaultService{
		basePath: basePath,
		client:   client,
	}
	s.Clusters = newVaultClustersService("clusters", s)
	return s
}

type VaultClustersService struct {
	vaultService *VaultService
	basePath     string
}

func newVaultClustersService(basePath string, vault *VaultService) *VaultClustersService {
	return &VaultClustersService{
		basePath:     basePath,
		vaultService: vault,
	}
}

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

func (v VaultClustersService) buildURL(p string) string {
	return path.Join(v.vaultService.basePath, v.basePath, p)
}

func (v VaultClustersService) Create(ctx context.Context, cluster VaultCluster) (VaultCluster, error) {
	b, err := json.Marshal(cluster)
	if err != nil {
		return VaultCluster{}, fmt.Errorf("error serialising cluster: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := v.vaultService.client.NewRequest(ctx, http.MethodPost, v.buildURL("/"), buf)
	if err != nil {
		return VaultCluster{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := v.vaultService.client.Do(req)
	if err != nil {
		return VaultCluster{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return VaultCluster{}, err
	}

	if resp.Errors.Contains(serverError) {
		return VaultCluster{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return VaultCluster{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return VaultCluster{}, errors.New("cluster must have a name")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrAccessDenied,
		Field: "/region",
	}) {
		return VaultCluster{}, ErrVaultClusterRegionAccessDenied
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrInvalidValue,
		Field: "/region",
	}) {
		return VaultCluster{}, ErrVaultClusterRegionNotFound
	}
	if len(resp.Errors) > 0 {
		return VaultCluster{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.VaultClusters) < 1 {
		return VaultCluster{}, errors.New("no Vault cluster returned in response")
	}
	return resp.VaultClusters[0], nil
}

func (v VaultClustersService) Get(ctx context.Context, id string) (VaultCluster, error) {
	if id == "" {
		return VaultCluster{}, errors.New("id must be specified")
	}
	req, err := v.vaultService.client.NewRequest(ctx, http.MethodGet, v.buildURL("/"+id), nil)
	if err != nil {
		return VaultCluster{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := v.vaultService.client.Do(req)
	if err != nil {
		return VaultCluster{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return VaultCluster{}, err
	}

	if resp.Errors.Contains(serverError) {
		return VaultCluster{}, errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return VaultCluster{}, ErrVaultClusterNotFound
	}
	if len(resp.Errors) > 0 {
		return VaultCluster{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.VaultClusters) < 1 {
		return VaultCluster{}, errors.New("no Vault cluster returned in response")
	}
	return resp.VaultClusters[0], nil
}

func (v VaultClustersService) Update(ctx context.Context, cluster VaultCluster) (VaultCluster, error) {
	if cluster.ID == "" {
		return VaultCluster{}, errors.New("id must be specified")
	}
	b, err := json.Marshal(cluster)
	if err != nil {
		return VaultCluster{}, fmt.Errorf("error serialising cluster: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := v.vaultService.client.NewRequest(ctx, http.MethodPut, v.buildURL("/"+cluster.ID), buf)
	if err != nil {
		return VaultCluster{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := v.vaultService.client.Do(req)
	if err != nil {
		return VaultCluster{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return VaultCluster{}, err
	}

	if resp.Errors.Contains(serverError) {
		return VaultCluster{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return VaultCluster{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return VaultCluster{}, ErrVaultClusterNotFound
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return VaultCluster{}, errors.New("cluster must have a name")
	}
	if len(resp.Errors) > 0 {
		return VaultCluster{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.VaultClusters) < 1 {
		return VaultCluster{}, errors.New("no Vault cluster returned in response")
	}
	return resp.VaultClusters[0], nil
}

func (v VaultClustersService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id must be specified")
	}
	req, err := v.vaultService.client.NewRequest(ctx, http.MethodGet, v.buildURL("/"+id), nil)
	if err != nil {
		return fmt.Errorf("error constructing request: %w", err)
	}
	res, err := v.vaultService.client.Do(req)
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
		return ErrVaultClusterNotFound
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	return nil
}
