package dadcorp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
)

type RegionsService struct {
	basePath string
	client   *Client
}

func newRegionsService(basePath string, client *Client) *RegionsService {
	return &RegionsService{
		basePath: basePath,
		client:   client,
	}
}

type Region struct {
	ID       string         `json:"id"`
	Products RegionProducts `json:"products"`
}

type RegionProducts struct {
	Vault     bool `json:"vault"`
	Terraform bool `json:"terraform"`
	Nomad     bool `json:"nomad"`
	Consul    bool `json:"consul"`
}

func (r RegionsService) buildURL(p string) string {
	return path.Join(r.basePath, p)
}

func (r RegionsService) List(ctx context.Context) ([]Region, error) {
	req, err := r.client.NewRequest(ctx, http.MethodGet, r.buildURL("/"), nil)
	if err != nil {
		return nil, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return nil, err
	}

	if resp.Errors.Contains(serverError) {
		return nil, errors.New("server error")
	}
	if len(resp.Errors) > 0 {
		return nil, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	return resp.Regions, nil
}
