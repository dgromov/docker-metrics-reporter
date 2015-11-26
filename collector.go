package main

import (
	"os"
	"github.com/fsouza/go-dockerclient"
	"time"
	"fmt"
	"log"
	"sync"
)

// Map of container id to that container's done channel
var runningContainers = make(map[string]chan bool)
var runningLock = &sync.Mutex{}

type ContainerStat struct {
	ID string
	Config *docker.Config
	BaseStat *docker.Stats
	CPUPercent float32
	DiskUsage float32
}

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

func collect(id string, client *docker.Client, metricChannel chan *ContainerStat) {
	if doneChannel := knownContainer(id); doneChannel != nil {
		return
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
			ID: id,
			Stream: true,
			// This timeout is only for the intial connection
			Timeout: (time.Second * time.Duration(5)),
			Stats: stats,
			Done: doneChannel,
		})
	}()

	container, err := client.InspectContainer(id)
	if err != nil {
		log.Fatal(err)
	}

	i := 4
	// TODO: Stats get added to once a second. Abstract that so it can change.
	for stat := range stats {
		if i == 4 {
			metricChannel <- calculateStat(container, stat)
			i = 0
		}
		i += 1
	}
}

func calculateStat(cont *docker.Container, stat *docker.Stats) *ContainerStat {
	s := ContainerStat{
		BaseStat: stat,
		ID: cont.ID,
		Config: cont.Config,
		CPUPercent: 0,
		DiskUsage: 0,
	}

	return &s
}

func stopCollecting(id string) {
	fmt.Println("stopped")
	runningLock.Lock()
	defer runningLock.Unlock()
	if doneChannel, ok := runningContainers[id]; ok {
		doneChannel <- true
	}
}

func listenForContainers(client *docker.Client, metricChannel chan *ContainerStat) {
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
			go collect(event.ID, client, metricChannel)
		case "die":
			stopCollecting(event.ID)
		}
	}
}

func Collect(dockerEndpoint string, interval int, metricChannel chan *ContainerStat) {

	client, err := getClient(dockerEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	containers, _ := client.ListContainers(docker.ListContainersOptions{All: false})
	for _, container := range containers {
		fmt.Println(container.ID)
		go collect(container.ID, client, metricChannel)
	}

	listenForContainers(client, metricChannel)
}