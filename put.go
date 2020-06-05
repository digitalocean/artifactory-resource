package resource

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

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

	b := buildInfo(filepath.Join(dir, req.Params.RepositoryPath), req.Params.Repository)
	err = c.PublishBuildInfo(b)
	if err != nil {
		log.Println(err)
		return get, err
	}
	rlog.StdErr("build published", []string{b.Name, b.Number})

	props := properties(b)

	rlog.StdErr("pattern", filepath.Join(dir, req.Params.Pattern))
	rlog.StdErr("target", req.Params.Target)
	rlog.StdErr("artifact properties", props)

	artifacts, uploaded, err := c.UploadItems(filepath.Join(dir, req.Params.Pattern), req.Params.Target, props)
	if err != nil {
		rlog.StdErr("failed to upload", err)
		log.Println(err)
		return get, err
	}

	if req.Params.MinimumUpload > uploaded {
		err = fmt.Errorf("failed to upload minimum (%v) count: uploaded %v artifacts", req.Params.MinimumUpload, uploaded)
		rlog.StdErr("failed to upload", err)
		log.Println(err)
		return get, err
	}

	rlog.StdErr("upload count", uploaded)

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

func buildInfo(path, repo string) buildinfo.BuildInfo {
	b := buildinfo.BuildInfo{
		Name:       os.Getenv("BUILD_TEAM_NAME") + "-" + os.Getenv("BUILD_PIPELINE_NAME") + "-" + os.Getenv("BUILD_JOB_NAME"),
		Number:     os.Getenv("BUILD_ID"),
		Started:    time.Now().Format("2006-01-02T15:04:05.000-0700"),
		BuildAgent: &buildinfo.Agent{Name: "Concourse"},
	}

	if path != "" {
		b.Vcs = vcsInfo(path, repo)
	}

	return b
}

func vcsInfo(path, repo string) *buildinfo.Vcs {
	vcs := buildinfo.Vcs{}

	g := git.Client{}
	r, err := g.Open(path)
	if err != nil {
		rlog.StdErr("failed to open repository", err)
		return &vcs
	}

	rev, err := r.Head()
	if err != nil {
		rlog.StdErr("failed to read vcs revision", err)
		return &vcs
	}
	vcs.Revision = rev.Hash().String()

	if repo != "" {
		vcs.Url = repo
		return &vcs
	}

	remotes, err := r.Remotes()
	if err != nil {
		rlog.StdErr("failed to read vcs info", err)
		return &vcs
	}

	if len(remotes) > 0 && len(remotes[0].Config().URLs) > 0 {
		url := remotes[0].Config().URLs[0]

		vcs.Url = url
	}

	return &vcs
}
