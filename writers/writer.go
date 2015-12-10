package writers

import (
	"github.com/dgromov/docker-metrics-reporter/common"
)

type Writer interface {
	Label(*common.ContainerStat, string) string
	SendInt(string, uint64)
	SendFloat(string, float64)
}

// This assumes that the writer is already connected to whatever its writing to.
func Write(w Writer, metricChannel chan *common.ContainerStat) {
	for metric := range metricChannel {
		for name, value := range metric.Collected.Raw {
			w.SendInt(w.Label(metric, name), value)
		}

		for name, value := range metric.Collected.Calculated {
			w.SendFloat(w.Label(metric, name), value)
		}
	}
}
