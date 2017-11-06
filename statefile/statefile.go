package statefile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/docker/docker/api/types"
)

const (
	stateFileName   = "containerlist"
	stateFileDirEnv = "MUNIN_PLUGSTATE"
)

type StateFile struct {
	filename string
}

func NewStateFile() (*StateFile, error) {
	dir := os.Getenv(stateFileDirEnv)
	if dir == "" {
		return nil, fmt.Errorf("os env [%s] is not set:", stateFileDirEnv)
	}

	return &StateFile{
		filename: path.Join(dir, stateFileName),
	}, nil
}

func (s *StateFile) exists() bool {
	_, err := os.Stat(s.filename)
	return err == nil
}

func (s *StateFile) SaveContainerList(c []types.Container) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f, err := os.Create(s.filename)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (s *StateFile) GetContainerList() ([]types.Container, error) {
	var savedContainerList []types.Container
	if !s.exists() {
		return savedContainerList, nil
	}
	f, err := os.Open(s.filename)
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bs, &savedContainerList)
	if err != nil {
		return nil, err
	}
	return savedContainerList, nil
}
