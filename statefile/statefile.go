package statefile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
)

const (
	stateFileDirEnv = "MUNIN_PLUGSTATE"
	stateFileName   = "containerlist"
)

var instance *stateFile

type stateFile struct {
	filename string
}

func Get() *stateFile {
	if instance == nil {
		instance = &stateFile{
			filename: fmt.Sprint(os.Getenv(stateFileDirEnv), "/", stateFileName),
		}
	}
	return instance
}

func (s *stateFile) exists() bool {
	_, err := os.Stat(s.filename)
	return err == nil
}

func (s *stateFile) SaveContainerList(c []types.Container) error {
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

func (s *stateFile) GetContainerList() ([]types.Container, error) {
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
