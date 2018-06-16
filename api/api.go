package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/EwanValentine/container-scheduler/container"
	"github.com/julienschmidt/httprouter"
)

// Arg is a key value pair,
// query params for example
type Arg struct {
	Key   string
	Value string
}

// Request is used throughout this codebase
// be wary of changes
type Request struct {
	Endpoint string
	Module   string
	Args     []Arg
	Payload  []byte
}

type invoker interface {
	Invoke(*Request) ([]byte, error)
	Register(string, *container.Container) error
}

// HTTPAPI is the main api server
type HTTPAPI struct {
	invoker
}

func (httpapi *HTTPAPI) health(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder.Encode(map[string]string{"status": "OK"})
}

func (httpapi *HTTPAPI) call(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	payload, err := httpapi.invoker.Invoke(&Request{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		encoder := json.NewEncoder(w)
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Write(payload)
}

func (httpapi *HTTPAPI) register(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	moduleName := params.ByName("name")
	if err := httpapi.invoker.Register(moduleName, &container.Container{
		Endpoint: moduleName,
		Host:     "todo",
		ImageID:  "todo",
		Status:   false,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	encoder.Encode(map[string]bool{
		"success": true,
	})
}

// Start server
func (httpapi *HTTPAPI) Start() {
	router := httprouter.New()
	router.GET("/_health", httpapi.health)
	router.GET("/modules/:name", httpapi.call)
	router.POST("/modules", httpapi.register)
	log.Fatal(http.ListenAndServe(":8080", router))
}
