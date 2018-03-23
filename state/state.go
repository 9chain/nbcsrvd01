package state

import (
	"github.com/9chain/nbcsrvd01/state/filestate")

var (
	State GlobalState
)

type GlobalState interface {
	CheckLogin(username, password string) error
}

func Init() {
	s, err := filestate.Init()
	if err != nil {
		panic(err)
	}

	State = s
}