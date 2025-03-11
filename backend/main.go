package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/HexCardGames/HexDeck/api"
	"github.com/HexCardGames/HexDeck/db"
	"github.com/HexCardGames/HexDeck/game"
	"github.com/HexCardGames/HexDeck/utils"
)

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

	api.InitApi()
}
