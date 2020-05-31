package main

import (
	"gallery-api/handlers"
	"gallery-api/hash"
	"gallery-api/models"
	"gallery-api/mw"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const hmacKey = "secret"

func main() {
	db, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3307)/gallerydb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.LogMode(true) // dev only!

	// err = models.AutoMigrate(db)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	hmac := hash.NewHMAC(hmacKey)
	gs := models.NewGalleryService(db)
	ims := models.NewImageService(db)
	us := models.NewUserService(db, hmac)

	gh := handlers.NewGalleryHandler(gs)
	imh := handlers.NewImageHandler(gs, ims)
	uh := handlers.NewUserHandler(us)

	r := gin.Default()

	r.Static("/upload", "./upload")

	r.POST("/signup", uh.Signup)
	r.POST("/login", uh.Login)

	auth := r.Group("/")
	auth.Use(mw.RequireUser(us))
	{
		auth.POST("/logout", uh.Logout)
		auth.POST("/galleries", gh.Create)
		auth.GET("/galleries", gh.List)
		auth.GET("/galleries/:id", gh.GetOne)
		auth.DELETE("/galleries/:id", gh.Delete)
		auth.PATCH("/galleries/:id/names", gh.UpdateName)
		auth.PATCH("/galleries/:id/publishes", gh.UpdatePublishing)
		auth.POST("/galleries/:id/images", imh.CreateImage)
		auth.GET("/galleries/:id/images", imh.ListGalleryImages)
		auth.DELETE("/images/:id", imh.DeleteImage)
	}

	r.Run(":8080")

}
