package invoker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/EwanValentine/container-scheduler/api"
	"github.com/EwanValentine/container-scheduler/container"
	"github.com/stretchr/testify/assert"
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

type mockContainerService struct{}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test")
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func (c *mockContainerService) Run(id string) (string, error) {
	http.HandleFunc("/test-module", handler)
	http.HandleFunc("/_health", health)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	return "some_id", nil
}

func (c *mockContainerService) Terminate(id string) error {
	return nil
}

func setup() *Invoker {
	config := &Config{
		Logger:           logging(),
		ContainerService: &mockContainerService{},
	}
	return NewInvoker(config)
}

func TestCanInvokeContainer(t *testing.T) {
	i := setup()
	i.Register("test-module", &container.Container{
		ImageID:  "test",
		Host:     "0.0.0.0",
		Endpoint: "/test-module",
	})
	payload, _, err := i.Invoke(&api.Request{
		Endpoint: "/test-module",
		Module:   "test-module",
	})
	assert.NoError(t, err)
	assert.Equal(t, "test", string(payload))
}
