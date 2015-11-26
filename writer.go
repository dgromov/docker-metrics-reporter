package main

import (
	"fmt"
)

func Write(dest string, port int, prefix string, metricChannel chan *ContainerStat) {
	fmt.Println("IAM HERE")
	for metric := range metricChannel {
		fmt.Printf("Recieved metric for %s at %v\n", metric.ID, metric.BaseStat.Read)

//		go writeToGraphite()
	}
}
