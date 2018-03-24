package filestate

import (
	"bytes"
	"fmt"
	"github.com/9chain/nbcsrvd01/config"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type UserChain struct {
	ChainName string
}

type User struct {
	Username string
	Password string
	Email    string
	ApiKey   string
	State    int
	Chains   []UserChain
	Created  time.Time
	Modified time.Time
}

var rwlock = sync.Mutex{}

func loadUserConfig() (map[string]User, error) {
	cfgPath := path.Join(config.Cfg.User.ConfigDir, "user.toml")
	users := make(map[string]User)

	if _, err := os.Stat(cfgPath); err != nil {
		fmt.Println("path not exists", cfgPath, err.Error())
		return users, nil
	}

	if _, err := toml.DecodeFile(cfgPath, &users); err != nil {
		fmt.Println("toml.DecodeFile fail", cfgPath, err.Error())
		return nil, err
	}

	return users, nil
}

func saveUserConfig(users map[string]User) error {
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(users); err != nil {
		panic(err)
	}

	rwlock.Lock()
	defer rwlock.Unlock()

	filename := time.Now().Format("20060102-150405.000")
	cfgPath := path.Join(config.Cfg.User.ConfigDir, "user.toml")
	bakfile := path.Join(config.Cfg.User.ConfigDir, "user.toml."+filename)
	if err := os.Rename(cfgPath, bakfile); err != nil {
		return err
	}

	if err := ioutil.WriteFile(cfgPath, buf.Bytes(), 0644); err != nil {
		return err
	}

	go clearBackupFiles()
	return nil
}

type stringSlice []string

func (s stringSlice) Len() int           { return len(s) }
func (s stringSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s stringSlice) Less(i, j int) bool { return s[i] < s[j] }

func clearBackupFiles() {
	files, err := filepath.Glob(path.Join(config.Cfg.User.ConfigDir, "user.toml.*"))
	if err != nil {
		fmt.Println("glob fail", err)
	}

	deleteN := len(files) - config.Cfg.User.MaxUserFiles
	if deleteN <= 0 {
		return
	}

	sort.Sort(stringSlice(files))
	toDelete := files[0:deleteN]
	for _, filename := range toDelete {
		os.Remove(filename)
	}
}
