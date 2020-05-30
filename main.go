package main

import (
	"gallery-api/handlers"
	"gallery-api/models"
	"gallery-api/mw"
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
	defer db.Close()

	db.LogMode(true) // dev only!

	err = models.AutoMigrate(db)
	if err != nil {
		log.Fatal(err)
	}

	gs := models.NewGalleryService(db)
	ims := models.NewImageService(db)
	us := models.NewUserService(db)

	gh := handlers.NewGalleryHandler(gs)
	imh := handlers.NewImageHandler(gs, ims)
	uh := handlers.NewUserHandler(us)

	r := gin.Default()

	r.Static("/upload", "./upload")

	r.POST("/signup", uh.Signup)
	r.POST("/login", uh.Login) // Success => 200, Fail => 401

	// auth := func(c *gin.Context) {
	// 	header := c.GetHeader("Authorization")
	// 	token := header[8:]
	// 	user = GetUserByToken(token)
	// 	if ไม่มี user {
	// 		// Bail out
	// 	}

	// }

	auth := r.Group("/")
	auth.Use(mw.RequireUser(us))
	{
		auth.POST("/logout", uh.Logout)
		auth.GET("/sessions", uh.GetSession)
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

	// r.Use()

	// xg := r.Group("/groupx")
	// xg.Use(func(c *gin.Context) {
	// 	fmt.Println("X")
	// })
	// {
	// 	xg.GET("/testx1", func(c *gin.Context) {
	// 		fmt.Println("x1-1")
	// 		// c.Set("abc", "xyz")
	// 		c.JSON(200, gin.H{
	// 			"message": "hello X1-1",
	// 		})
	// 	}, func(c *gin.Context) {
	// 		fmt.Println("x1-2")
	// 		// value := c.Get("abc") // xyz
	// 		c.JSON(200, gin.H{
	// 			"message": "hello X1-2",
	// 		})
	// 	})
	// }

	// yg := r.Group("/groupy")
	// yg.Use(func(c *gin.Context) {
	// 	fmt.Println("Y")
	// })
	// {
	// 	yg.GET("/testy1", func(c *gin.Context) {
	// 		fmt.Println("y1")
	// 		c.JSON(200, gin.H{
	// 			"message": "hello Y1",
	// 		})
	// 	})
	// }

	r.Run()

}
