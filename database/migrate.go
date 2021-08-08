package database

import "gorm.io/gorm"

type constraintChild struct {
	KeyName string
	FkName  string
}

func Migrate(db *gorm.DB, dt dbType) {
	if dt == MYSQL {
		db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	migrator := db.Migrator()

	createTable(migrator, []interface{}{&IndieApp{}, &Movie{}, &Screenshot{}, &IndieAppDetail{}})

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
