package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)
import . "web-app-backend/utils"
import . "web-app-backend/initializers"

func TxMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := DB.Begin()

		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				log.Println("Recovered from panic:", r)
			} else if c.Writer.Status() >= http.StatusBadRequest {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()

		c.Set(TxKey, tx)

		c.Next()
	}
}
