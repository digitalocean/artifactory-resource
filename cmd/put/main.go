package main

import (
	"log"

	resource "github.com/digitalocean/artifactory-resource"
	rlog "github.com/digitalocean/concourse-resource-library/log"
)

func main() {
	input := rlog.WriteStdin()
	defer rlog.Close()

	var request resource.PutRequest
	err := request.Read(input)
	if err != nil {
		log.Fatalf("failed to read request input: %s", err)
	}

	response, err := resource.Put(request)
	if err != nil {
		log.Fatalf("failed to perform put: %s", err)
	}

	err = response.Write()
	if err != nil {
		log.Fatalf("failed to write response to stdout: %s", err)
	}

	log.Println("Put complete")
}
