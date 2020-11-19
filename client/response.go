package dadcorp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	requestErrAccessDenied  = "access_denied"
	requestErrInsufficient  = "insufficient"
	requestErrOverflow      = "overflow"
	requestErrInvalidValue  = "invalid_value"
	requestErrInvalidFormat = "invalid_format"
	requestErrMissing       = "missing"
	requestErrNotFound      = "not_found"
	requestErrConflict      = "conflict"
	requestErrActOfGod      = "act_of_god"
)

var (
	serverError        = RequestError{Slug: requestErrActOfGod}
	invalidFormatError = RequestError{Slug: requestErrInvalidFormat, Field: "/"}
)

type Response struct {
	Regions             []Region             `json:"regions,omitempty"`
	VaultClusters       []VaultCluster       `json:"vaultClusters,omitempty"`
	TerraformWorkspaces []TerraformWorkspace `json:"terraformWorkspaces,omitempty"`
	ConsulClusters      []ConsulCluster      `json:"consulClusters,omitempty"`
	NomadClusters       []NomadCluster       `json:"nomadClusters,omitempty"`
	AccessPolicies      []AccessPolicy       `json:"accessPolicies,omitempty"`
	Errors              RequestErrors        `json:"errors,omitempty"`
	Status              int                  `json:"-"`
}

func responseFromBody(resp *http.Response) (Response, error) {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("error reading response body: %w", err)
	}
	var res Response
	err = json.Unmarshal(b, &res)
	if err != nil {
		return Response{}, fmt.Errorf("error parsing response body: %w", err)
	}
	return res, nil
}

type RequestError struct {
	Slug   string `json:"error,omitempty"`
	Field  string `json:"field,omitempty"`
	Param  string `json:"param,omitempty"`
	Header string `json:"header,omitempty"`
}

func (e RequestError) Equal(other RequestError) bool {
	if e.Slug != other.Slug {
		return false
	}
	if e.Field != other.Field {
		return false
	}
	if e.Param != other.Param {
		return false
	}
	if e.Header != other.Header {
		return false
	}
	return true
}

type RequestErrors []RequestError

func (e RequestErrors) Contains(err RequestError) bool {
	for _, candidate := range e {
		if candidate.Equal(err) {
			return true
		}
	}
	return false
}

func (e RequestErrors) FieldMatches(slug string, re *regexp.Regexp) [][]string {
	for _, candidate := range e {
		if candidate.Slug != slug {
			continue
		}
		if re.MatchString(candidate.Field) {
			return re.FindAllStringSubmatch(candidate.Field, -1)
		}
	}
	return nil
}
