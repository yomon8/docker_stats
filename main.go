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

const (
	stateFileName = "containerlist"
)

var version string

func getContainerList(cli *client.Client, cl *statefile.ContainerList) ([]types.Container, error) {
	containerListFile, err := cl.Load()
	if err != nil {
		return nil, err
	}

	containerListApi, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	containerList := make([]types.Container, len(containerListFile))
	copy(containerList, containerListFile)
	for _, c := range containerListApi {
		isNewContainer := true
		for _, sc := range containerListFile {
			if sc.Names[0] == c.Names[0] {
				isNewContainer = false
				break
			}
		}
		if isNewContainer {
			containerList = append(containerList, c)
		}
	}

	err = cl.Save(containerList)
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

	cl, err := statefile.NewContainerList(stateFileName)
	if err != nil {
		log.Fatal("State File Error:", err)
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal("Create Docker Client Error:", err)
	}

	containers, err := getContainerList(cli, cl)
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
			cs := new(stats.ContainerStats)
			if err := json.Unmarshal(jsonBytes, cs); err != nil {
				log.Println("JSON Unmarshal error:", err)
			}
			p, err := statefile.NewContainerStatsFile(graph.GetKey(c.Names))
			if err != nil {
				log.Println("Stats state setting error:", err)
			}

			ps, err := p.Load()
			if err != nil {
				log.Println("Stats state read error:", err)
			}

			var (
				pbr, pbw uint64
				pnr, pnt float64
			)

			if ps != nil {
				pbr, pbw = ps.BlockIO()
				pnr, pnt = ps.NetworkIO()
			}
			cbr, cbw := cs.BlockIO()
			cnr, cnt := cs.NetworkIO()

			br, bw := cbr-pbr, cbw-pbw
			nr, nt := cnr-pnr, cnt-pnt
			mu, ml := cs.MemUsage()
			graph.AddMetricValues(&graph.MetricValues{
				ID:         c.ID[:12],
				Image:      c.Image,
				Names:      c.Names,
				CPUPercent: cs.CPUPercent(),
				MemUsage:   mu,
				MemLimit:   ml,
				BlockRD:    br,
				BlockWR:    bw,
				NetRCV:     nr,
				NetTRN:     nt,
			})

			p.Save(cs)
		}
	}
	graph.PrintMetricsValues()
	os.Exit(0)
}
