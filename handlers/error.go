package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, err error) {
	log.Println(err)
	c.JSON(status, gin.H{
		"message": err.Error(),
	})
}
