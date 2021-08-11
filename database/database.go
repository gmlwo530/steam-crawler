package database

import (
	"fmt"
	"log"
	"os"

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

	dbName := os.Getenv("DATABASE_NAME")

	if dt == MYSQL {
		user := os.Getenv("MYSQL_USER")
		password := os.Getenv("MYSQL_PASSWORD")
		dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=UTC", user, password, dbName)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else if dt == SQLITE3 {
		db, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", dbName)), &gorm.Config{})
	} else {
		log.Fatal("Wrong database type")
	}

	if err != nil {
		log.Fatal(err)
	}

	Migrate(db, dt)

	return db
}
