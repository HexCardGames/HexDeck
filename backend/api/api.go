package api

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/HexCardGames/HexDeck/db"
	"github.com/HexCardGames/HexDeck/game"
	"github.com/HexCardGames/HexDeck/types"
	"github.com/HexCardGames/HexDeck/utils"
	"github.com/gin-gonic/gin"
)

type ErrorReply struct {
	StatusCode string
	Message    string
}
type StatsReply struct {
	TotalGamesPlayed  int
	RunningGames      int
	OnlinePlayerCount int
}
type ImprintReply struct {
	Content string
}
type CreateRoomReply struct {
	JoinCode string
}
type JoinRoomRequest struct {
	JoinCode string
	Username string
}
type LeaveRoomRequest struct {
	SessionToken string
}

func InitApi() {
	server := gin.Default()
	server.SetTrustedProxies(nil)

	server.GET("/api/stats", func(c *gin.Context) {
		stats := game.CalculateStats()
		c.JSON(http.StatusOK, StatsReply{
			TotalGamesPlayed:  db.Conn.QueryGlobalStats().GamesPlayed,
			RunningGames:      stats.RunningGames,
			OnlinePlayerCount: stats.OnlinePlayerCount,
		})
	})
	server.GET("/api/imprint", func(c *gin.Context) {
		// TODO: Implement imprint endpoint
		c.JSON(http.StatusOK, ImprintReply{
			Content: "Not implemented yet",
		})
	})

	server.POST("/api/room/create", func(c *gin.Context) {
		request := JoinRoomRequest{}
		c.BindJSON(&request)
		room := game.CreateRoom()
		player := game.JoinRoom(room, request.Username)
		player.SetPermissionBit(types.PermissionHost)
		slog.Debug("New room created", "username", player.Username, "sessionToken", player.SessionToken, "roomId", room.RoomId.Hex())
		c.JSON(http.StatusOK, player)
	})

	server.POST("/api/room/join", func(c *gin.Context) {
		request := JoinRoomRequest{}
		c.BindJSON(&request)
		room := game.FindRoomByJoinCode(request.JoinCode)
		if room == nil {
			slog.Debug("Client tried joining room using an invalid joinCode", "joinCode", request.JoinCode)
			c.JSON(http.StatusBadRequest, ErrorReply{
				StatusCode: "invalid_join_code",
				Message:    "No valid joinCode was provided",
			})
			return
		}
		if room.GameState != types.StateLobby {
			slog.Debug("Client tried joining room not in lobby state", "joinCode", request.JoinCode)
			c.JSON(http.StatusBadRequest, ErrorReply{
				StatusCode: "game_already_running",
				Message:    "You cannot join this room as the game has already started",
			})
			return
		}
		player := game.JoinRoom(room, request.Username)
		slog.Debug("New session created", "username", player.Username, "sessionToken", player.SessionToken, "roomId", room.RoomId.Hex(), "joinCode", request.JoinCode)
		c.JSON(http.StatusOK, player)
	})

	server.GET("/api/check/session", func(c *gin.Context) {
		sessionToken := c.Query("sessionToken")
		if sessionToken == "" {
			c.JSON(http.StatusBadRequest, ErrorReply{
				StatusCode: "missing_parameter",
				Message:    "Parameter sessionToken is missing",
			})
		}
		_, player := game.FindSession(sessionToken)
		if player == nil {
			c.Status(401)
		} else {
			c.Status(200)
		}
	})

	server.GET("/api/check/joinCode", func(c *gin.Context) {
		joinCode := c.Query("JoinCode")
		if joinCode == "" {
			c.JSON(http.StatusBadRequest, ErrorReply{
				StatusCode: "missing_parameter",
				Message:    "Parameter JoinCode is missing",
			})
		}
		room := game.FindRoomByJoinCode(joinCode)
		if room == nil {
			c.Status(401)
		} else {
			c.Status(200)
		}
	})

	server.POST("/api/room/leave", func(c *gin.Context) {
		request := LeaveRoomRequest{}
		c.BindJSON(&request)
		room, player := game.FindSession(request.SessionToken)
		if player == nil {
			c.JSON(http.StatusBadRequest, ErrorReply{
				StatusCode: "invalid_session",
				Message:    "No user was found with the provided sessionToken",
			})
			return
		}
		room.RemovePlayer(*player)
		game.OnRoomUpdate(room)
		c.Status(http.StatusOK)
	})

	// Handle WebSocket connections using Socket.io
	wsHandler := initWS()
	server.Any("/socket.io/", gin.WrapH(wsHandler))

	listenHost := utils.Getenv("LISTEN_HOST", "0.0.0.0")
	listenPort, err := strconv.Atoi(utils.Getenv("LISTEN_PORT", "3000"))
	if err != nil {
		log.Fatal("Value of variable PORT is not a valid integer!")
	}
	slog.Info(fmt.Sprintf("HexDeck server listening on http://%s:%d", listenHost, listenPort))
	server.Run(fmt.Sprintf("%s:%d", listenHost, listenPort))
}
