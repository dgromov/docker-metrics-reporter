package collectors

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/dgromov/docker-metrics-reporter/common"
)

var DiskUsageCollector Collector = Collector{
	collectFunc:          BasicCollector.collectFunc,
	aggregateAndSendFunc: BasicCollector.aggregateAndSendFunc,
	measureFunc:          CPUAndDiskMeasure,
	shouldMeasureFunc:    BasicCollector.shouldMeasureFunc,
}

func calculateDiskUsage(cont *docker.Container) (string, float64, error) {
	return "disk.usage", 0.0, nil
}

func CPUAndDiskMeasure(cont *docker.Container, stat *docker.Stats) *common.ContainerStat {
	stats := BasicCollector.measureFunc(cont, stat)

	disk_usage_name, disk_usage_value, _ := calculateDiskUsage(cont)
	stats.Collected.AddCalculated(disk_usage_name, disk_usage_value)

	return stats
}
