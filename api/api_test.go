package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/EwanValentine/container-scheduler/container"
	"github.com/stretchr/testify/assert"
)

type fakeInvoker struct{}

func (fake *fakeInvoker) Invoke(request *Request) ([]byte, error) {
	return []byte("test"), nil
}

func (fake *fakeInvoker) Register(module string, container *container.Container) error {
	return nil
}

type testReq struct {
	expected string
	method   string
	body     []byte
}

func TestEndpoints(t *testing.T) {
	golden := map[string]testReq{
		"/_health":             testReq{expected: "{\"status\":\"OK\"}\n", method: "GET"},
		"/modules/test-module": testReq{expected: "test", method: "GET"},
		"/modules": testReq{
			expected: "{\"success\":true}\n",
			method:   "POST",
			body:     []byte(`{"image_id": "test123", "endpoint": "/test-module", "name": "test-module"}`),
		},
	}

	httpapi := &HTTPAPI{
		invoker: &fakeInvoker{},
	}

	go httpapi.Start()
	for endpoint, req := range golden {
		log.Println(endpoint, req.expected)
		url := fmt.Sprintf("http://localhost:8080/%s", endpoint)
		var err error
		var res *http.Response
		switch req.method {
		case "POST":
			request, err := http.NewRequest("POST", url, bytes.NewBuffer(req.body))
			assert.NoError(t, err)
			client := &http.Client{}
			res, err = client.Do(request)

		case "GET":
			res, err = http.Get(url)
		}

		defer res.Body.Close()
		actual, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, req.expected, string(actual))
	}
}
