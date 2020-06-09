package resource

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	m "github.com/digitalocean/concourse-resource-library/metadata"
)

// Source represents the configuration for the resource
type Source struct {
	Endpoint    string `json:"endpoint"`           // Endpoint for Artifactory AQL API (leave blank for cloud)
	User        string `json:"user,omitempty"`     // User for Artifactory API with permissions to Repository
	Password    string `json:"password,omitempty"` // Password for Artifactory API with permissions to Repository
	AccessToken string `json:"access_token"`       // AccessToken for Artifactory API with permissions to Repository
	APIKey      string `json:"api_key,omitempty"`  // APIKey for Artifactory API with permissions to Repository
	AQL         string `json:"aql"`                // AQL to filter versions on
}

// Validate ensures that the source configuration is valid
func (s Source) Validate() error {
	return nil
}

// Version contains the version data Concourse uses to determine if a build should run
type Version struct {
	Repo     string    `json:"repo"`
	Path     string    `json:"path"`
	Name     string    `json:"name"`
	Modified time.Time `json:"modified"`
}

// Pattern returns the string needed to fetch the artifact
func (v *Version) Pattern() string {
	return fmt.Sprintf("%s/%s/%s", v.Repo, v.Path, v.Name)
}

// Empty returns true if the version is empty
func (v *Version) Empty() bool {
	if v.Repo == "" || v.Path == "" {
		return true
	}

	return false
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
	SkipDownload bool `json:"skip_download"` // SkipDownload is used with `put` steps to skip `get` step that Concourse runs by default
}

// GetRequest is the data struct received from Concoruse by the resource get operation
type GetRequest struct {
	Source  Source        `json:"source"`
	Version Version       `json:"version"`
	Params  GetParameters `json:"params"`
}

// Read will read the json response from Concourse via stdin
func (r *GetRequest) Read(input []byte) error {
	return json.Unmarshal(input, r)
}

// GetResponse ...
type GetResponse struct {
	Version  Version    `json:"version"`
	Metadata m.Metadata `json:"metadata,omitempty"`
}

// Write will write the json response to stdout for Concourse to parse
func (r GetResponse) Write() error {
	return json.NewEncoder(os.Stdout).Encode(r)
}

// PutParameters for the resource
type PutParameters struct {
	Pattern        string        `json:"pattern"`               // Pattern to find artifacts within inputs
	Target         string        `json:"target"`                // Target to upload artifacts too
	Module         string        `json:"module,omitempty"`      // Module ID to associate the artifacts of the build to
	BuildEnv       string        `json:"build_env,omitempty"`   // BuildEnv is path to file containing build environment values in `key=value\n` form, e.g. `env > env.txt`
	EnvInclude     string        `json:"env_include,omitempty"` // EnvInclude case insensitive patterns in the form of "value1;value2;..." will be included
	EnvExclude     string        `json:"env_exclude,omitempty"` // EnvExclude case insensitive patterns in the form of "value1;value2;..." will be excluded, defaults to `*password*;*psw*;*secret*;*key*;*token*`
	Properties     string        `json:"properties,omitempty"`  // Properties is path to file containing artifact properties in `key=value\n` form
	MinimumUpload  int           `json:"min_upload,omitempty"`  // MinimumUpload sets the minimum number of uploads expected & will error if not met
	RepositoryPath string        `json:"repo_path,omitempty"`   // RepositoryPath sets the path to the input containing the repository (git support only)
	Repository     string        `json:"repo,omitempty"`        // Repository set the repository url explicitly for compatibility with the git resource
	Get            GetParameters `json:"get,omitempty"`         // Get parameters for explicit get step after put
}

// PutRequest is the data struct received from Concoruse by the resource put operation
type PutRequest struct {
	Source Source        `json:"source"`
	Params PutParameters `json:"params"`
}

// Read will read the json response from Concourse via stdin
func (r *PutRequest) Read(input []byte) error {
	return json.Unmarshal(input, r)
}
