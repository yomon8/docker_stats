package statefile

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/docker/docker/api/types"
)

type ContainerList struct {
	f statefile
}

func NewContainerList(filename string) (*ContainerList, error) {
	dir := os.Getenv(stateFileDirEnv)
	if dir == "" {
		return nil, fmt.Errorf("os env [%s] is not set:", stateFileDirEnv)
	}

	return &ContainerList{
		f: statefile{
			filename: path.Join(dir, filename),
		},
	}, nil
}

func (cl *ContainerList) Save(c []types.Container) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err := cl.f.save(b); err != nil {
		return err
	}

	return nil
}

func (cl *ContainerList) Load() ([]types.Container, error) {
	var savedContainerList []types.Container
	bs, err := cl.f.load()
	if err == ErrNoFile {
		return savedContainerList, nil
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bs, &savedContainerList)
	if err != nil {
		return nil, err
	}
	return savedContainerList, nil
}
