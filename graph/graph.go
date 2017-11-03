package graph

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
)

const (
	category      = "docker"
	graphkeyCPU   = "doker_cpu"
	graphkeyMEM   = "doker_mem"
	graphkeyNWTX  = "doker_nw_tx"
	graphkeyNWRX  = "doker_nw_rx"
	graphkeyBLKWR = "doker_blk_wr"
	graphkeyBLKRD = "doker_blk_rd"
	cpukey        = "cpu"
	memkey        = "mem"
	nwkey         = "nw"
	blkkey        = "blk"
)

func getKey(names []string) string {
	var key string
	name := names[0]
	if strings.Contains(name, ".") {
		name = strings.Join(strings.Split(name, ".")[:2], "_")
	}

	key = strings.Replace(strings.Replace(name,
		"@", "", -1),
		"/", "", -1)

	return key
}

func PrintGraphDefinition(hostname string, containers []types.Container) {
	fmt.Printf("host_name %s\n\n", hostname)
	//CPU
	printGraphMetadata(graphkeyCPU, "CPU Usage", "1000", "%")
	printCPUGraph(containers)

	//Memory
	printGraphMetadata(graphkeyMEM, "MEM Usage", "1024", "Byte")
	printMemoryGraph(containers)

	//NW I/O
	printGraphMetadata(graphkeyNWTX, "Net I/O Transmit Exchange", "1024", "Byte")
	printNWGraphTX(containers)

	printGraphMetadata(graphkeyNWRX, "Net I/O Received Exchange", "1024", "Byte")
	printNWGraphRX(containers)

	//Block I/O
	printGraphMetadata(graphkeyBLKWR, "Block I/O Write", "1024", "Byte")
	printBLKGraphWR(containers)
	printGraphMetadata(graphkeyBLKRD, "Block I/O Read", "1024", "Byte")
	printBLKGraphRD(containers)
	fmt.Println("")
}

func printGraphMetadata(graphkey, title, base, vlabel string) {
	fmt.Println("")
	fmt.Printf("multigraph %s\n", graphkey)
	fmt.Printf("graph_title %s\n", title)
	fmt.Printf("graph_args --base %s\n", base)
	fmt.Printf("graph_vlabel %s\n", vlabel)
	fmt.Printf("graph_category %s\n", category)
	fmt.Println("")
}

func printCPUGraph(containers []types.Container) {
	for _, c := range containers {
		key := getKey(c.Names)
		fmt.Printf("%s_%s.label %s\n", key, cpukey, key)
		fmt.Printf("%s_%s.type %s\n", key, cpukey, "GAUGE")
		fmt.Printf("%s_%s.min %d\n", key, cpukey, 0)
		fmt.Printf("%s_%s.draw %s\n", key, cpukey, "AREASTACK")
	}
}

func printMemoryGraph(containers []types.Container) {
	fmt.Printf("limit_%s.label limit\n", memkey)
	fmt.Printf("limit_%s.type %s\n", memkey, "GAUGE")
	fmt.Printf("limit_%s.min %d\n", memkey, 0)
	fmt.Printf("limit_%s.draw %s\n", memkey, "AREA")
	for _, c := range containers {
		key := getKey(c.Names)
		fmt.Printf("%s_%s.label %s\n", key, memkey, key)
		fmt.Printf("%s_%s.type %s\n", key, memkey, "GAUGE")
		fmt.Printf("%s_%s.min %d\n", key, memkey, 0)
		fmt.Printf("%s_%s.draw %s\n", key, memkey, "AREASTACK")
	}
}

func printNWGraphTX(containers []types.Container) {
	for _, c := range containers {
		key := getKey(c.Names)
		fmt.Printf("%s_%stx.label %s_tx\n", key, nwkey, key)
		fmt.Printf("%s_%stx.type %s\n", key, nwkey, "GAUGE")
		fmt.Printf("%s_%stx.min %d\n", key, nwkey, 0)
		fmt.Printf("%s_%stx.draw %s\n", key, nwkey, "AREASTACK")
	}
}
func printNWGraphRX(containers []types.Container) {
	for _, c := range containers {
		key := getKey(c.Names)
		fmt.Printf("%s_%srx.label %s_rx\n", key, nwkey, key)
		fmt.Printf("%s_%srx.type %s\n", key, nwkey, "GAUGE")
		fmt.Printf("%s_%srx.min %d\n", key, nwkey, 0)
		fmt.Printf("%s_%srx.draw %s\n", key, nwkey, "AREASTACK")
	}
}

func printBLKGraphWR(containers []types.Container) {
	for _, c := range containers {
		key := getKey(c.Names)
		fmt.Printf("%s_%swr.label %s_wr\n", key, blkkey, key)
		fmt.Printf("%s_%swr.type %s\n", key, blkkey, "GAUGE")
		fmt.Printf("%s_%swr.min %d\n", key, blkkey, 0)
		fmt.Printf("%s_%swr.draw %s\n", key, blkkey, "AREASTACK")
	}
}
func printBLKGraphRD(containers []types.Container) {
	for _, c := range containers {
		key := getKey(c.Names)
		fmt.Printf("%s_%srd.label %s_rd\n", key, blkkey, key)
		fmt.Printf("%s_%srd.type %s\n", key, blkkey, "GAUGE")
		fmt.Printf("%s_%srd.min %d\n", key, blkkey, 0)
		fmt.Printf("%s_%srd.draw %s\n", key, blkkey, "AREASTACK")
	}
}

type MetricValues struct {
	ID         string
	Image      string
	CPUPercent float64
	Names      []string
	MemUsage   int64
	MemLimit   int64
	BlockRD    uint64
	BlockWR    uint64
	NetRCV     float64
	NetTRN     float64
}

var values = make(map[string]*MetricValues, 10)

func AddMetricValues(m *MetricValues) {
	key := getKey(m.Names)
	values[key] = m
}

func PrintMetricsValues() {
	var keys []string
	for k := range values {
		keys = append(keys, k)
	}

	//cpu
	fmt.Println("")
	fmt.Printf("multigraph %s\n", graphkeyCPU)
	for _, k := range keys {
		fmt.Printf("%s_%s.value %f\n", k, cpukey, values[k].CPUPercent)
	}
	//memory
	fmt.Println("")
	fmt.Printf("multigraph %s\n", graphkeyMEM)
	for i, k := range keys {
		if i == 0 {
			fmt.Printf("limit_%s.value %d\n", memkey, values[k].MemLimit)
		}
		fmt.Printf("%s_%s.value %d\n", k, memkey, values[k].MemUsage)
	}
	//network
	fmt.Println("")
	fmt.Printf("multigraph %s\n", graphkeyNWTX)
	for _, k := range keys {
		fmt.Printf("%s_%stx.value %f\n", k, nwkey, values[k].NetTRN)
	}
	fmt.Println("")
	fmt.Printf("multigraph %s\n", graphkeyNWRX)
	for _, k := range keys {
		fmt.Printf("%s_%srx.value %f\n", k, nwkey, values[k].NetRCV)
	}
	//block i/o
	fmt.Println("")
	fmt.Printf("multigraph %s\n", graphkeyBLKWR)
	for _, k := range keys {
		fmt.Printf("%s_%swr.value %d\n", k, blkkey, values[k].BlockWR)
	}
	fmt.Println("")
	fmt.Printf("multigraph %s\n", graphkeyBLKRD)
	for _, k := range keys {
		fmt.Printf("%s_%srd.value %d\n", k, blkkey, values[k].BlockRD)
	}
	fmt.Println("")
}
