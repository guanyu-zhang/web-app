package routers

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"log"
	"net/http"
	"web-app-backend/controllers"
	"web-app-backend/middlewares"
)

func SetupRouter(r *gin.Engine, authMiddleWare *jwt.GinJWTMiddleware, m *melody.Melody) *gin.Engine {
	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./templates/signup.html")
	})

	r.GET("/login", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./templates/login.html")
	})

	r.NoRoute(handleNoRoute())
	visitorGroup := r.Group("/visitor", middlewares.TxMiddleware())
	{
		visitorGroup.POST("/login", authMiddleWare.LoginHandler)
		visitorGroup.POST("/signup", controllers.SignUpHandler) // redirect to login through frontend
	}

	userGroup := r.Group("/auth", middlewares.TxMiddleware(), authMiddleWare.MiddlewareFunc())
	{
		userGroup.POST("/create_chat_room_for_two", controllers.CreateChatRoomForTwoHandler)
		userGroup.POST("/get_my_chat_groups", controllers.GetMyGroupsHandler)
		userGroup.POST("/upvote_message", controllers.UpVoteHandler)
		userGroup.POST("/downvote_message", controllers.DownVoteHandler)
		userGroup.POST("/view_current_messages", controllers.ViewCurrentMessagesHandler)
		userGroup.POST("logout", authMiddleWare.LogoutHandler)
		userGroup.GET("/refresh_token", authMiddleWare.RefreshHandler)
		userGroup.GET("/all_users", controllers.GetAllUsersHandler)
		userGroup.GET("/unknown_people", controllers.UnknownPeopleHandler)
		userGroup.GET("/hello", controllers.HelloHandler)
		userGroup.GET("/chats", func(c *gin.Context) {
			http.ServeFile(c.Writer, c.Request, "./templates/chats.html")
		})
		userGroup.GET("/channel/:chan", middlewares.PermissionMiddleware(), func(c *gin.Context) {
			http.ServeFile(c.Writer, c.Request, "./templates/chan.html")
		})
		userGroup.GET("/channel/:chan/ws",
			func(c *gin.Context) {
				c.Set("websocket-melody", m)
				err := m.HandleRequest(c.Writer, c.Request)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"status":  "fail",
						"message": "unknown error",
					})
					return
				}
			})
	}
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		err := m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
		})
		if err != nil {
			return
		}
		controllers.RecordMessage(string(msg))
	})

	return r
}

func handleNoRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}
