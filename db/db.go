package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func GetDB() *sql.DB {
	db, err := sql.Open("sqlite3", os.Getenv("SQLITE3_PATH"))

	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	create table if not exists indie_app (
		id integer not null primary key autoincrement,
		app_id integer not null,
		average_forever integer not null,
		ccu integer not null
	);
	create table if not exists app_detail (
		app_id integer not null primary key,
		name string,
		is_free boolean,
		detailed_description string,
		release_date string
	);
	create table if not exists genre (
		id integer not null primary key autoincrement,
		genre_id integer not null,
		description string not null,
		app_id integer not null,
		foreign key (app_id)
       		REFERENCES app_detail (app_id) 
	);
	create table if not exists screenshot (
		id integer not null primary key autoincrement,
		secreenshot_id integer not null,
		path_thumb_nail string,
		path_full string,
		app_id integer not null,
		foreign key (app_id)
       		REFERENCES app_detail (app_id) 
	);
	create table if not exists movie (
		id integer not null primary key autoincrement,
		movie_id integer not null,
		name string,
		mp4 string not null,
		app_id integer not null,
		foreign key (app_id)
       		REFERENCES app_detail (app_id) 
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func InsertIndieApp(db *sql.DB, indieApps []IndieApp) {
	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into indie_app(app_id, average_forever, ccu) values(?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	for _, indieApp := range indieApps {
		_, err = stmt.Exec(indieApp.AppId, indieApp.AverageForever, indieApp.Ccu)
		if err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit()
}

func SelectIndieApp(db *sql.DB, limit int, offset int) []IndieApp {
	indieApps := make([]IndieApp, 0, limit)

	rows, err := db.Query(fmt.Sprintf("SELECT app_id, average_forever, ccu FROM indie_app ORDER BY average_forever ASC LIMIT %d OFFSET %d", limit, offset))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var appId int
		var averageForever int
		var ccu int
		err = rows.Scan(&appId, &averageForever, &ccu)
		if err != nil {
			log.Fatal(err)
		}
		indieApps = append(indieApps, IndieApp{
			AppId:          appId,
			AverageForever: averageForever,
			Ccu:            ccu,
		})
	}

	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	return indieApps
}

func InsertAppDetail(db *sql.DB, appDetails []AppDetail) {
	var genres []AppGenre
	var screenshots []AppScreenshot
	var movies []AppMovie

	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into app_detail(app_id, name, is_free, detailed_description, release_date) values(?, ?, ?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	for _, appDetail := range appDetails {
		_, err = stmt.Exec(appDetail.AppId, appDetail.Name, appDetail.IsFree, appDetail.DetailedDescription, appDetail.ReleaseDate.Date)
		if err != nil {
			log.Fatal(err)
		}
		for _, genre := range appDetail.Genres {
			genre.AppId = appDetail.AppId
			genres = append(genres, genre)
		}
		for _, screenshot := range appDetail.Screenshots {
			screenshot.AppId = appDetail.AppId
			screenshots = append(screenshots, screenshot)
		}
		for _, movie := range appDetail.Movies {
			movie.AppId = appDetail.AppId
			movies = append(movies, movie)
		}
	}

	tx.Commit()

	insertGenre(db, genres)
	insertScreenshot(db, screenshots)
	insertMovie(db, movies)
}

func insertGenre(db *sql.DB, genres []AppGenre) {
	tx, err := db.Begin()

	stmt, err := tx.Prepare("insert into genre(genre_id, description, app_id) values(?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	for _, genre := range genres {
		_, err = stmt.Exec(genre.Id, genre.Description, genre.AppId)
		if err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit()
}

func insertScreenshot(db *sql.DB, screenshots []AppScreenshot) {
	tx, err := db.Begin()

	stmt, err := tx.Prepare("insert into screenshot(secreenshot_id, path_thumb_nail, path_full, app_id) values(?, ?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	for _, screenshot := range screenshots {
		_, err = stmt.Exec(screenshot.Id, screenshot.PathThumbnail, screenshot.PathFull, screenshot.AppId)
		if err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit()
}

func insertMovie(db *sql.DB, movies []AppMovie) {
	tx, err := db.Begin()

	stmt, err := tx.Prepare("insert into movie(movie_id, name, mp4, app_id) values(?, ?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	for _, movie := range movies {
		_, err = stmt.Exec(movie.Id, movie.Name, movie.Mp4["480"], movie.AppId)
		if err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit()
}
