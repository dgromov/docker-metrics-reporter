package writers

import (
	"fmt"
	"github.com/dgromov/docker-metrics-reporter/common"
)

type ConsoleWriter struct{}

func (c ConsoleWriter) SendInt(label string, value uint64) {
	fmt.Printf("%s -> %d\n", label, value)
}

func (c ConsoleWriter) SendFloat(label string, value float64) {
	fmt.Printf("%s -> %f\n", label, value)
}

func (c ConsoleWriter) Label(container *common.ContainerStat, name string) string {
	return fmt.Sprintf("%s - %s.%s", container.Timestamp, container.ID, name)
}
