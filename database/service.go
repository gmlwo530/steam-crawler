package database

import (
	"log"

	"gorm.io/gorm"
)

func CountIndieApp(db *gorm.DB) int64 {
	var count int64

	result := db.Model(&IndieApp{}).Distinct("indie_app_id").Count(&count)

	log.Printf("Count: %d", count)

	checkError(result)

	return count
}

func getCrawledIndieAppIds(db *gorm.DB) []int64 {
	var indieAppIds []int64

	result := db.Table("indie_app_details").Select("indie_app_id").Find(&indieAppIds)

	checkError(result)

	return indieAppIds
}

func GetNotCrawledIndieApps(db *gorm.DB) []IndieApp {
	var indieApps []IndieApp

	result := db.Not(getCrawledIndieAppIds(db)).Find(&indieApps)

	checkError(result)

	return indieApps
}

func CreateIndieApps(db *gorm.DB, indieApps []IndieApp) {
	batchSize := 1000
	loopCount := len(indieApps) / batchSize

	for i := 0; i < loopCount+1; i++ {
		rangeEnd := i*batchSize + batchSize
		if rangeEnd < len(indieApps) {
			rangeEnd = len(indieApps)
		}

		result := db.CreateInBatches(indieApps[i*batchSize:rangeEnd], batchSize)

		checkError(result)
	}
}

func checkError(result *gorm.DB) {
	if result.Error != nil {
		log.Fatal(result.Error)
	}
}
