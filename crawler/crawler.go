package crawler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gmlwo530/steam-crawler/config"
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

func UpdateIndieApp(db *gorm.DB, timeSleep time.Duration, offset int) {
	c := make(chan database.IndieApp)
	errC := make(chan string)

	notCrawledIndieApps := database.GetNotCrawledIndieApps(db)

	if config.GetConfig().Debug {
		debugOffset := 100
		if len(notCrawledIndieApps) < debugOffset {
			debugOffset = len(notCrawledIndieApps)
		}
		notCrawledIndieApps = notCrawledIndieApps[:debugOffset]
	}

	languages := []string{"korean", "english"}

	loopCount := len(notCrawledIndieApps) / offset

	for i := 0; i <= loopCount; i++ {
		endSlice := i*offset + offset

		if len(notCrawledIndieApps) < endSlice {
			endSlice = len(notCrawledIndieApps)
		}

		notCrawledIndieAppsPartial := notCrawledIndieApps[i*offset : endSlice]

		for _, indieApp := range notCrawledIndieAppsPartial {
			for _, lang := range languages {
				go getAppDetail(indieApp, lang, c, errC)
			}
		}

		for i := 0; i < len(notCrawledIndieAppsPartial)*len(languages); i++ {
			select {
			case indieApp := <-c:
				db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&indieApp)
			case errStr := <-errC:
				log.Println(errStr)
			}
		}
		time.Sleep(timeSleep)
	}
}

func getAppDetail(indieApp database.IndieApp, lang string, c chan<- database.IndieApp, errC chan<- string) {
	strAppId := strconv.Itoa(int(indieApp.AppId))

	resp, err := http.Get(storeApiUrl + "/api/appdetails?appids=" + strAppId + "&l=" + lang)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	var apr map[string]AppDetailRes
	err = json.NewDecoder(resp.Body).Decode(&apr)

	if err != nil {
		errC <- fmt.Sprintf("Error appId: %s, err: %+v", strAppId, err)
	} else {
		appDetail := apr[strAppId].Data
		for _, val := range appDetail.Movies {
			indieApp.Movies = append(indieApp.Movies, database.Movie{
				SteamMovieId: uint(val.Id),
				Name:         val.Name,
				Mp4:          val.Mp4["480"],
			})
		}

		for _, val := range appDetail.Screenshots {
			indieApp.Screenshots = append(indieApp.Screenshots, database.Screenshot{
				SteamScreenshotId: uint(val.Id),
				PathThumbnail:     val.PathThumbnail,
				PathFull:          val.PathFull,
			})
		}

		for _, val := range appDetail.Genres {
			id, err := strconv.Atoi(val.Id)
			if err != nil {
				log.Fatalf("Wrong genre ID: %s", val.Id)
				continue
			}
			indieApp.Genres = append(indieApp.Genres, database.Genre{
				SteamGenreId: uint(id),
				Description:  val.Description,
				Language:     lang,
			})
		}

		indieApp.IndieAppDetails = append(indieApp.IndieAppDetails, database.IndieAppDetail{
			Name:             appDetail.Name,
			ReleaseDate:      appDetail.ReleaseDate.Date,
			ShortDescription: appDetail.ShortDescription,
			Language:         lang,
		})

		indieApp.HeaderImage = sql.NullString{String: appDetail.HeaderImage, Valid: true}
		indieApp.IsFree = sql.NullBool{Bool: appDetail.IsFree, Valid: true}

		c <- indieApp
	}
}
