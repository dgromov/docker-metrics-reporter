package main

import (
	"os"
	"github.com/fsouza/go-dockerclient"
	"time"
	"fmt"
)

const DOCKER_ENDPOINT = "unix:///var/run/docker.sock"

var running_containers = make(map[string]chan *docker.Stats)

func getClient() (*docker.Client, error) {
	endpoint, _ := docker.DefaultDockerHost()
	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost != "" {
		return docker.NewClientFromEnv()
	}

	return docker.NewClient(endpoint)
}

func fetchStats(id string, client *docker.Client) (chan *docker.Stats, chan bool) {
	statChannel := make(chan *docker.Stats)
	doneChannel := make(chan bool)
	go client.Stats(docker.StatsOptions{
		ID: id,
		Stream: true,
		Timeout: (time.Second*time.Duration(5)),
		Stats: statChannel,
		Done: doneChannel,
	})

	return statChannel, doneChannel
}

func Collect() {
	client, err := getClient()
	if err != nil {
		fmt.Println(err)
	}

	for {
		fmt.Println("Trying... ")
		containers, _ := client.ListContainers(docker.ListContainersOptions{All: false})
		for _, container := range containers {
			statChan := running_containers[container.ID]
			if running_containers[container.ID] == nil {
				fmt.Printf("Fetching new")
				statChan, _ = fetchStats(container.ID, client)
				running_containers[container.ID] = statChan
			}
			stat := <- statChan

			fmt.Printf("ID: %s, READ: %v\n", container.ID, stat.Read)
		}
		time.Sleep(time.Second*time.Duration(3))
	}
}