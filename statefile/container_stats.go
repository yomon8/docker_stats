package statefile

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/yomon8/docker_stats/stats"
)

type ContainerStatsFile struct {
	f statefile
}

func NewContainerStatsFile(filename string) (*ContainerStatsFile, error) {
	dir := os.Getenv(stateFileDirEnv)
	if dir == "" {
		return nil, fmt.Errorf("os env [%s] is not set:", stateFileDirEnv)
	}

	return &ContainerStatsFile{
		f: statefile{
			filename: path.Join(dir, filename),
		},
	}, nil
}

func (cs *ContainerStatsFile) Save(c *stats.ContainerStats) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err := cs.f.save(b); err != nil {
		return err
	}

	return nil
}

func (cs *ContainerStatsFile) Load() (*stats.ContainerStats, error) {
	var savedContainerStats *stats.ContainerStats
	bs, err := cs.f.load()
	if err == ErrNoFile {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bs, &savedContainerStats)
	if err != nil {
		return nil, err
	}
	return savedContainerStats, nil
}
