package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Docker is a wrapper for the docker client
type Docker struct {
	client *client.Client
}

// Run runs the docker run command with an image ID
// returns container ID
func (d *Docker) Run(id string) (string, error) {
	cnf := &container.Config{
		Image: id,
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
	}

	hostCnf := &container.HostConfig{
		AutoRemove: true,
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "8080",
				},
			},
		},
	}

	cnt, err := d.client.ContainerCreate(context.Background(), cnf, hostCnf, nil, "")
	if err != nil {
		return "", err
	}
	return cnt.ID, d.client.ContainerStart(context.Background(), cnt.ID, types.ContainerStartOptions{})
}

// Terminate kills a running container
func (d *Docker) Terminate(id string) error {
	return d.client.ContainerKill(context.Background(), id, "SIGKILL")
}
