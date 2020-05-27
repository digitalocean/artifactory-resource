package resource

import (
	"log"

	"github.com/digitalocean/concourse-resource-library/artifactory"
)

// Check performs the check operation for the resource
func Check(req CheckRequest) (CheckResponse, error) {
	c, err := artifactory.NewClient(
		artifactory.Endpoint(req.Source.Endpoint),
		artifactory.Authentication(req.Source.User, req.Source.Password, req.Source.APIKey, req.Source.AccessToken),
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("query:", req.Source.AQL)

	data, err := c.AQL(req.Source.AQL)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(data)

	return nil, nil
}
