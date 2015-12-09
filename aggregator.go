package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
)

type Aggregator interface {
	// Should take from the input
	Aggregate(AggregatorConfig) error
	GetRawStats(*docker.Stats) map[string]uint64
	GetCalculatedStats(*docker.Stats) map[string]float64
}

type AggregatorConfig struct {
	container  *docker.Container
	inChannel  chan *docker.Stats
	outChannel chan *ContainerStat
}

type BasicAggregator struct {
	interval int
}

func (b BasicAggregator) GetCalculatedStats(s *docker.Stats) map[string]float64 {
	return map[string]float64{}
}

func (b BasicAggregator) GetRawStats(s *docker.Stats) map[string]uint64 {

}

func calculateStat(cont *docker.Container, stat *docker.Stats) *ContainerStat {
	s := ContainerStat{
		BaseStat: stat,
		ID:       cont.ID,
		Config:   cont.Config,
	}

	return &s
}
