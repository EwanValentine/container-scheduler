package container

// Container -
type Container struct {
	ID       string `json:"id"`
	Host     string `json:"host"`
	ImageID  string `json:"image_id"`
	Endpoint string `json:"endpoint"`

	// When called, status will be set to
	// true, after 30 seconds, will be set to false
	Status bool `json:"status"`
}
