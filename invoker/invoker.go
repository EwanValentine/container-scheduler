package invoker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/EwanValentine/container-scheduler/api"
	"github.com/EwanValentine/container-scheduler/container"
	"go.uber.org/zap"
)

// ContainerService -
type containerService interface {
	Run(id string) (string, error)
	Terminate(id string) error
}

// Config config for the invoker
type Config struct {
	Logger           *zap.Logger
	ContainerService containerService
}

// NewInvoker returns a new Invoker instance
func NewInvoker(config *Config) *Invoker {
	return &Invoker{
		logger:           config.Logger,
		containerService: config.ContainerService,
		containers:       make(map[string]*container.Container, 0),
	}
}

// Invoker houses the main functionality
// used to call containers
type Invoker struct {
	mu sync.Mutex
	containerService
	containers map[string]*container.Container
	logger     *zap.Logger
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
func (s *Invoker) Invoke(request *api.Request) ([]byte, http.Header, error) {
	s.mu.Lock()
	container, ok := s.containers[request.Module]
	s.mu.Unlock()

	if !ok {
		return nil, nil, errors.New("No container found")
	}

	if container.Status == false {

		// Run container, and wait for it to start
		containerID, err := s.containerService.Run(container.ImageID)
		if err != nil {
			return nil, nil, err
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
func (s *Invoker) Register(name string, container *container.Container) error {
	s.mu.Lock()
	s.containers[name] = container
	s.mu.Unlock()
	return nil
}

// call makes
func (s *Invoker) call(container *container.Container, request *api.Request) ([]byte, http.Header, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:8080%s", container.Host, container.Endpoint))
	if err != nil {
		return nil, nil, err
	}
	s.logger.Info(request.Endpoint)
	defer resp.Body.Close()
	payload, err := ioutil.ReadAll(resp.Body)
	return payload, resp.Header, err
}

// poll polls the health-check endpoint recursively until a valid response is returned
func (s *Invoker) poll(container *container.Container) bool {
	retries := 0
	_, err := http.Get("http://localhost:8080/_health")
	if err != nil {
		log.Println("Retrying:", retries)
		time.Sleep(1 * time.Second)
		retries++
		if retries > 10 {
			return false
		}
		s.poll(container)
	}
	return true
}
