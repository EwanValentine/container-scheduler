package listener

import (
	"sync"

	"github.com/EwanValentine/container-scheduler/api"
)

type Listener struct {
	mu       sync.Mutex
	Requests <-chan *api.Request
}
