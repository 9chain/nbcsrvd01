package filestate

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type FileState struct {
	users map[string]User
}

var (
	once      sync.Once
	fileState *FileState
)

func (s *FileState) CheckLogin(username, password string) error {
	fmt.Println(username, password)
	user, ok := s.users[username]
	fmt.Println(ok, user.Username, user.Password)
	if user, ok := s.users[username]; ok && user.Password == password {
		return nil
	}
	return errors.New("invalid username/password")
}

func (s *FileState) GetUserApiKey(username string) (string, bool) {
	if user, ok := s.users[username]; ok {
		return user.ApiKey, ok
	}

	return "", false
}

func (s *FileState) ResetUserApiKey(username string) (string, error) {
	user, ok := s.users[username]
	if !ok {
		return "", errors.New("no such user")
	}

	user.ApiKey = genRandomApiKey()
	user.Modified = time.Now()

	s.users[username] = user

	if err := saveUserConfig(s.users); err != nil {
		return "", err
	}

	return user.ApiKey, nil
}

func (s *FileState) GetUserState(username string) (int, error) {
	user, ok := s.users[username]
	if ok {
		return user.State, nil
	}

	return -1, errors.New("not exists")
}

func (s *FileState) UpdateUserState(username string, state int) error {
	user, ok := s.users[username]
	if !ok {
		return errors.New("no such user")
	}

	user.State = state
	user.Modified = time.Now()

	s.users[username] = user

	if err := saveUserConfig(s.users); err != nil {
		return err
	}

	return nil
}

func (s *FileState) UpdateUserPassword(username, password string) error {
	user, ok := s.users[username]
	if !ok {
		return errors.New("no such user")
	}

	user.Password = password
	user.Modified = time.Now()

	s.users[username] = user

	if err := saveUserConfig(s.users); err != nil {
		return err
	}

	return nil
}

func (s *FileState) AddNewUser(username, password, email string) error {
	if _, ok := s.users[username]; ok {
		return errors.New("why user exists!")
	}

	user := User{
		Email:    email,
		Username: username,
		Password: password,
		State:    0,
		ApiKey:   genRandomApiKey(),
		Created:  time.Now(),
		Modified: time.Now(),
		Chains:   []UserChain{},
	}

	s.users[username] = user
	if err := saveUserConfig(s.users); err != nil {
		return err
	}

	return nil
}

func (s *FileState) UpdateModified(username string) error {
	user, ok := s.users[username]
	if !ok {
		return errors.New("not found user")
	}
	user.Modified = time.Now()
	s.users[username] = user
	return nil
}

func (s *FileState) ShouldSendEmail(username, email string) bool {
	user, ok := s.users[username]
	if !ok {
		return false
	}

	if user.Email != email {
		return true
	}

	now := time.Now()
	least := user.Modified.Add(time.Minute * 2)
	if now.After(least) {
		return true
	}

	return false
}

func (s *FileState) EmailInfo(username string) (string, time.Time, error) {
	user, ok := s.users[username]
	if !ok {
		return "", time.Now(), errors.New("no such user")
	}

	return user.Email, user.Modified, nil
}

func (s *FileState) UpdateEmailTime(username string) error {
	user, ok := s.users[username]
	if !ok {
		return errors.New("no such user")
	}
	user.Modified = time.Now()
	s.users[username] = user
	if err := saveUserConfig(s.users); err != nil {
		return err
	}
	return nil
}

func (s *FileState) UpdateUserInfo(username, password, email string) error {
	user, ok := s.users[username]
	if !ok {
		return errors.New("not found user")
	}
	user.Modified = time.Now()
	user.Email = email
	user.Password = password
	s.users[username] = user
	if err := saveUserConfig(s.users); err != nil {
		return err
	}
	return nil
}

func (s *FileState) ResetPassword(username, password string) error {
	user, ok := s.users[username]
	if !ok {
		return errors.New("not found user")
	}

	if password == user.Password {
		return nil
	}

	user.Modified = time.Now()
	user.Password = password
	s.users[username] = user
	if err := saveUserConfig(s.users); err != nil {
		return err
	}

	return nil
}

func Init() (*FileState, error) {
	once.Do(func() {
		var err error
		users, err := loadUserConfig()
		if err != nil {
			panic(err)
		}

		fileState = &FileState{
			users: users,
		}
	})

	return fileState, nil
}
