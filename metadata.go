package resource

import (
	"encoding/json"

	"github.com/digitalocean/concourse-resource-library/artifactory"
	meta "github.com/digitalocean/concourse-resource-library/metadata"
)

func metadata(a artifactory.Artifact) meta.Metadata {
	var m meta.Metadata

	m.Add("ArtifactoryPath", a.File.ArtifactoryPath)
	m.Add("LocalPath", a.File.LocalPath)
	if a.File.FileHashes != nil {
		m.Add("SHA1", a.File.Sha1)
		m.Add("SHA256", a.File.Sha256)
	}

	m.Add("Created", a.Item.Created)
	m.Add("Modified", a.Item.Modified)
	m.Add("Name", a.Item.Name)
	m.Add("Repo", a.Item.Repo)
	m.Add("Size", a.Item.Size)
	m.Add("Type", a.Item.Type)

	props, _ := json.Marshal(a.Item.Properties)
	m.Add("Properties", string(props))

	return m
}
