package main

import "testing"

type ContainerService interface {
	Run(name string) error
}

type Docker struct{}

// Run runs the docker run command
func (d *Docker) Run(name string) error {

}

type Container struct {
	Endoint string

	// When called, status will be set to
	// true, after 30 seconds, will be set to false
	Status bool
}

type Scheduler struct {
	mu sync.Mutex
	ContainerService
	Containers map[string]Container
}

func (s *Scheduler) Timeout(name string) {
	time.Sleep(30 * time.Second)
	s.mu.Lock()
	s.Containers[name].Status = false
	s.mu.Unlock()
	s.Terminate(name)
}

func (s *Scheduler) Poll(name string) bool {
	// Poll container until ready
	// Poll localhost:8080/_health until 200

	// Wait until poll status is ready
	// @todo - should also include a time-out
	for {
		select {
			ready := <-s.poll
				if ready {
					return true
				}
		}
	}
}

func (s *Scheduler) Invoke(request *Request) error {
	if s.ready == false {
		return errors.New("Endpoint is not ready")
	}

	// Actually call endpoint here
	s.mu.Lock()
	container := s.Containers[request.Module]
	s.mu.Unlock()

	if container.Status == false {
		s.ContainerService.Run(request.Module)
		ready := s.Poll()
		if ready == false {
			return errors.New("Timed out")
		}

		s.mu.Lock()
		s.Containers[request.Module].Status = true
		s.mu.Unlock()

		// Timeout in 30 seconds, sets container status 
		// to false again
		go container.Timeout()
		return nil
	}

	return nil
}

func (s *Scheduler) Terminate(name string) error {
	// Run command to terminate container
}

func TestCanInvokeContainer(t *testing.T) {

}
