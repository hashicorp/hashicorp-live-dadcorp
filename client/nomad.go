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
	ErrNomadClusterNotFound = errors.New("nomad cluster not found")
)

type NomadService struct {
	basePath string
	client   *Client
	Clusters *NomadClustersService
}

func newNomadService(basePath string, client *Client) *NomadService {
	s := &NomadService{
		basePath: basePath,
		client:   client,
	}
	s.Clusters = newNomadClustersService("clusters", s)
	return s
}

type NomadClustersService struct {
	nomadService *NomadService
	basePath     string
}

func newNomadClustersService(basePath string, nomad *NomadService) *NomadClustersService {
	return &NomadClustersService{
		basePath:     basePath,
		nomadService: nomad,
	}
}

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

func (n NomadClustersService) buildURL(p string) string {
	return path.Join(n.nomadService.basePath, n.basePath, p)
}

func (n NomadClustersService) Create(ctx context.Context, cluster NomadCluster) (NomadCluster, error) {
	b, err := json.Marshal(cluster)
	if err != nil {
		return NomadCluster{}, fmt.Errorf("error serialising cluster: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := n.nomadService.client.NewRequest(ctx, http.MethodPost, n.buildURL("/"), buf)
	if err != nil {
		return NomadCluster{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := n.nomadService.client.Do(req)
	if err != nil {
		return NomadCluster{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return NomadCluster{}, err
	}

	if resp.Errors.Contains(serverError) {
		return NomadCluster{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return NomadCluster{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return NomadCluster{}, errors.New("cluster must have a name")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/datacenter",
	}) {
		return NomadCluster{}, errors.New("cluster must have a datacenter name")
	}
	if len(resp.Errors) > 0 {
		return NomadCluster{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.NomadClusters) < 1 {
		return NomadCluster{}, errors.New("no Nomad cluster returned in response")
	}
	return resp.NomadClusters[0], nil
}

func (n NomadClustersService) Get(ctx context.Context, id string) (NomadCluster, error) {
	if id == "" {
		return NomadCluster{}, errors.New("id must be specified")
	}
	req, err := n.nomadService.client.NewRequest(ctx, http.MethodGet, n.buildURL("/"+id), nil)
	if err != nil {
		return NomadCluster{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := n.nomadService.client.Do(req)
	if err != nil {
		return NomadCluster{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return NomadCluster{}, err
	}

	if resp.Errors.Contains(serverError) {
		return NomadCluster{}, errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return NomadCluster{}, ErrNomadClusterNotFound
	}
	if len(resp.Errors) > 0 {
		return NomadCluster{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.NomadClusters) < 1 {
		return NomadCluster{}, errors.New("no Nomad cluster returned in response")
	}
	return resp.NomadClusters[0], nil
}

func (n NomadClustersService) Update(ctx context.Context, cluster NomadCluster) (NomadCluster, error) {
	if cluster.ID == "" {
		return NomadCluster{}, errors.New("id must be specified")
	}
	b, err := json.Marshal(cluster)
	if err != nil {
		return NomadCluster{}, fmt.Errorf("error serialising cluster: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := n.nomadService.client.NewRequest(ctx, http.MethodPut, n.buildURL("/"+cluster.ID), buf)
	if err != nil {
		return NomadCluster{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := n.nomadService.client.Do(req)
	if err != nil {
		return NomadCluster{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return NomadCluster{}, err
	}

	if resp.Errors.Contains(serverError) {
		return NomadCluster{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return NomadCluster{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return NomadCluster{}, ErrNomadClusterNotFound
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return NomadCluster{}, errors.New("cluster must have a name")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/datacenter",
	}) {
		return NomadCluster{}, errors.New("cluster must have a datacenter name")
	}
	if len(resp.Errors) > 0 {
		return NomadCluster{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.NomadClusters) < 1 {
		return NomadCluster{}, errors.New("no Nomad cluster returned in response")
	}
	return resp.NomadClusters[0], nil
}

func (n NomadClustersService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id must be specified")
	}
	req, err := n.nomadService.client.NewRequest(ctx, http.MethodGet, n.buildURL("/"+id), nil)
	if err != nil {
		return fmt.Errorf("error constructing request: %w", err)
	}
	res, err := n.nomadService.client.Do(req)
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
		return ErrNomadClusterNotFound
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	return nil
}
