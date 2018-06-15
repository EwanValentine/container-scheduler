package api

// Arg is a key value pair,
// query params for example
type Arg struct {
	Key   string
	Value string
}

// Request is used throughout this codebase
// be wary of changes
type Request struct {
	Endpoint string
	Module   string
	Args     []Arg
	Payload  []byte
}

// Server is the main api server
type Server struct{}

func (server *Server) Start() {}
