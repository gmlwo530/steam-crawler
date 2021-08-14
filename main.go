package main

import (
	"log"
	"os"
	"time"

	"github.com/gmlwo530/steam-crawler/config"
	"github.com/gmlwo530/steam-crawler/crawler"
	"github.com/gmlwo530/steam-crawler/database"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config.InitConfig()

	db := database.GetDB(database.MYSQL)

	crawler.GetIndieAppList(db)
	crawler.UpdateIndieApp(db, time.Second*2)

	log.Println("Crawling is Done!")
	os.Exit(100)
}
