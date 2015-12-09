package writers

import (
	"fmt"
	"github.com/dgromov/docker-to-graphite/common"
	"time"
)

type ConsoleWriter struct{}

func (c ConsoleWriter) SendInt(label string, value uint64) {
	fmt.Printf("%s -> %d\n", label, value)
}

func (c ConsoleWriter) SendFloat(label string, value float64) {
	fmt.Printf("%s -> %f\n", label, value)
}

func (c ConsoleWriter) Label(container *common.ContainerStat) string {
	return container.ID
}

func (c ConsoleWriter) Timestamp(container *common.ContainerStat) time.Time {
	return container.Timestamp
}
