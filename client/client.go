package dadcorp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-cleanhttp"
)

type Client struct {
	client  *http.Client
	baseURL *url.URL

	username, password string

	Terraform      *TerraformService
	Vault          *VaultService
	Nomad          *NomadService
	Consul         *ConsulService
	AccessPolicies *AccessPoliciesService
	Regions        *RegionsService
}

func NewClient(baseURL, username, password string) (*Client, error) {
	if username == "" {
		return nil, fmt.Errorf("username must be set")
	}
	if password == "" {
		return nil, fmt.Errorf("password must be set")
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	c := &Client{
		client:   cleanhttp.DefaultPooledClient(),
		baseURL:  base,
		username: username,
		password: password,
	}
	c.Terraform = newTerraformService("terraform", c)
	c.Vault = newVaultService("vault", c)
	c.Nomad = newNomadService("nomad", c)
	c.Consul = newConsulService("consul", c)
	c.Regions = newRegionsService("regions", c)
	c.AccessPolicies = newAccessPoliciesService("accessPolicies", c)
	return c, nil
}

func (c Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing path: %w", err)
	}
	reqURL := c.baseURL.ResolveReference(u)
	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
