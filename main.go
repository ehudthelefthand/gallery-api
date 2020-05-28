package main

import (
	"gallery-api/handlers"
	"gallery-api/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3307)/gallerydb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&models.Gallery{},
	).Error
	if err != nil {
		log.Fatal(err)
	}

	gs := models.NewGalleryService(db)
	gh := handlers.NewGalleryHandler(gs)

	r := gin.Default()

	r.POST("/galleries", gh.Create)
	r.GET("/galleries", gh.List)
	r.DELETE("/galleries/:id", gh.Delete)
	r.PATCH("/galleries/:id/names", gh.UpdateName)
	r.PATCH("/galleries/:id/publishes", gh.UpdatePublishing)
	r.Run()

}
