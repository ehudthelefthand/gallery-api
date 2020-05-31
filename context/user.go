package context

import (
	"gallery-api/models"

	"github.com/gin-gonic/gin"
)

const key string = "user"

func SetUser(c *gin.Context, user *models.User) {
	c.Set(key, user)
}

func User(c *gin.Context) *models.User {
	user, ok := c.Value(key).(*models.User)
	if !ok {
		return nil
	}
	return user
}
