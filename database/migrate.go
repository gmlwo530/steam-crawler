package database

import (
	"github.com/gmlwo530/steam-crawler/config"
	"gorm.io/gorm"
)

type constraintChild struct {
	KeyName string
	FkName  string
}

func Migrate(db *gorm.DB, dt dbType) {
	if dt == MYSQL {
		db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	migrator := db.Migrator()

	tables := []interface{}{&IndieApp{}, &Movie{}, &Screenshot{}, &IndieAppDetail{}, &Genre{}}

	if config.GetConfig().Debug {
		dropTables(migrator, tables)
	}

	createTable(migrator, tables)

	createConstraint(migrator, &IndieApp{}, []constraintChild{
		{
			KeyName: "Movies",
			FkName:  "fk_indie_apps_movies",
		},
		{
			KeyName: "Screenshots",
			FkName:  "fk_indie_apps_screenshots",
		},
		{
			KeyName: "IndieAppDetails",
			FkName:  "fk_indie_apps_indie_app_details",
		},
		{
			KeyName: "Genres",
			FkName:  "fk_indie_apps_genres",
		},
	})
}

func createTable(migrator gorm.Migrator, tables []interface{}) {
	for _, table := range tables {
		if !migrator.HasTable(table) {
			migrator.CreateTable(table)
		}
	}
}

func createConstraint(migrator gorm.Migrator, parent interface{}, childs []constraintChild) {
	for _, child := range childs {
		if !migrator.HasConstraint(parent, child.KeyName) {
			migrator.CreateConstraint(parent, child.KeyName)
		}

		if !migrator.HasConstraint(parent, child.FkName) {
			migrator.CreateConstraint(parent, child.FkName)
		}
	}
}

func dropTables(migrator gorm.Migrator, tables []interface{}) {
	for _, table := range tables {
		if migrator.HasTable(table) {
			migrator.DropTable(table)
		}
	}
}
