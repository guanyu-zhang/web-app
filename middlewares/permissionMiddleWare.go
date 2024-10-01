package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	. "web-app-backend/controllers"
	. "web-app-backend/utils"
)

func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tx, err := GetTx(c)
		if err != nil {
			log.Println("can't get transaction from gin context")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "unknown error",
			})
			return
		}
		chanParam := c.Param("chan")
		userId, _ := c.Get(idKey)
		chanId, err := strconv.ParseUint(chanParam, 10, 0)
		isAllowed, _ := IsUserAllowed(tx, userId.(uint), uint(chanId))
		log.Println("isallowed = %v", isAllowed)
		if !isAllowed {
			log.Println("not allowed to visit this page")
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "fail",
				"message": "you are not allowed to visit this page",
			})
			return
		}

		c.Next()
	}
}
