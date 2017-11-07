package statefile

import (
	"errors"
	"io/ioutil"
	"os"
)

const (
	stateFileDirEnv = "MUNIN_PLUGSTATE"
)

var (
	ErrNoFile = errors.New("state file not exist")
)

type statefile struct {
	filename string
}

func (f *statefile) exists() bool {
	_, err := os.Stat(f.filename)
	return err == nil
}

func (f *statefile) save(b []byte) error {
	fl, err := os.Create(f.filename)
	if err != nil {
		return err
	}

	_, err = fl.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (f *statefile) load() ([]byte, error) {
	if !f.exists() {
		return nil, ErrNoFile
	}
	fl, err := os.Open(f.filename)
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(fl)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
