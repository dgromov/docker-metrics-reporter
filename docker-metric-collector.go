package main
import (
	"flag"
	"os"
	"log"
)

func main() {
	docker := flag.String("endpoint", "", "Docker endpoint")
	dest := flag.String("graphite-host", "", "Graphite host")
	port := flag.Int("graphite-port", 8888, "Graphite port")
	prefix := flag.String("metric-prefix", "", "graphite metric prefix")
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
	metricChannel := make(chan *ContainerStat)


	go Write(*dest, *port, *prefix, metricChannel)
	Collect(*docker, *interval, metricChannel)
}

