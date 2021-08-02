package crawler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gmlwo530/steam-crawler/db"
)

const steamSpyUrl = "https://steamspy.com/api.php?request=genre&genre=Indie"
const apiUrl = "https://api.steampowered.com"
const storeApiUrl = "https://store.steampowered.com"

func GetIndieAppList(dbObj *sql.DB) {
	resp, err := http.Get(steamSpyUrl)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	var indieAppsMap map[string]db.IndieApp
	err = json.NewDecoder(resp.Body).Decode(&indieAppsMap)

	if err != nil {
		log.Fatal(err)
	}

	indieApps := make([]db.IndieApp, 0, len(indieAppsMap))
	for _, val := range indieAppsMap {
		indieApps = append(indieApps, val)
	}

	db.InsertIndieApp(dbObj, indieApps)

	// i := 0
	// for _, val := range indieApps {
	// 	if i == 10 {
	// 		break
	// 	}

	// 	log.Printf("IndieApp : %+v", val)
	// 	i++
	// }
}

func CreateAppDetail(dbObj *sql.DB) {
	indieApps := db.SelectIndieApp(dbObj, 100, 0)

	var appDetails []db.AppDetail

	c := make(chan db.AppDetail)

	for _, indieApp := range indieApps {
		go GetAppDetail(indieApp.AppId, "english", c)
	}

	for i := 0; i < len(indieApps); i++ {
		appDetail := <-c
		appDetails = append(appDetails, appDetail)
	}

	db.InsertAppDetail(dbObj, appDetails)
}

func GetAppDetail(appId int, lang string, c chan<- db.AppDetail) {
	strAppId := strconv.Itoa(appId)

	resp, err := http.Get(storeApiUrl + "/api/appdetails?appids=" + strAppId + "&l=" + lang)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	var apr map[string]db.AppDetailResp
	err = json.NewDecoder(resp.Body).Decode(&apr)

	if err != nil {
		log.Fatal(err)
	}

	c <- apr[strAppId].Data
}
