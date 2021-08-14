package database

import (
	"database/sql"
)

type IndieApp struct {
	AppId           uint `gorm:"primarykey;column:indie_app_id"`
	Name            string
	AverageForever  int
	Ccu             int
	IsFree          sql.NullBool
	HeaderImage     sql.NullString
	Movies          []Movie          `gorm:"foreignKey:IndieAppId;references:AppId"`
	Screenshots     []Screenshot     `gorm:"foreignKey:IndieAppId;references:AppId"`
	IndieAppDetails []IndieAppDetail `gorm:"foreignKey:IndieAppId;references:AppId"`
	Genres          []Genre          `gorm:"foreignKey:IndieAppId;references:AppId"`
}

type Movie struct {
	MovieId      uint `gorm:"primarykey;column:movie_id"`
	Name         string
	Mp4          string
	SteamMovieId uint
	IndieAppId   uint
}

type Screenshot struct {
	ScreenshotId      uint `gorm:"primarykey;column:screenshot_id"`
	PathThumbnail     string
	PathFull          string
	SteamScreenshotId uint
	IndieAppId        uint
}

type IndieAppDetail struct {
	AppDetailId      uint `gorm:"primarykey;column:app_detail_id"`
	Name             string
	ReleaseDate      string
	ShortDescription string
	Language         string
	IndieAppId       uint
}

type Genre struct {
	GenreId      uint `gorm:"primarykey;column:genre_id"`
	Description  string
	Language     string
	SteamGenreId uint
	IndieAppId   uint
}
