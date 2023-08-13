package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db *db
)

type db struct {
	Dsn string
	DbController *gorm.DB
}

func InitDB(dsn string) error {
	if dbObj, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return err
	} else {
		Db = &db{
			Dsn:        dsn,
			DbController: dbObj,
		}
		return nil
	}

}
