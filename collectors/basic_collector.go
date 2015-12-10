package collectors

import (
	"fmt"

	"github.com/fsouza/go-dockerclient"
	"log"
	"time"
	"github.com/dgromov/docker-metrics-reporter/common"
)

var BasicCollector Collector = Collector{
	collectFunc:          basicCollect,
	measureFunc:          basicMeasure,
	aggregateAndSendFunc: basicAggregate,
	shouldMeasureFunc:    basicShouldMeasure,
}

func basicCollect(id string, client *docker.Client, metricChannel chan *common.ContainerStat) error {
	if doneChannel := knownContainer(id); doneChannel != nil {
		return nil
	}
	cont, err := client.InspectContainer(id)
	if err != nil {
		log.Fatal(err)
	}

	if should_measure, _ := shouldMeasureFunc(cont); !should_measure {
		return nil
	}

	fmt.Println("starting", id)
	stats := make(chan *docker.Stats)
	doneChannel := make(chan bool)

	runningLock.Lock()
	runningContainers[id] = doneChannel
	runningLock.Unlock()
	defer fmt.Println("Ending collect for ", id)

	go func() {
		defer fmt.Println("Exited func for", id)
		client.Stats(docker.StatsOptions{
			ID:     id,
			Stream: true,
			// This timeout is only for the intial connection
			Timeout: (time.Second * time.Duration(5)),
			Stats:   stats,
			Done:    doneChannel,
		})
	}()



	err = basicAggregate(AggregationConfig{
		container:  cont,
		inChannel:  stats,
		outChannel: metricChannel,
		interval:   5,
	})

	if err != nil {
		fmt.Printf("Could not collect this round")
	}

	return nil
}

func flattenDockerStats(s *docker.Stats) (map[string]uint64, error) {
	return map[string]uint64{
		"memory.max":           s.MemoryStats.MaxUsage,
		"memory.usage":         s.MemoryStats.Usage,
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
		"memory.limit":         s.MemoryStats.Limit,

		"net.rx_bytes":   s.Network.RxBytes,
		"net.rx_dropped": s.Network.RxDropped,
		"net.rx_errors":  s.Network.RxErrors,
		"net.rx_packets": s.Network.RxPackets,
		"net.tx_bytes":   s.Network.TxBytes,
		"net.tx_dropped": s.Network.TxDropped,
		"net.tx_errors":  s.Network.TxErrors,
		"net.tx_packets": s.Network.TxPackets,

		"cpu.total":  s.CPUStats.CPUUsage.TotalUsage,

	}, nil
}

func calculateCPUUsage(stat *docker.Stats) (string, float64, error) {

	containerDelta := stat.CPUStats.CPUUsage.TotalUsage - stat.PreCPUStats.CPUUsage.TotalUsage
	totalDelta := stat.CPUStats.SystemCPUUsage - stat.PreCPUStats.SystemCPUUsage

	usage := (float64(containerDelta) / float64(totalDelta) * 100.0) * float64(len(stat.CPUStats.CPUUsage.PercpuUsage))

	return "cpu.usage", usage, nil
}

func basicMeasure(cont *docker.Container, stat *docker.Stats) *common.ContainerStat  {
	raw_metrics, _ := flattenDockerStats(stat)

	collectables := common.Collectables{
		Raw:        raw_metrics,
		Calculated: make(map[string]float64),
	}

	cpu_usage_percent, cpu_usage_value, _ := calculateCPUUsage(stat)
	collectables.AddCalculated(cpu_usage_percent, cpu_usage_value)

	s := common.ContainerStat{
		Collected: &collectables,
		ID:        cont.ID,
		Config:    cont.Config,
		Timestamp: stat.Read,
	}

	return &s
}

func basicAggregate(conf AggregationConfig) error {
	counter := conf.interval
	fmt.Println(counter)

	for stat := range conf.inChannel {
		if counter == conf.interval {
			conf.outChannel <- basicMeasure(conf.container, stat)
			counter = 0
		}
		counter += 1
	}
	return nil
}

func basicShouldMeasure(cont *docker.Container) (bool, error) {
	return true, nil
}