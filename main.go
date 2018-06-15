package main

import (
	"encoding/json"
	"log"

	"github.com/EwanValentine/container-scheduler/invoker"
	"github.com/docker/docker/client"
	"github.com/moby/moby/container"
	"go.uber.org/zap"
)

var rawJSON = []byte(`{
	"level": "debug",
	"encoding": "json",
	"outputPaths": ["stdout", "/tmp/logs"],
	"errorOutputPaths": ["stderr"],
	"initialFields": {"foo": "bar"},
	"encoderConfig": {
		"messageKey": "message",
		"levelKey": "level",
		"levelEncoder": "lowercase"
	}
}`)

func logging() *zap.Logger {
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	return logger
}

func main() {
	wait := make(chan bool)

	// Docker instance
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	docker := &Docker{
		client: cli,
	}

	// New scheduler instance
	s := New(&invoker.Config{
		Logger:           logger(),
		ContainerService: docker,
	})

	// @todo - make endpoint
	// CI process will build the image
	// fire the endpoint, image id etc at this
	// service.
	s.Register("module-a", &container.Container{
		ImageID: "efa78366f8b0",
		Endoint: "/",
		Status:  false,
	})

	// This would be on the end of a listener or something
	response, err := s.Invoke(&Request{
		Module: "module-a",
	})
	if err != nil {
		panic(err)
	}

	response, err = s.Invoke(&Request{
		Module: "module-a",
	})
	if err != nil {
		panic(err)
	}

	log.Println(string(response))

	<-wait
}
