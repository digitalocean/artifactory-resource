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

	log.Println("working directory:", dir)
	log.Printf("put parameters: %+v", req.Params)

	b := buildInfo(req.Params, dir)

	props := properties(b)
	if req.Params.Properties != "" {
		err = props.FromFile(filepath.Join(dir, req.Params.Properties))
		if err != nil {
			rlog.StdErr("failed to read properties file", err)
		}
	}

	pattern := filepath.Join(dir, req.Params.Pattern)
	rlog.StdErr("pattern", pattern)
	rlog.StdErr("artifact properties", props)

	artifacts, uploaded, err := c.UploadItems(pattern, req.Params.Target, props)
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

	mod := buildinfo.Module{Id: moduleID(req.Params.Module, b.Name), Artifacts: []buildinfo.Artifact{}}

	for _, a := range artifacts {
		mod.Artifacts = append(mod.Artifacts, a.ToBuildArtifacts())
		rlog.StdErr("artifact uploaded", a)
	}

	b.Modules = []buildinfo.Module{mod}

	err = c.PublishBuildInfo(b)
	if err != nil {
		log.Println(err)
		return get, err
	}
	rlog.StdErr("build published", []string{b.Name, b.Number})

	return get, nil
}

func properties(b buildinfo.BuildInfo) artifactory.Properties {
	props := artifactory.Properties{
		artifactory.Property{Name: "build.name", Value: b.Name},
		artifactory.Property{Name: "build.number", Value: b.Number},
		artifactory.Property{Name: "build.started", Value: b.Started},
	}

	if b.Vcs != nil {
		props = append(props, artifactory.Property{Name: "vcs.revision", Value: b.Vcs.Revision})
		props = append(props, artifactory.Property{Name: "vcs.url", Value: b.Vcs.Url})
	}

	return props
}

func buildInfo(params PutParameters, dir string) buildinfo.BuildInfo {
	b := buildinfo.BuildInfo{
		Name:       os.Getenv("BUILD_TEAM_NAME") + "-" + os.Getenv("BUILD_PIPELINE_NAME") + "-" + os.Getenv("BUILD_JOB_NAME"),
		Number:     os.Getenv("BUILD_ID"),
		Started:    time.Now().Format("2006-01-02T15:04:05.000-0700"),
		Agent:      &buildinfo.Agent{Name: "Concourse"},
		BuildAgent: &buildinfo.Agent{Name: "digitalocean/artifactory-resource"},
		BuildUrl:   os.Getenv("ATC_EXTERNAL_URL") + "/builds/" + os.Getenv("BUILD_ID"),
	}

	if params.BuildEnv != "" {
		p := artifactory.Properties{}
		err := p.FromFile(filepath.Join(dir, params.BuildEnv))
		if err != nil {
			rlog.StdErr("failed to read build environment file", err)
		}

		log.Println("build env:", p)

		b.Properties = p.Env()
	}

	if params.RepositoryPath != "" {
		b.Vcs = vcsInfo(filepath.Join(dir, params.RepositoryPath), params.Repository)
	}

	return b
}

func vcsInfo(path, repo string) *buildinfo.Vcs {
	vcs := buildinfo.Vcs{}

	log.Println("vcs path:", path)

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

func moduleID(m, b string) string {
	if m != "" {
		return m
	}

	return b
}
