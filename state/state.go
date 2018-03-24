package state

import (
	"github.com/9chain/nbcsrvd01/state/filestate"
	"time"
)

var (
	State GlobalState
)

type GlobalState interface {
	CheckLogin(username, password string) error
	GetUserApiKey(username string) (string, bool)
	ResetUserApiKey(username string) (string, error)
	GetUserState(username string) (int, error)
	AddNewUser(username, password, email string) error
	UpdateModified(username string) error
	ShouldSendEmail(username, email string) bool
	UpdateUserInfo(username, password, email string) error
	ResetPassword(username, password string) error
	EmailInfo(username string) (string, time.Time, error)
	UpdateEmailTime(username string) error
	UpdateUserState(username string, state int) error
	UpdateUserPassword(username, password string) error
}

func Init() {
	s, err := filestate.Init()
	if err != nil {
		panic(err)
	}

	State = s
}
