package container

import (
	"context"

	"github.com/docker/docker/api/types"
	dcontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Docker is a wrapper for the docker client
type Docker struct {
	Client *client.Client
}

// Run runs the docker run command with an image ID
// returns container ID
func (d *Docker) Run(id string) (string, error) {
	cnf := &dcontainer.Config{
		Image: id,
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
	}

	hostCnf := &dcontainer.HostConfig{
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

	cnt, err := d.Client.ContainerCreate(context.Background(), cnf, hostCnf, nil, "")
	if err != nil {
		return "", err
	}
	return cnt.ID, d.Client.ContainerStart(context.Background(), cnt.ID, types.ContainerStartOptions{})
}

// Terminate kills a running container
func (d *Docker) Terminate(id string) error {
	return d.Client.ContainerKill(context.Background(), id, "SIGKILL")
}
