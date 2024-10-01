package middlewares

import (
	"log"
	"time"
	. "web-app-backend/controllers"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var (
	identityKey = "username"
	idKey       = "id"
)

func GetAuthMiddleWare() *jwt.GinJWTMiddleware {
	middleware, err := jwt.New(initParams())
	if err != nil {
		log.Fatal("Error when creating jwt middleware:" + err.Error())
	}
	return middleware
}

func HandlerMiddleWare(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(context *gin.Context) {
		errInit := authMiddleware.MiddlewareInit()
		if errInit != nil {
			log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		}
	}
}

func initParams() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: PayloadFunc(),

		IdentityHandler: IdentityHandler(),
		Authenticator:   Authenticator(),
		Authorizator:    Authorizator(),
		Unauthorized:    Unauthorized(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
		SendCookie:      true,
	}
}
