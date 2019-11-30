package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type rating struct {
	Filename string `gorm:"PRIMARY_KEY"`
	Rating int
}

var db *gorm.DB 

func initDB() error {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}
	defer db.Close()

	db.AutoMigrate(&rating{})

	return nil
}

func readRating(files []string) []rating{
	ratings := make([]rating, len(files))
	for i, f := range files{
		db.FirstOrCreate(&ratings[i], rating{Filename: f})
	}

	return ratings
}

func setRating(r rating) {
	db.Save(&r)
}