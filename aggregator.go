package main

import (
	"github.com/fsouza/go-dockerclient"
	"fmt"
)

type Aggregator interface {
	// Should take from the input
	Aggregate(AggregatorConfig) error
}


type AggregatorConfig struct {
	container *docker.Container
	inChannel chan *docker.Stats
	outChannel chan *ContainerStat
}


type BasicAggregator struct{
	interval int
}

func calculateStat(cont *docker.Container, stat *docker.Stats) *ContainerStat {
	s := ContainerStat{
		BaseStat: stat,
		ID: cont.ID,
		Config: cont.Config,
		CPUPercent: 0,
		DiskUsage: 0,
	}

	return &s
}

func (agg BasicAggregator) Aggregate(conf AggregatorConfig) error {
	counter := agg.interval
	fmt.Println(counter)

	for stat := range conf.inChannel {
		if counter == agg.interval {
			conf.outChannel <- calculateStat(conf.container, stat)
			counter = 0
		}
		counter += 1
	}
	return nil
}

