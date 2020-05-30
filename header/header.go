package header

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	header = strings.TrimSpace(header)
	min := len("Bearer ")
	if len(header) <= min {
		return ""
	}
	return header[min:]
}
