package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/HexCardGames/HexDeck/game"
	"github.com/HexCardGames/HexDeck/types"
	"github.com/zishang520/socket.io/v2/socket"
	socketio "github.com/zishang520/socket.io/v2/socket"
)

var io *socketio.Server

func initWS() http.Handler {
	io = socketio.NewServer(nil, nil)

	io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		remoteAddr := client.Request().Request().RemoteAddr

		sessionToken, exists := client.Request().Query().Get("sessionToken")
		room, player := game.FindSession(sessionToken)
		if !exists || player == nil {
			slog.Debug("New WebSocket connection from didn't provide a valid sessionToken -> disconnecting", "remoteAddress", remoteAddr, "sessionToken", sessionToken)
			client.Emit("Status", types.S2C_Status{
				IsError:    true,
				StatusCode: "invalid_session",
				Message:    "No valid sessionToken was provided",
			})
			client.Disconnect(true)
			return
		}
		if player.Connection.IsConnected && player.Connection.Socket != nil {
			slog.Debug("User already connected to WebSocket -> disconnecting old socket", "remoteAddress", remoteAddr, "sessionToken", sessionToken)
			player.Connection.Socket.Emit("Status", types.S2C_Status{
				IsError:    true,
				StatusCode: "connection_from_different_socket",
				Message:    "User connected from a different socket",
			})
			player.Connection.Socket.Disconnect(true)
		}
		player.Connection.Socket = client
		player.Connection.IsConnected = true
		player.ResetInactivity()
		slog.Debug("New WebSocket connection", "username", player.Username, "remoteAddress", remoteAddr, "playerId", player.PlayerId, "sessionToken", sessionToken, "roomId", room.RoomId.Hex())
		game.OnRoomUpdate(room)

		onPlayerJoin(client, room, player)
	})

	return io.ServeHandler(nil)
}

func unpackData(datas []any, target interface{}) bool {
	if len(datas) < 1 {
		slog.Warn("Unexpected length of WebSocket data; ignoring message")
		return false
	}
	request, _ := datas[0].(string)
	ok := json.Unmarshal([]byte(request), &target)
	return ok != nil
}

func verifyPlayerIsActivePlayer(room *types.Room, target *types.Player) bool {
	if room.GameState != types.StateRunning {
		target.Connection.Socket.Emit("Status", types.S2C_Status{
			IsError:    true,
			StatusCode: "game_not_running",
			Message:    "The game is not running",
		})
		return false
	}

	if !room.CardDeck.IsPlayerActive(target) {
		target.Connection.Socket.Emit("Status", types.S2C_Status{
			IsError:    true,
			StatusCode: "player_not_active",
			Message:    "You can't execute this action while you are not the active player",
		})
		return false
	}
	return true
}

func onPlayerJoin(client *socket.Socket, room *types.Room, player *types.Player) {
	client.On("disconnect", func(...any) {
		player.Connection.IsConnected = false
		player.Connection.Socket = nil
		slog.Debug("Player disconnected from WebSocket", "username", player.Username, "remoteAddress", client.Conn().RemoteAddress(), "sessionToken", player.SessionToken, "roomId", room.RoomId.Hex())
		game.OnRoomUpdate(room)
	})

	client.On("UpdatePlayer", func(datas ...any) {
		updatePlayerRequest := types.C2S_UpdatePlayer{}
		unpackData(datas, &updatePlayerRequest)
		if updatePlayerRequest.PlayerId != player.PlayerId && !player.HasPermissionBit(types.PermissionHost) {
			client.Emit("Status", types.S2C_Status{
				IsError:    true,
				StatusCode: "insufficient_permission",
				Message:    "You can't update other users unless you are host",
			})
			return
		}
		targetPlayer := room.FindPlayer(updatePlayerRequest.PlayerId)
		if targetPlayer == nil {
			client.Emit("Status", types.S2C_Status{
				IsError:    true,
				StatusCode: "invalid_player",
				Message:    "No player with the requested playerId was found",
			})
			return
		}
		slog.Debug("Updating player data", "roomId", room.RoomId, "playerId", updatePlayerRequest.PlayerId, "username", targetPlayer.Username, "request", updatePlayerRequest)

		if updatePlayerRequest.Username != nil {
			if room.IsUsernameAvailable(*updatePlayerRequest.Username) {
				targetPlayer.Username = *updatePlayerRequest.Username
			} else {
				client.Emit("Status", types.S2C_Status{
					IsError:    true,
					StatusCode: "username_taken",
					Message:    "The requested username is not available",
				})
			}
		}
		if updatePlayerRequest.Permissions != nil {
			targetPlayer.Permissions = *updatePlayerRequest.Permissions
		}

		game.OnRoomUpdate(room)
	})

	client.On("KickPlayer", func(datas ...any) {
		kickPlayerRequest := types.C2S_KickPlayer{}
		unpackData(datas, &kickPlayerRequest)
		if !player.HasPermissionBit(types.PermissionHost) {
			client.Emit("Status", types.S2C_Status{
				IsError:    true,
				StatusCode: "insufficient_permission",
				Message:    "You can't update other users unless you are host",
			})
			return
		}
		targetPlayer := room.FindPlayer(kickPlayerRequest.PlayerId)
		if targetPlayer == nil {
			client.Emit("Status", types.S2C_Status{
				IsError:    true,
				StatusCode: "invalid_player",
				Message:    "No player with the requested playerId was found",
			})
			return
		}

		if room.RemovePlayer(*targetPlayer) {
			slog.Debug("Player was kicked from room", "playerId", player.PlayerId, "targetPlayerId", kickPlayerRequest.PlayerId, "roomId", room.RoomId)
			if targetPlayer.Connection.IsConnected && targetPlayer.Connection.Socket != nil {
				targetPlayer.Connection.Socket.Emit("Status", types.S2C_Status{
					IsError:    true,
					StatusCode: "player_kicked",
					Message:    "You were kicked from the room",
				})
			}
		}
		game.OnRoomUpdate(room)
	})

	client.On("StartGame", func(datas ...any) {
		if !player.HasPermissionBit(types.PermissionHost) {
			client.Emit("Status", types.S2C_Status{
				IsError:    true,
				StatusCode: "insufficient_permission",
				Message:    "You can't start the game unless you are host",
			})
			return
		}
		if room.GameState != types.StateLobby {
			client.Emit("Status", types.S2C_Status{
				IsError:    true,
				StatusCode: "game_already_started",
				Message:    "The game has already started",
			})
			return
		}
		game.StartGame(room)
	})

	client.On("DrawCard", func(datas ...any) {
		if !verifyPlayerIsActivePlayer(room, player) {
			return
		}
		card := room.CardDeck.DrawCard()
		if card == nil {
			// TODO: Handle empty card deck
			return
		}
		game.OnPlayerStateUpdate(room, player, false)
	})

	client.On("PlayCard", func(datas ...any) {
		if !verifyPlayerIsActivePlayer(room, player) {
			return
		}

		updatePlayerRequest := types.C2S_PlayCard{}
		unpackData(datas, &updatePlayerRequest)
		if updatePlayerRequest.CardIndex == nil {
			client.Emit("Status", types.S2C_Status{
				IsError:    true,
				StatusCode: "missing_parameter",
				Message:    "CardIndex parameter is missing",
			})
			return
		}
		if *updatePlayerRequest.CardIndex < 0 || *updatePlayerRequest.CardIndex >= len(player.Cards) {
			client.Emit("Status", types.S2C_Status{
				IsError:    true,
				StatusCode: "invalid_card_index",
				Message:    "Provided CardIndex is out of bounds",
			})
			return
		}
		card := player.Cards[*updatePlayerRequest.CardIndex]
		if !room.CardDeck.PlayCard(card) {
			client.Emit("Status", types.S2C_Status{
				IsError:    true,
				StatusCode: "card_not_playable",
				Message:    "You can't play this card now",
			})
			return
		}
		player.Cards = append(player.Cards[:*updatePlayerRequest.CardIndex], player.Cards[*updatePlayerRequest.CardIndex+1:]...)
		game.OnPlayCard(room, player, *updatePlayerRequest.CardIndex, card)

		if len(player.Cards) == 0 {
			room.Winner = &player.PlayerId
			game.UpdateGameState(room, types.StateEnded)
		}
	})

	client.On("UpdatePlayedCard", func(datas ...any) {
		if !verifyPlayerIsActivePlayer(room, player) {
			return
		}

		updatePlayerRequest := types.C2S_UpdatePlayedCard{}
		unpackData(datas, &updatePlayerRequest)
		card := room.CardDeck.UpdatePlayedCard(updatePlayerRequest.CardData)
		if card == nil {
			client.Emit("Status", types.S2C_Status{
				IsError:    true,
				StatusCode: "card_not_updatable",
				Message:    "You can't update this card now",
			})
			return
		}
		game.OnPlayedCardUpdate(room, player, card)
	})
}
