package collectors

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fsouza/go-dockerclient"
	"github.com/dgromov/docker-metrics-reporter/common"
)

// Map of container id to that container's done channel
var runningContainers = make(map[string]chan bool)
var runningLock = &sync.Mutex{}

func getClient(endpoint string) (*docker.Client, error) {
	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost != "" {
		return docker.NewClientFromEnv()
	}

	return docker.NewClient(endpoint)
}

func knownContainer(id string) chan bool {
	runningLock.Lock()
	defer runningLock.Unlock()

	if doneChannel, ok := runningContainers[id]; ok {
		return doneChannel
	}
	return nil
}

type AggregationConfig struct {
	container  *docker.Container
	inChannel  chan *docker.Stats
	outChannel chan *common.ContainerStat
	interval   int
}

type Collector struct {
	collectFunc          func(string, *docker.Client, chan *common.ContainerStat) error
	measureFunc          func(*docker.Container, *docker.Stats) *common.ContainerStat
	aggregateAndSendFunc func(AggregationConfig) error
	shouldMeasureFunc  func(*docker.Container) (bool, error)
}

func (c *Collector) listenForContainers(client *docker.Client, metricChannel chan *common.ContainerStat) {
	eventStream := make(chan *docker.APIEvents)
	err := client.AddEventListener(eventStream)
	if err != nil {
		log.Fatal(err)
	}

	defer client.RemoveEventListener(eventStream)

	for event := range eventStream {
		fmt.Println(" ", event.ID, event.Status)
		switch event.Status {
		case "start":
			go c.collectFunc(event.ID, client, metricChannel)
		case "die":
			c.stopCollecting(event.ID)
		}
	}
}

func (c *Collector) stopCollecting(id string) {
	fmt.Println("stopped")
	runningLock.Lock()
	defer runningLock.Unlock()
	if doneChannel, ok := runningContainers[id]; ok {
		doneChannel <- true
	}
}

func (c *Collector) Collect(dockerEndpoint string, interval int, metricChannel chan *common.ContainerStat) {

	client, err := getClient(dockerEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	containers, _ := client.ListContainers(docker.ListContainersOptions{All: false})
	for _, container := range containers {
		fmt.Println(container.ID)
		go c.collectFunc(container.ID, client, metricChannel)
	}

	c.listenForContainers(client, metricChannel)
}
