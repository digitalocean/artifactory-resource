package resource

import (
	"log"
	"time"

	"github.com/digitalocean/concourse-resource-library/artifactory"
)

// Check performs the check operation for the resource
func Check(req CheckRequest) (CheckResponse, error) {
	var res CheckResponse

	c, err := artifactory.NewClient(
		artifactory.Endpoint(req.Source.Endpoint),
		artifactory.Authentication(req.Source.User, req.Source.Password, req.Source.APIKey, req.Source.AccessToken),
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("query:", req.Source.AQL)

	data, err := c.SearchItems(req.Source.AQL)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(data)

	for _, i := range data {
		m, err := time.Parse(time.RFC3339, i.Modified)
		if err != nil {
			return nil, err
		}

		res = append(res, Version{SHA1: i.Actual_Sha1, Modified: m})
	}

	// If there are no new but an old version = return the old
	if len(res) == 0 && req.Version.SHA1 != "" {
		log.Println("no new versions, use old")
		res = append(res, req.Version)
	}

	// If there are new versions and no previous = return just the latest
	if len(res) != 0 && req.Version.SHA1 == "" {
		res = CheckResponse{res[len(res)-1]}
	}

	log.Println("version count in response:", len(res))
	log.Println("versions:", res)

	return res, nil
}
