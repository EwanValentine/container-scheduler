package invoker

import "testing"

type containerService struct{}

func (c *containerService) Run(id string) (string, error) {
	return "some_id", nil
}

func (c *containerService) Terminate(id string) error {
	return nil
}

func TestCanInvokeContainer(t *testing.T) {
	config := &Config{
		ContainerService: &containerService{},
	}
	i := NewInvoker(config)
}
