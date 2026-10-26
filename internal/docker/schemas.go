package docker

type DockerContainer struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
