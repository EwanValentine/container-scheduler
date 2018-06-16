package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
	Endpoint string `json:"endpoint"`
	Module   string `json:"module"`
	Args     []Arg  `json:"args"`
	Payload  []byte `json:"payload"`
}

type invoker interface {
	Invoke(*Request) ([]byte, http.Header, error)
	Register(string, *container.Container) error
}

// HTTPAPI is the main api server
type HTTPAPI struct {
	Invoker invoker
}

func (httpapi *HTTPAPI) health(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder.Encode(map[string]string{"status": "OK"})
}

func (httpapi *HTTPAPI) get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	endpoint := params.ByName("name")
	module := strings.TrimPrefix(endpoint, "/")

	request := &Request{
		Endpoint: "/" + params.ByName("name"),
		Module:   module,
	}
	payload, headers, err := httpapi.Invoker.Invoke(request)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		encoder := json.NewEncoder(w)
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	contentType := headers.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func (httpapi *HTTPAPI) post(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	var request *Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	payload, headers, err := httpapi.Invoker.Invoke(request)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		encoder := json.NewEncoder(w)
		encoder.Encode(map[string]string{"error": err.Error()})
		return
	}

	contentType := headers.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func (httpapi *HTTPAPI) register(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	var cnt *container.Container
	if err := json.NewDecoder(r.Body).Decode(&cnt); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	moduleName := strings.TrimPrefix(cnt.Endpoint, "/")
	if err := httpapi.Invoker.Register(moduleName, cnt); err != nil {
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

	// Healthcheck
	router.GET("/_health", httpapi.health)

	// Invokers
	router.GET("/modules/:name", httpapi.get)
	router.POST("/modules/:name", httpapi.post)

	// Register module
	router.POST("/modules", httpapi.register)

	// Start server
	log.Fatal(http.ListenAndServe(":3000", router))
}
