package main

import (
	"log"
	"os"
	"time"

	"github.com/gmlwo530/steam-crawler/crawler"
	"github.com/gmlwo530/steam-crawler/database"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db := database.GetDB(database.SQLITE3)

	crawler.GetIndieAppList(db)
	crawler.UpdateIndieApp(db, 10, time.Second*5, true)

	log.Println("Crawling is Done!")
	os.Exit(100)
}
