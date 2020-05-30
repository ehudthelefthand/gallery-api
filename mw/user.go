package mw

import (
	"gallery-api/header"
	"gallery-api/models"

	"github.com/gin-gonic/gin"
)

func RequireUser(us models.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := header.GetToken(c)
		if token == "" {
			c.Status(401)
			c.Abort()
			return
		}
		user, err := us.GetByToken(token)
		if err != nil {
			c.Status(401)
			c.Abort()
			return
		}
		c.Set("user", user)
	}
}
