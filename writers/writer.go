package writers

import (
	"github.com/dgromov/docker-to-graphite/common"
)

type Writer interface {
	Label(*common.ContainerStat) string
	SendInt(string, uint64)
	SendFloat(string, float64)
}

// This assumes that the writer is already connected to whatever its writing to.
func Write(w Writer, metricChannel chan *common.ContainerStat) {
	for metric := range metricChannel {
		for name, value := range metric.Collected.Raw {
			w.SendInt(name, value)
		}

		for name, value := range metric.Collected.Calculated {
			w.SendFloat(name, value)
		}
	}
}
