package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Arg struct {
	Key   string
	Value string
}

type Request struct {
	Module  string
	Args    []Arg
	Payload []byte
}

type Listener struct {
	mu       sync.Mutex
	Requests <-chan *Request
}

func TestCanPickUpJob(t *testing.T) {
	listener := &Listener{}

	listener.Requests <- &Request{
		Module: "module-a",
	}

	req := <-listener.Requests

	assert.Equal(t, "module-a", req.Module)
}
