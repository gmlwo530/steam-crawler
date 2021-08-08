package database

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type dbType int

const (
	SQLITE3 dbType = iota
	MYSQL
)

func GetDB(dt dbType) *gorm.DB {
	var db *gorm.DB
	var err error

	if dt == MYSQL {
		dsn := "root:root@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=UTC"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else if dt == SQLITE3 {
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	} else {
		log.Fatal("Wrong database type")
	}

	if err != nil {
		log.Fatal(err)
	}

	Migrate(db, dt)

	return db
}
