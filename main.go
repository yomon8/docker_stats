package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/yomon8/docker_stats/graph"
	"github.com/yomon8/docker_stats/stats"
)

var version string

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal("Create Docker Client Error:", err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatal("Get Container Error:", err)
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "config":
			hostname, err := os.Hostname()
			if err != nil {
				log.Fatal("Get Hostname Error:", err)
			}
			graph.PrintGraphDefinition(hostname, containers)
			os.Exit(0)
		case "-v", "--version", "version":
			fmt.Println("version: ", version)
			os.Exit(0)
		}
	}

	for _, c := range containers {
		cstat, err := cli.ContainerStats(context.Background(), c.ID, false)
		if err != nil {
			continue
		}
		sc := bufio.NewScanner(cstat.Body)
		if sc.Scan() {
			jsonBytes := []byte(sc.Text())
			sf := new(stats.ContainerStats)
			if err := json.Unmarshal(jsonBytes, sf); err != nil {
				log.Println("JSON Unmarshal error:", err)
			}
			mu, ml := sf.MemUsage()
			br, bw := sf.BlockIO()
			nr, nt := sf.NetworkIO()
			graph.AddMetricValues(&graph.MetricValues{
				ID:         c.ID[:12],
				Image:      c.Image,
				Names:      c.Names,
				CPUPercent: sf.CPUPercent(),
				MemUsage:   mu,
				MemLimit:   ml,
				BlockRD:    br,
				BlockWR:    bw,
				NetRCV:     nr,
				NetTRN:     nt,
			})
		}
	}
	graph.PrintMetricsValues()
	os.Exit(0)
}
