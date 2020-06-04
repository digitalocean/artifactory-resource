package resource

import (
	"log"
	"os"

	"github.com/digitalocean/concourse-resource-library/artifactory"
	"github.com/digitalocean/concourse-resource-library/git"
	rlog "github.com/digitalocean/concourse-resource-library/log"
	meta "github.com/digitalocean/concourse-resource-library/metadata"
	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
)

// Put performs the Put operation for the resource
func Put(req PutRequest, dir string) (GetResponse, error) {
	get := GetResponse{
		Version:  Version{},
		Metadata: meta.Metadata{},
	}

	c, err := newClient(req.Source)
	if err != nil {
		log.Println(err)
		return get, err
	}

	log.Println(dir)

	b := buildInfo(req.Params.RepositoryPath)
	err = c.PublishBuildInfo(b)
	if err != nil {
		log.Println(err)
		return get, err
	}

	props := properties(b)
	artifacts, err := c.UploadItems(req.Params.Pattern, req.Params.Target, props)
	if err != nil {
		log.Println(err)
		return get, err
	}

	for _, a := range artifacts {
		rlog.StdErr("artifact uploaded", a)
	}

	return get, nil
}

func properties(b buildinfo.BuildInfo) artifactory.Properties {
	props := artifactory.Properties{
		artifactory.Property{
			Name:  "build.name",
			Value: b.Name,
		},
		artifactory.Property{Name: "build.number", Value: b.Number},
	}

	if b.Vcs != nil {
		props = append(props, artifactory.Property{Name: "vcs.revision", Value: b.Vcs.Revision})
		props = append(props, artifactory.Property{Name: "vcs.url", Value: b.Vcs.Url})
	}

	return props
}

func buildInfo(p string) buildinfo.BuildInfo {
	b := buildinfo.BuildInfo{
		Name:   os.Getenv("BUILD_TEAM_NAME") + "-" + os.Getenv("BUILD_PIPELINE_NAME") + "-" + os.Getenv("BUILD_JOB_NAME"),
		Number: os.Getenv("BUILD_ID"),
	}

	if p != "" {
		b.Vcs = vcsInfo(p)
	}

	return b
}

func vcsInfo(p string) *buildinfo.Vcs {
	vcs := buildinfo.Vcs{}

	g := git.Client{}
	r, err := g.Open(p)
	if err != nil {
		rlog.StdErr("failed to read vcs info", err)
		return &vcs
	}

	remotes, err := r.Remotes()
	if err != nil {
		rlog.StdErr("failed to read vcs info", err)
		return &vcs
	}

	if len(remotes) > 0 && len(remotes[0].Config().URLs) > 0 {
		url := remotes[0].Config().URLs[0]

		rev, err := r.Head()
		if err != nil {
			rlog.StdErr("failed to read vcs info", err)
			return &vcs
		}

		vcs.Url = url
		vcs.Revision = rev.Hash().String()
	}

	return &vcs
}
