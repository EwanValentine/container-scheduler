package cmd

import (
	"encoding/json"

	"github.com/EwanValentine/container-scheduler/api"
	"github.com/EwanValentine/container-scheduler/container"
	"github.com/EwanValentine/container-scheduler/invoker"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

var rawJSON = []byte(`{
	"level": "debug",
	"encoding": "json",
	"outputPaths": ["stdout", "/tmp/logs"],
	"errorOutputPaths": ["stderr"],
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

// Execute -
func Execute() {
	// Docker instance
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// New scheduler instance
	inv := invoker.NewInvoker(&invoker.Config{
		Logger: logging(),
		ContainerService: &container.Docker{
			Client: cli,
		},
	})

	httpapi := &api.HTTPAPI{
		Invoker: inv,
	}
	httpapi.Start()
}
