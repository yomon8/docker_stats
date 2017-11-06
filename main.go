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
	"github.com/yomon8/docker_stats/statefile"
	"github.com/yomon8/docker_stats/stats"
)

var version string

func getContainerList(cli *client.Client, statefile *statefile.StateFile) ([]types.Container, error) {
	containerList, err := statefile.GetContainerList()
	if err != nil {
		return nil, err
	}

	containerListApi, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	for _, c := range containerListApi {
		isNewContainer := true
		for _, sc := range containerList {
			if sc.Names[0] == c.Names[0] {
				isNewContainer = false
				break
			}
		}
		if isNewContainer {
			containerList = append(containerList, c)
		}
	}

	err = statefile.SaveContainerList(containerList)
	if err != nil {
		return nil, err
	}

	return containerList, nil
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version", "version":
			fmt.Println("version: ", version)
			os.Exit(0)
		}
	}

	statefile, err := statefile.NewStateFile()
	if err != nil {
		log.Fatal("State File Error:", err)
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal("Create Docker Client Error:", err)
	}

	containers, err := getContainerList(cli, statefile)
	if err != nil {
		log.Fatal("Get ContainerList Error:", err)
	}

	if len(os.Args) > 1 {
		if "config" == os.Args[1] {
			hostname, err := os.Hostname()
			if err != nil {
				log.Fatal("Get Hostname Error:", err)
			}
			graph.PrintGraphDefinition(hostname, containers)
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
