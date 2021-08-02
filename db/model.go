package db

type IndieApp struct {
	AppId          int `json:"appid"`
	AverageForever int `json:"average_forever"` // average playtime since March 2009. In minutes.
	Ccu            int `json:"ccu"`             // peak CCU yesterday.
}

type AppDetailResp struct {
	Success bool
	Data    AppDetail
}

type AppGenre struct {
	Id          string
	Description string
	AppId       int
}

type AppScreenshot struct {
	Id            int
	PathThumbnail string
	PathFull      string
	AppId         int
}

type AppMovie struct {
	Id    int
	Name  string
	Mp4   map[string]string
	AppId int
}

type AppReleaseDate struct {
	Date       string `json:"date"`
	ComingSoon bool   `json:"coming_soon"`
}

type AppDetail struct {
	Name                string
	AppId               int `json:"steam_appid"`
	IsFree              bool
	DetailedDescription string         `json:"detailed_description"`
	ReleaseDate         AppReleaseDate `json:"release_date"`
	Genres              []AppGenre
	Screenshots         []AppScreenshot
	Movies              []AppMovie
}
