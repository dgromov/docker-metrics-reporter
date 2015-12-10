package main

import (
	"flag"
	"log"
	"os"

	"github.com/dgromov/docker-metrics-reporter/common"
	"github.com/dgromov/docker-metrics-reporter/writers"
	"github.com/dgromov/docker-metrics-reporter/collectors"
)

func main() {
	docker := flag.String("endpoint", "", "Docker endpoint")
	interval := flag.Int("interval", 60, "interval to report")

	//	if *dest == "" {
	//		flag.PrintDefaults()
	//		os.Exit(1)
	//	}

	if *docker == "" {
		host := os.Getenv("DOCKER_HOST")
		if host == "" {
			message :=
				`Your environment does not have docker configured.
			Please supply an endpoint`

			log.Fatal(message)
		}
	}

	flag.Parse()
	metricChannel := make(chan *common.ContainerStat)
	go writers.Write(writers.ConsoleWriter{}, metricChannel)
	collectors.DiskUsageCollector.Collect(*docker, *interval, metricChannel)

	// TODO: Calculate CPU usage percent
	// TODO: Calculate Disk usage
	// TODO: Aggregate metrics that can fluctuate, i.e. memory usage
	// TODO: Come up with abstraction for paths
	// TODO: Dry Run mode
	// TODO: Actually Send
	// TODO: runit script (or other)
	// TODO: FPM -> DEB -> S3
	// TODO: UNIT TESTS! BUT HAO!?
}
