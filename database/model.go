package database

import (
	"database/sql"
)

type IndieApp struct {
	AppId           uint `gorm:"primarykey;column:indie_app_id"`
	AverageForever  int
	Ccu             int
	AppReleaseDate  sql.NullString
	IsFree          sql.NullBool
	Movies          []Movie          `gorm:"foreignKey:IndieAppId;references:AppId"`
	Screenshots     []Screenshot     `gorm:"foreignKey:IndieAppId;references:AppId"`
	IndieAppDetails []IndieAppDetail `gorm:"foreignKey:IndieAppId;references:AppId"`
}

type Movie struct {
	MovieId    uint `gorm:"primarykey;column:movie_id"`
	Name       string
	Mp4        string
	IndieAppId uint
}

type Screenshot struct {
	ScreenshotId  uint `gorm:"primarykey;column:screenshot_id"`
	PathThumbnail string
	PathFull      string
	IndieAppId    uint
}

type IndieAppDetail struct {
	AppDetailId uint `gorm:"primarykey;column:app_detail_id"`
	Name        string
	IndieAppId  uint
}
