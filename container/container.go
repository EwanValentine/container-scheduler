package container

// Container -
type Container struct {
	ID       string
	Host     string
	ImageID  string
	Endpoint string

	// When called, status will be set to
	// true, after 30 seconds, will be set to false
	Status bool
}
