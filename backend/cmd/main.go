package main

import (
	"log"

	"github.com/vandit1604/go-realtime-chat/db"
	"github.com/vandit1604/go-realtime-chat/internal/user"
	"github.com/vandit1604/go-realtime-chat/internal/ws"
	"github.com/vandit1604/go-realtime-chat/router"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatal("couldn't initialize database connection", err)
	}

	userRepo := user.NewRepository(dbConn.GetDb())
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)

	hub := ws.NewHub()
	wshandler := ws.NewHandler(hub)
	go hub.Run()

	router.InitRouter(userHandler, wshandler)
	router.Start("0.0.0.0:8080")
}
