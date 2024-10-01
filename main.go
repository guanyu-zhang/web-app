package main

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"log"
	"net/http"
	"web-app-backend/initializers"
	"web-app-backend/middlewares"
	"web-app-backend/routers"
)

var (
	port string
)

func init() {
	port = "8000"
	initializers.SyncDatabase()
}

func main() {
	initializers.ConnectToDb()
	engine := gin.Default()
	m := melody.New()
	gin.SetMode(gin.DebugMode)
	authMiddleware := middlewares.GetAuthMiddleWare()

	engine.Use(middlewares.HandlerMiddleWare(authMiddleware))

	routers.SetupRouter(engine, authMiddleware, m)

	if err := http.ListenAndServe("0.0.0.0:"+port, engine); err != nil {
		log.Fatal(err)
	}
}
