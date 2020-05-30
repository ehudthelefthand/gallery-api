package models

import (
	"log"

	"github.com/jinzhu/gorm"
)

// Reset will drop and re-create table
func Reset(db *gorm.DB) error {
	err := db.DropTableIfExists(
		&Gallery{},
		&Image{},
	).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return AutoMigrate(db)
}

// AutoMigrate will create or update table
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Gallery{},
		&Image{},
		&User{},
	).Error
}
