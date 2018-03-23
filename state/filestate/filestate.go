package filestate

import (
	"errors"
)

type User struct {
	Username string
	Password string
}

type FileState struct {
	users map[string]User
}

func (s *FileState) CheckLogin(username, password string) error {
	if pwd, ok := s.users[username]; ok && pwd.Password == password {
		return nil
	}
	return errors.New("invalid username/password")
}


func Init() (*FileState, error) {
	return &FileState{}, nil
}