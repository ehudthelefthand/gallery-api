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

	err = models.AutoMigrate(db)
	if err != nil {
		log.Fatal(err)
	}

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
	r.GET("/galleries", gh.ListPublish)

	auth := r.Group("/")
	auth.Use(mw.RequireUser(us))
	{
		auth.POST("/logout", uh.Logout)
		admin := auth.Group("/admin")
		{
			admin.POST("/galleries", gh.Create)
			admin.GET("/galleries", gh.List)
			admin.GET("/galleries/:id", gh.GetOne)
			admin.DELETE("/galleries/:id", gh.Delete)
			admin.PATCH("/galleries/:id/names", gh.UpdateName)
			admin.PATCH("/galleries/:id/publishes", gh.UpdatePublishing)
			admin.POST("/galleries/:id/images", imh.CreateImage)
			admin.GET("/galleries/:id/images", imh.ListGalleryImages)
			admin.DELETE("/images/:id", imh.DeleteImage)
		}

	}

	r.Run(":8080")

}
