package resource

import (
	"errors"
	"log"

	"github.com/digitalocean/concourse-resource-library/artifactory"
)

// Get performs the get operation for the resource
func Get(req GetRequest, dir string) (GetResponse, error) {
	c, err := artifactory.NewClient(
		artifactory.Endpoint(req.Source.Endpoint),
		artifactory.Authentication(req.Source.User, req.Source.Password, req.Source.APIKey, req.Source.AccessToken),
	)
	if err != nil {
		log.Println(err)
		return GetResponse{}, err
	}

	log.Println("version pattern:", req.Version.Pattern())

	artifacts, err := c.DownloadItems(req.Version.Pattern(), dir)
	if err != nil {
		log.Println(err)
		return GetResponse{}, err
	}

	if len(artifacts) == 0 {
		err := errors.New("no artifacts found")
		log.Println(err)
		return GetResponse{}, err
	}

	res := GetResponse{
		Version:  req.Version,
		Metadata: metadata(artifacts[0]),
	}

	return res, nil
}
