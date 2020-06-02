package resource

import (
	"fmt"
	"log"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

// Check performs the check operation for the resource
func Check(req CheckRequest) (CheckResponse, error) {
	c, err := newClient(req.Source)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	aql := addModifiedTime(req.Source.AQL, req.Version)

	log.Println("query:", aql)

	data, err := c.SearchItems(aql)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(data)

	res, err := processItems(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res = selectVersions(req.Version, res)

	log.Println("version count in response:", len(res))
	log.Println("versions:", res)

	return res, nil
}

func processItems(s []utils.ResultItem) (CheckResponse, error) {
	var res CheckResponse

	for _, i := range s {
		v, err := processItem(i)
		if err != nil {
			return nil, err
		}

		res = append(res, v)
	}

	return res, nil
}

func processItem(i utils.ResultItem) (Version, error) {
	var v Version

	m, err := time.Parse(time.RFC3339, i.Modified)
	if err != nil {
		return v, err
	}

	v = Version{Repo: i.Repo, Path: i.Path, Name: i.Name, Modified: m}

	return v, nil
}

// selectVersions handles business logic based on input version
// 	from Concourse and versions found in external resource
func selectVersions(v Version, res CheckResponse) CheckResponse {
	// If there are no new but an input version, return the input
	if len(res) == 0 && v.Repo != "" {
		log.Println("no new versions, use input version")
		res = append(res, v)

	}

	// If there are new versions and no input version, return latest new version
	if len(res) != 0 && v.Repo == "" {
		log.Println("new versions but no input version, use latest")
		res = CheckResponse{res[len(res)-1]}
	}

	return res
}

func addModifiedTime(aql string, v Version) string {
	m := time.Now().AddDate(-2, 0, 0)

	if !v.Modified.IsZero() {
		m = v.Modified
	}

	if len(aql) < 1 {
		return ""
	}

	return fmt.Sprintf("%s, \"modified\": {\"$gt\": \"%s\"}}", aql[:len(aql)-1], m.Format(time.RFC3339))
}
