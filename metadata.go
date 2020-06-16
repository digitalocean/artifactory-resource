package resource

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/digitalocean/concourse-resource-library/artifactory"
	meta "github.com/digitalocean/concourse-resource-library/metadata"
)

func metadata(a artifactory.Artifact) meta.Metadata {
	var m meta.Metadata

	m.Add("artifactory-path", a.File.ArtifactoryPath)

	// trim `/tmp/build/get/` as that path is only available during the `get` step
	m.Add("local-path", strings.TrimPrefix(a.File.LocalPath, "/tmp/build/get/"))
	if a.File.FileHashes != nil {
		m.Add("sha1", a.File.Sha1)
	}

	m.Add("created", a.Item.Created)
	m.Add("modified", a.Item.Modified)
	m.Add("name", a.Item.Name)
	m.Add("repo", a.Item.Repo)
	m.Add("size", strconv.FormatInt(a.Item.Size, 10))
	m.Add("type", a.Item.Type)

	props, _ := json.Marshal(a.Item.Properties)
	m.AddJSON("properties", string(props))

	return m
}
