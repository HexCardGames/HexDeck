package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/HexCardGames/HexDeck/api"
	"github.com/HexCardGames/HexDeck/db"
	"github.com/HexCardGames/HexDeck/game"
	"github.com/HexCardGames/HexDeck/utils"
	"github.com/gin-gonic/gin"
)

//go:embed all:public/*
var public embed.FS

func main() {
	logHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(logHandler))

	mongoUri := utils.Getenv("MONGO_URI", "")
	if mongoUri == "" {
		slog.Error("MONGO_URI environment variable not set!")
		return
	}
	ok := db.InitDB(mongoUri)
	if !ok {
		slog.Error("Initializing MongoDB database failed")
		return
	}
	game.LoadRooms()

	roomTicker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-roomTicker.C:
				game.TickRooms(1000)
			}
		}
	}()

	server := gin.Default()
	server.SetTrustedProxies(nil)

	api.RegisterApi(server)
	server.Use(api.SPAMiddleware(public, "public", "/"))

	listenHost := utils.Getenv("LISTEN_HOST", "0.0.0.0")
	listenPort, err := strconv.Atoi(utils.Getenv("LISTEN_PORT", "3000"))
	if err != nil {
		log.Fatal("Value of variable PORT is not a valid integer!")
	}
	slog.Info(fmt.Sprintf("HexDeck server listening on http://%s:%d", listenHost, listenPort))
	server.Run(fmt.Sprintf("%s:%d", listenHost, listenPort))
}
