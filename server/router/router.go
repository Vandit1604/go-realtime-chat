package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vandit1604/go-realtime-chat/internal/user"
	"github.com/vandit1604/go-realtime-chat/internal/ws"
)

var engine *gin.Engine

func InitRouter(userHandler *user.Handler, wsHandler *ws.Handler) {
	engine = gin.Default()

	engine.POST("/signup", userHandler.CreateUser)
	engine.POST("/login", userHandler.Login)
	engine.GET("/logout", userHandler.Logout)

	engine.POST("/ws/createRoom", wsHandler.CreateRoom)
	engine.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	engine.GET("/ws/getRooms", wsHandler.GetRooms)
}

func Start(addr string) error {
	return engine.Run(addr)
}
