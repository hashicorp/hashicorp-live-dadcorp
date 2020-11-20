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
	ErrConsulClusterNotFound = errors.New("consul cluster not found")
)

type ConsulService struct {
	basePath string
	client   *Client
	Clusters *ConsulClustersService
}

func newConsulService(basePath string, client *Client) *ConsulService {
	s := &ConsulService{
		basePath: basePath,
		client:   client,
	}
	s.Clusters = newConsulClustersService("clusters", s)
	return s
}

type ConsulClustersService struct {
	consulService *ConsulService
	basePath      string
}

func newConsulClustersService(basePath string, consul *ConsulService) *ConsulClustersService {
	return &ConsulClustersService{
		basePath:      basePath,
		consulService: consul,
	}
}

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

func (c ConsulClustersService) buildURL(p string) string {
	return path.Join(c.consulService.basePath, c.basePath, p)
}

func (c ConsulClustersService) Create(ctx context.Context, cluster ConsulCluster) (ConsulCluster, error) {
	b, err := json.Marshal(cluster)
	if err != nil {
		return ConsulCluster{}, fmt.Errorf("error serialising cluster: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := c.consulService.client.NewRequest(ctx, http.MethodPost, c.buildURL("/"), buf)
	if err != nil {
		return ConsulCluster{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := c.consulService.client.Do(req)
	if err != nil {
		return ConsulCluster{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return ConsulCluster{}, err
	}

	if resp.Errors.Contains(serverError) {
		return ConsulCluster{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return ConsulCluster{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return ConsulCluster{}, errors.New("cluster must have a name")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrConflict,
		Field: "/id",
	}) {
		return ConsulCluster{}, errors.New("cluster already exists")
	}
	if len(resp.Errors) > 0 {
		return ConsulCluster{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.ConsulClusters) < 1 {
		return ConsulCluster{}, errors.New("no Consul cluster returned in response")
	}
	return resp.ConsulClusters[0], nil
}

func (c ConsulClustersService) Get(ctx context.Context, id string) (ConsulCluster, error) {
	if id == "" {
		return ConsulCluster{}, errors.New("id must be specified")
	}
	req, err := c.consulService.client.NewRequest(ctx, http.MethodGet, c.buildURL("/"+id), nil)
	if err != nil {
		return ConsulCluster{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := c.consulService.client.Do(req)
	if err != nil {
		return ConsulCluster{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return ConsulCluster{}, err
	}

	if resp.Errors.Contains(serverError) {
		return ConsulCluster{}, errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return ConsulCluster{}, ErrConsulClusterNotFound
	}
	if len(resp.Errors) > 0 {
		return ConsulCluster{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.ConsulClusters) < 1 {
		return ConsulCluster{}, errors.New("no Consul cluster returned in response")
	}
	return resp.ConsulClusters[0], nil
}

func (c ConsulClustersService) Update(ctx context.Context, cluster ConsulCluster) (ConsulCluster, error) {
	if cluster.ID == "" {
		return ConsulCluster{}, errors.New("id must be specified")
	}
	b, err := json.Marshal(cluster)
	if err != nil {
		return ConsulCluster{}, fmt.Errorf("error serialising cluster: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := c.consulService.client.NewRequest(ctx, http.MethodPut, c.buildURL("/"+cluster.ID), buf)
	if err != nil {
		return ConsulCluster{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := c.consulService.client.Do(req)
	if err != nil {
		return ConsulCluster{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return ConsulCluster{}, err
	}

	if resp.Errors.Contains(serverError) {
		return ConsulCluster{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return ConsulCluster{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return ConsulCluster{}, ErrConsulClusterNotFound
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return ConsulCluster{}, errors.New("cluster must have a name")
	}
	if len(resp.Errors) > 0 {
		return ConsulCluster{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.ConsulClusters) < 1 {
		return ConsulCluster{}, errors.New("no Consul cluster returned in response")
	}
	return resp.ConsulClusters[0], nil
}

func (c ConsulClustersService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id must be specified")
	}
	req, err := c.consulService.client.NewRequest(ctx, http.MethodGet, c.buildURL("/"+id), nil)
	if err != nil {
		return fmt.Errorf("error constructing request: %w", err)
	}
	res, err := c.consulService.client.Do(req)
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
		return ErrConsulClusterNotFound
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	return nil
}
