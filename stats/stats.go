package stats

import (
	"strings"

	"github.com/docker/docker/api/types"
)

type ContainerStats types.StatsJSON

func (s *ContainerStats) CPUPercent() float64 {
	var (
		cpuPercent  = 0.0
		cpuDelta    = float64(s.CPUStats.CPUUsage.TotalUsage) - float64(s.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta = float64(s.CPUStats.SystemUsage) - float64(s.PreCPUStats.SystemUsage)
		onlineCPUs  = float64(s.CPUStats.OnlineCPUs)
	)
	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(s.CPUStats.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}
	return cpuPercent
}

func (s *ContainerStats) MemUsage() (int64, int64) {
	return int64(s.MemoryStats.Usage - s.MemoryStats.Stats["cache"]), int64(s.MemoryStats.Limit)
}

func (s *ContainerStats) BlockIO() (uint64, uint64) {
	var blkRead, blkWrite uint64
	for _, bioEntry := range s.BlkioStats.IoServiceBytesRecursive {
		switch strings.ToLower(bioEntry.Op) {
		case "read":
			blkRead = blkRead + uint64(bioEntry.Value)
		case "write":
			blkWrite = blkWrite + uint64(bioEntry.Value)
		}
	}
	return blkRead, blkWrite
}

func (s *ContainerStats) NetworkIO() (float64, float64) {
	var rx, tx float64
	for _, v := range s.Networks {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}
	return rx, tx
}
