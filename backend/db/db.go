package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/HexCardGames/HexDeck/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DatabaseConnection struct {
	client *mongo.Client
}

type GlobalStatsCollection struct {
	GamesPlayed int `bson:"games_played"`
}

func (conn *DatabaseConnection) QueryRunningRooms() []*types.Room {
	res, err := conn.client.Database("hexdeck").Collection("games").Find(context.TODO(), bson.D{{Key: "gamestate", Value: bson.D{{Key: "$ne", Value: types.StateEnded}}}})
	if err != nil {
		slog.Error("Loading rooms from database failed", "error", err)
		return make([]*types.Room, 0)
	}

	var serializableRooms []SerializableRoom
	err = res.All(context.TODO(), &serializableRooms)
	if err != nil {
		slog.Error("Decoding rooms from database failed", "error", err)
		return make([]*types.Room, 0)
	}
	var rooms []*types.Room = make([]*types.Room, len(serializableRooms))
	for i, serializableRoom := range serializableRooms {
		room := serializableRoom.ToRoom()
		rooms[i] = room
	}
	return rooms
}

func (conn *DatabaseConnection) InsertRoom(room *types.Room) {
	_, err := conn.client.Database("hexdeck").Collection("games").InsertOne(context.TODO(), room)
	if err != nil {
		slog.Error("Error while inserting room into database", "error", err)
	}
}

func (conn *DatabaseConnection) UpdateRoom(room *types.Room) {
	result, err := conn.client.Database("hexdeck").Collection("games").UpdateByID(context.TODO(), room.RoomId, bson.D{{Key: "$set", Value: room}})
	if err != nil {
		slog.Error("Error while updating room in database", "error", err)
	}
	if result.MatchedCount < 1 {
		slog.Warn(fmt.Sprintf("No collections were found while trying to update room data for room '%s'", room.RoomId))
	}
}

func (conn *DatabaseConnection) IncrementGamesPlayed() {
	conn.client.Database("hexdeck").Collection("global_stats").UpdateOne(context.TODO(), bson.D{}, bson.D{
		{Key: "$inc", Value: bson.D{{Key: "games_played", Value: 1}}},
	}, options.UpdateOne().SetUpsert(true))
}

func (conn *DatabaseConnection) QueryGlobalStats() GlobalStatsCollection {
	res := conn.client.Database("hexdeck").Collection("global_stats").FindOne(context.TODO(), bson.D{})
	var stats GlobalStatsCollection
	res.Decode(&stats)
	return stats
}

func CreateDBConnection(uri string) DatabaseConnection {
	client, _ := mongo.Connect(options.Client().ApplyURI(uri))
	return DatabaseConnection{client}
}

var Conn DatabaseConnection

func InitDB(uri string) {
	Conn = CreateDBConnection(uri)
}
