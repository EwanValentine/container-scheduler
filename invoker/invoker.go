package invoker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/EwanValentine/container-scheduler/api"
	"go.uber.org/zap"
)

// ContainerService -
type containerService interface {
	Run(id string) (string, error)
	Terminate(id string) error
}

type container struct {
	ID      string
	ImageID string
	Endoint string

	// When called, status will be set to
	// true, after 30 seconds, will be set to false
	Status bool
}

// Config config for the invoker
type Config struct {
	Logger *zap.Logger
	ContainerService
}

// NewInvoker returns a new Invoker instance
func NewInvoker(config *Config) *Invoker {
	return &Invoker{
		logger:           config.logger,
		containerService: config.cs,
		containers:       make(map[string]*container, 0),
	}
}

// Invoker houses the main functionality
// used to call containers
type Invoker struct {
	mu sync.Mutex
	containerService
	containers map[string]*container
}

// Timeout removes containers after a set time-period
func (s *Invoker) Timeout(name string) {
	time.Sleep(30 * time.Second)
	s.mu.Lock()
	s.containers[name].Status = false
	id := s.containers[name].ID
	s.mu.Unlock()
	s.containerService.Terminate(id)
}

// Invoke a module
func (s *Invoker) Invoke(request *api.Request) ([]byte, error) {
	s.mu.Lock()
	container, ok := s.containers[request.Module]
	s.mu.Unlock()

	if !ok {
		return nil, errors.New("No container found")
	}

	if container.Status == false {

		// Run container, and wait for it to start
		containerID, err := s.containerService.Run(container.ImageID)
		if err != nil {
			return nil, err
		}

		s.mu.Lock()
		s.containers[request.Module].ID = containerID
		s.containers[request.Module].Status = true
		s.mu.Unlock()

		// Wait for container to start
		s.poll(container)

		// Timeout in 30 seconds, sets container status to false again
		go s.Timeout(request.Module)
	}

	return s.call(container, request)
}

// Register a new container
func (s *Invoker) Register(name string, container *container) error {
	s.mu.Lock()
	s.containers[name] = container
	s.mu.Unlock()
	return nil
}

// call makes
func (s *Invoker) call(container *container.Container, request *api.Request) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:8080/%s", "localhost", request.Endpoint))
	if err != nil {
		return nil, err
	}
	s.Logger.Info(request.Endpoint)
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// poll polls the health-check endpoint recursively until a valid response is returned
func (s *Invoker) poll(container *container) bool {
	retries := 0
	_, err := http.Get("http://localhost:8080/_health")
	if err != nil {
		time.Sleep(1 * time.Second)
		retries++
		if retries > 10 {
			return false
		}
		s.poll(container)
	}
	return true
}
