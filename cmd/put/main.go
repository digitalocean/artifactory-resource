package main

import (
	"log"
	"os"

	resource "github.com/digitalocean/artifactory-resource"
	rlog "github.com/digitalocean/concourse-resource-library/log"
	jlog "github.com/jfrog/jfrog-client-go/utils/log"
)

func main() {
	input := rlog.WriteStdin()
	defer rlog.Close()

	jlog.SetLogger(jlog.NewLogger(jlog.DEBUG, log.Writer()))

	var request resource.PutRequest
	err := request.Read(input)
	if err != nil {
		log.Fatalf("failed to read request input: %s", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("missing arguments")
	}
	dir := os.Args[1]

	response, err := resource.Put(request, dir)
	if err != nil {
		log.Fatalf("failed to perform put: %s", err)
	}

	err = response.Write()
	if err != nil {
		log.Fatalf("failed to write response to stdout: %s", err)
	}

	log.Println("Put complete")
}
