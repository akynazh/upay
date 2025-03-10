package model

import (
	"github.com/akynazh/upay/app/config"
	"github.com/akynazh/upay/app/help"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var _err error

func Init() error {
	var dbPath = config.GetDbPath()
	if !help.IsExist(dbPath) {
		DB, _err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if _err != nil {

			return _err
		}

		DB.Exec(installSql)
		addStartWalletAddress()

		return nil
	}

	DB, _err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if _err == nil {

		addStartWalletAddress()
	}

	return _err
}
