package resource

import (
	"errors"
	"log"
	"os"
)

// Get performs the get operation for the resource
func Get(req GetRequest, dir string) (GetResponse, error) {
	var res GetResponse

	c, err := newClient(req.Source)
	if err != nil {
		log.Println(err)
		return res, err
	}

	log.Println(dir)
	log.Println("version pattern:", req.Version.Pattern())

	artifacts, err := c.DownloadItems(req.Version.Pattern(), dir+string(os.PathSeparator))
	if err != nil {
		log.Println(err)
		return res, err
	}

	if len(artifacts) == 0 {
		err := errors.New("no artifacts found")
		log.Println(err)
		return res, err
	}

	a := artifacts[0]
	res = GetResponse{
		Version:  req.Version,
		Metadata: metadata(a),
	}

	return res, nil
}
