package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("postgres", "host=postgres port=5432 user=postgres password=orrr-14 dbname=postgres sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
}
