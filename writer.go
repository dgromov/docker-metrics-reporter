package main

import (
	"fmt"
	"time"
)

type Writer interface {
	Label(*ContainerStat) string
	Send(string, uint64)
}

type ConsoleWriter struct{}

func (c ConsoleWriter) Send(label string, value uint64) {
	fmt.Printf("%s -> %d\n", label, value)
}

func (c ConsoleWriter) Label(container *ContainerStat) string {
	return container.ID
}

func (c ConsoleWriter) Timestamp(container *ContainerStat) time.Time {
	return container.BaseStat.Read
}


// This assumes that the writer is already connected to whatever its writing to.
func Write(w Writer, metricChannel chan *ContainerStat) {
	for metric := range metricChannel {
		name := w.Label(metric) + "." + "cpu.used"
		fmt.Println(metric.Config.Volumes)
		w.Send(name, metric.BaseStat.CPUStats.CPUUsage.TotalUsage)
	}
}



