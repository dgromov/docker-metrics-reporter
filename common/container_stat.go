package common

import (
	"github.com/fsouza/go-dockerclient"
	"time"
)

type ContainerStat struct {
	ID        string
	Config    *docker.Config
	Collected *Collectables
	Timestamp time.Time
}
