package setting

import (
	"strings"

	"../service"

	"github.com/gentwolf-shen/gohelper/ginhelper"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(c.GetHeader("Authorization"))
		appKey := service.Auth.CheckToken(token)
		if appKey == "" {
			ginhelper.ShowNoAuth(c)
			c.Abort()
			return
		}

		c.Set("appKey", appKey)
	}
}
