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

//	metrics := map[string]uint64{
//		"cpu.user":   s.Stats.CPUStats.CPUUsage.UsageInUsermode,
//		"cpu.system": s.Stats.CPUStats.CPUUsage.UsageInKernelmode,
//		"cpu.total":  s.Stats.CPUStats.CPUUsage.TotalUsage,
//
//		"memory.limit": s.Stats.MemoryStats.Limit,
//		"memory.max":   s.Stats.MemoryStats.MaxUsage,
//		"memory.usage": s.Stats.MemoryStats.Usage,
//
//		"memory.active_anon":   s.Stats.MemoryStats.Stats.TotalActiveAnon,
//		"memory.active_file":   s.Stats.MemoryStats.Stats.TotalActiveFile,
//		"memory.cache":         s.Stats.MemoryStats.Stats.TotalCache,
//		"memory.inactive_anon": s.Stats.MemoryStats.Stats.TotalInactiveAnon,
//		"memory.inactive_file": s.Stats.MemoryStats.Stats.TotalInactiveFile,
//		"memory.mapped_file":   s.Stats.MemoryStats.Stats.TotalMappedFile,
//		"memory.pg_fault":      s.Stats.MemoryStats.Stats.TotalPgfault,
//		"memory.pg_in":         s.Stats.MemoryStats.Stats.TotalPgpgin,
//		"memory.pg_out":        s.Stats.MemoryStats.Stats.TotalPgpgout,
//		"memory.rss":           s.Stats.MemoryStats.Stats.TotalRss,
//		"memory.rss_huge":      s.Stats.MemoryStats.Stats.TotalRssHuge,
//		"memory.unevictable":   s.Stats.MemoryStats.Stats.TotalUnevictable,
//		"memory.writeback":     s.Stats.MemoryStats.Stats.TotalWriteback,
//
//		"net.rx_bytes":   s.Stats.Network.RxBytes,
//		"net.rx_dropped": s.Stats.Network.RxDropped,
//		"net.rx_errors":  s.Stats.Network.RxErrors,
//		"net.rx_packets": s.Stats.Network.RxPackets,
//		"net.tx_bytes":   s.Stats.Network.TxBytes,
//		"net.tx_dropped": s.Stats.Network.TxDropped,
//		"net.tx_errors":  s.Stats.Network.TxErrors,
//		"net.tx_packets": s.Stats.Network.TxPackets,
//	}

