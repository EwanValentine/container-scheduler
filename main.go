package main

import (
	docker "docker.io/go-docker"
)

func main() {

	// Listen for requests

	// Run docker container for incoming module request
	// Kill it after 30 seconds

	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

}
