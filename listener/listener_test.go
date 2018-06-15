package listener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanPickUpJob(t *testing.T) {
	listener := &Listener{}

	listener.Requests <- &Request{
		Module: "module-a",
	}

	req := <-listener.Requests

	assert.Equal(t, "module-a", req.Module)
}
