package main

import (
	"github.com/fsouza/go-dockerclient"
	"fmt"
)

type Aggregator interface {
	// Should take from the input
	Aggregate(AggregatorConfig) error
	GetRawStats(*docker.Stats) map[string]uint64
	GetCalculatedStats(*docker.Stats) map[string]float64

}

type AggregatorConfig struct {
	container *docker.Container
	inChannel chan *docker.Stats
	outChannel chan *ContainerStat
}

type BasicAggregator struct{
	interval int
}

func (b BasicAggregator) GetCalculatedStats(s *docker.Stats) map[string]float64 {
	return map[string]float64{}
}

func (b BasicAggregator) GetRawStats(s *docker.Stats) map[string]uint64 {
	return map[string]uint64{
		// Variable
		"memory.max":   s.MemoryStats.MaxUsage,
		"memory.usage": s.MemoryStats.Usage,
		"memory.active_anon":   s.MemoryStats.Stats.TotalActiveAnon,
		"memory.active_file":   s.MemoryStats.Stats.TotalActiveFile,
		"memory.cache":         s.MemoryStats.Stats.TotalCache,
		"memory.inactive_anon": s.MemoryStats.Stats.TotalInactiveAnon,
		"memory.inactive_file": s.MemoryStats.Stats.TotalInactiveFile,
		"memory.mapped_file":   s.MemoryStats.Stats.TotalMappedFile,
		"memory.pg_fault":      s.MemoryStats.Stats.TotalPgfault,
		"memory.pg_in":         s.MemoryStats.Stats.TotalPgpgin,
		"memory.pg_out":        s.MemoryStats.Stats.TotalPgpgout,
		"memory.rss":           s.MemoryStats.Stats.TotalRss,
		"memory.rss_huge":      s.MemoryStats.Stats.TotalRssHuge,
		"memory.unevictable":   s.MemoryStats.Stats.TotalUnevictable,
		"memory.writeback":     s.MemoryStats.Stats.TotalWriteback,

		// Only Grow
		"net.rx_bytes":   s.Network.RxBytes,
		"net.rx_dropped": s.Network.RxDropped,
		"net.rx_errors":  s.Network.RxErrors,
		"net.rx_packets": s.Network.RxPackets,
		"net.tx_bytes":   s.Network.TxBytes,
		"net.tx_dropped": s.Network.TxDropped,
		"net.tx_errors":  s.Network.TxErrors,
		"net.tx_packets": s.Network.TxPackets,
		"cpu.user":   s.CPUStats.CPUUsage.UsageInUsermode,
		"cpu.system": s.CPUStats.CPUUsage.UsageInKernelmode,
		"cpu.total":  s.CPUStats.CPUUsage.TotalUsage,
		"memory.limit": s.MemoryStats.Limit,
	}
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

