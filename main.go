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

	err = models.AutoMigrate(db)
	if err != nil {
		log.Fatal(err)
	}

	gs := models.NewGalleryService(db)
	ims := models.NewImageService(db)

	gh := handlers.NewGalleryHandler(gs)
	imh := handlers.NewImageHandler(gs, ims)

	r := gin.Default()
	r.Static("/upload", "./upload")

	r.POST("/galleries", gh.Create)
	r.GET("/galleries", gh.List)
	r.GET("/galleries/:id", gh.GetOne)
	r.DELETE("/galleries/:id", gh.Delete)
	r.PATCH("/galleries/:id/names", gh.UpdateName)
	r.PATCH("/galleries/:id/publishes", gh.UpdatePublishing)

	r.POST("/galleries/:id/images", imh.CreateImage)
	r.GET("/galleries/:id/images", imh.ListGalleryImages)
	r.DELETE("/images/:id", imh.DeleteImage)

	r.Run()

}
