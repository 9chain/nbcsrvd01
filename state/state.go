package state

import (
	"github.com/jinzhu/gorm"
	_  "github.com/jinzhu/gorm/dialects/sqlite"
	"fmt"
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
}