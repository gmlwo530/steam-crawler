package crawler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gmlwo530/steam-crawler/database"
	"gorm.io/gorm"
)

const steamSpyUrl = "https://steamspy.com/api.php?request=genre&genre=Indie"
const storeApiUrl = "https://store.steampowered.com"

func GetIndieAppList(db *gorm.DB) {
	if database.CountIndieApp(db) > 0 {
		return
	}

	resp, err := http.Get(steamSpyUrl)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	var indieAppsMap map[string]IndieAppRes
	err = json.NewDecoder(resp.Body).Decode(&indieAppsMap)

	if err != nil {
		log.Fatal(err)
	}

	indieApps := make([]database.IndieApp, 0, len(indieAppsMap))
	for _, val := range indieAppsMap {
		indieApps = append(indieApps, database.IndieApp{
			AppId:          uint(val.AppId),
			AverageForever: val.AverageForever,
			Ccu:            val.Ccu,
		})
	}

	database.CreateIndieApps(db, indieApps)
}

func UpdateIndieApp(db *gorm.DB, offset int, timeSleep time.Duration, debug bool) {
	notCrawledIndieApps := database.GetNotCrawledIndieApps(db)

	if debug {
		debugOffset := 100
		if len(notCrawledIndieApps) < debugOffset {
			debugOffset = len(notCrawledIndieApps)
		}
		notCrawledIndieApps = notCrawledIndieApps[:debugOffset]
	}

	for i := 0; i < len(notCrawledIndieApps); i += offset {
		indieApps := notCrawledIndieApps[i : i+offset]

		var crawledIndieApps []*database.IndieApp

		c := make(chan *database.IndieApp)
		errC := make(chan string)

		for _, indieApp := range indieApps {
			go getAppDetail(indieApp.AppId, "korean", &indieApp, c, errC)
		}

		for i := 0; i < len(indieApps); i++ {
			select {
			case indieApp := <-c:
				crawledIndieApps = append(crawledIndieApps, indieApp)
			case errStr := <-errC:
				log.Println(errStr)
			}
		}

		db.Save(crawledIndieApps)
		time.Sleep(timeSleep)
	}
}

func getAppDetail(appId uint, lang string, indieApp *database.IndieApp, c chan<- *database.IndieApp, errC chan<- string) {
	strAppId := strconv.Itoa(int(appId))

	resp, err := http.Get(storeApiUrl + "/api/appdetails?appids=" + strAppId + "&l=" + lang)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	var apr map[string]AppDetailRes
	err = json.NewDecoder(resp.Body).Decode(&apr)

	if err != nil {
		errC <- fmt.Sprintf("Error appId: %d, err: %+v", appId, err)
	} else {
		appDetail := apr[strAppId].Data

		for _, val := range appDetail.Movies {
			indieApp.Movies = append(indieApp.Movies, database.Movie{
				MovieId: uint(val.Id),
				Name:    val.Name,
				Mp4:     val.Mp4["480"],
			})
		}

		for _, val := range appDetail.Screenshots {
			indieApp.Screenshots = append(indieApp.Screenshots, database.Screenshot{
				ScreenshotId:  uint(val.Id),
				PathThumbnail: val.PathThumbnail,
				PathFull:      val.PathFull,
			})
		}

		indieApp.IndieAppDetails = append(indieApp.IndieAppDetails, database.IndieAppDetail{
			AppDetailId: uint(appDetail.AppId),
			Name:        appDetail.Name,
		})

		indieApp.AppReleaseDate = sql.NullString{String: appDetail.ReleaseDate.Date, Valid: true}
		indieApp.IsFree = sql.NullBool{Bool: appDetail.IsFree, Valid: true}

		c <- indieApp
	}
}
