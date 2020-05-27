package resource

import (
	"encoding/json"
	"os"
	"time"

	m "github.com/digitalocean/concourse-resource-library/metadata"
)

// Source represents the configuration for the resource
type Source struct {
	// Endpoint for Artifactory AQL API (leave blank for cloud)
	Endpoint string `json:"endpoint"`
	// User for Artifactory API with permissions to Repository
	User string `json:"user,omitempty"`
	// Password for Artifactory API with permissions to Repository
	Password string `json:"password,omitempty"`
	// AccessToken for Artifactory API with permissions to Repository
	AccessToken string `json:"access_token"`
	// APIKey for Artifactory API with permissions to Repository
	APIKey string `json:"api_key,omitempty"`
	// AQL to filter versions on
	AQL string `json:"aql"`
}

// Validate ensures that the source configuration is valid
func (s Source) Validate() error {
	return nil
}

// Version contains the version data Concourse uses to determine if a build should run
type Version struct {
	UpdatedDate time.Time `json:"updated"`
}

// CheckRequest is the data struct received from Concoruse by the resource check operation
type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// Read will read the json response from Concourse via stdin
func (r *CheckRequest) Read(input []byte) error {
	return json.Unmarshal(input, r)
}

// CheckResponse is the data struct returned to Concourse by the resource check operation
type CheckResponse []Version

// Len returns the number of versions in the response
func (r CheckResponse) Len() int {
	return len(r)
}

// Write will write the json response to stdout for Concourse to parse
func (r CheckResponse) Write() error {
	return json.NewEncoder(os.Stdout).Encode(r)
}

// GetParameters is the configuration for a resource step
type GetParameters struct {
}

// GetRequest is the data struct received from Concoruse by the resource get operation
type GetRequest struct {
}

// GetResponse ...
type GetResponse struct {
	Version  Version    `json:"version"`
	Metadata m.Metadata `json:"metadata,omitempty"`
}

// PutParameters for the resource
type PutParameters struct {
}

// PutRequest is the data struct received from Concoruse by the resource put operation
type PutRequest struct {
}