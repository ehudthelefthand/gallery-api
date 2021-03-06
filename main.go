package main

import (
	"gallery-api/config"
	"gallery-api/handlers"
	"gallery-api/hash"
	"gallery-api/models"
	"gallery-api/mw"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	conf := config.Load()

	db, err := gorm.Open("mysql", conf.Connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if conf.Mode == "dev" {
		db.LogMode(true) // dev only!
	}

	err = models.AutoMigrate(db)
	if err != nil {
		log.Fatal(err)
	}

	hmac := hash.NewHMAC(conf.HMACKey)
	gs := models.NewGalleryService(db)
	ims := models.NewImageService(db)
	us := models.NewUserService(db, hmac)

	gh := handlers.NewGalleryHandler(gs)
	imh := handlers.NewImageHandler(gs, ims)
	uh := handlers.NewUserHandler(us)

	if conf.Mode != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
