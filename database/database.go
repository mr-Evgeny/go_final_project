package database

import (
	"github.com/mr-Evgeny/go_final_project/config"
	"github.com/mr-Evgeny/go_final_project/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func Connect() {
	db, err := gorm.Open(sqlite.Open(config.DB_FILE), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database.\n", err)
		os.Exit(1)
	}

	db.Logger = logger.Default.LogMode(logger.Info)

	db.AutoMigrate(&model.ToDo{})

	DB = Dbinstance{
		Db: db,
	}
}
