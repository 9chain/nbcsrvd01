package state

import (
	"github.com/jinzhu/gorm"
	_  "github.com/jinzhu/gorm/dialects/sqlite"
	"fmt"
	"time"
	"bytes"
	"github.com/BurntSushi/toml"
	"path"
	"github.com/9chain/nbcsrvd01/config"
	"os"
	"io/ioutil"
	"path/filepath"
	"sort"
)

var (
	DB *gorm.DB
)

func Init() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database " + err.Error())
	}

	db.LogMode(true)
	DB = db
}

func BackupUserConfig() {
	fmt.Println("backup user config")
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		fmt.Println(err)
		return
	}

	var chains []UserChain
	if err := DB.Find(&chains).Error; err != nil {
		fmt.Println(err)
		return
	}

	type backupCfg struct {
		users []User
		chains []UserChain
	}
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(backupCfg{users:users, chains:chains}); err != nil {
		fmt.Println( err)
		return
	}

	filename := time.Now().Format("20060102-150405.000")
	cfgPath := path.Join(config.Cfg.User.ConfigDir, "backup.toml")
	bakfile := path.Join(config.Cfg.User.ConfigDir, "backup.toml."+filename)

	if _, err := os.Stat(cfgPath); err == nil {
		if err := os.Rename(cfgPath, bakfile); err != nil {
			fmt.Println(err)
			return
		}
	}

	if err := ioutil.WriteFile(cfgPath, buf.Bytes(), 0644); err != nil {
		fmt.Println(err)
		return
	}

	clearBackupFiles()
}


type stringSlice []string

func (s stringSlice) Len() int           { return len(s) }
func (s stringSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s stringSlice) Less(i, j int) bool { return s[i] < s[j] }

func clearBackupFiles() {
	files, err := filepath.Glob(path.Join(config.Cfg.User.ConfigDir, "backup.toml.*"))
	if err != nil {
		fmt.Println("glob fail", err)
		return
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
