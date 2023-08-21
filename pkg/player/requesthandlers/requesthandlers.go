package requesthandlers

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"playerapi/pkg/config"
	"playerapi/pkg/player/api"
)

const (
	playerCollectionPrefix = "player_"
)

type PlayerServiceRequestHandlers struct {
	mongoClient *mongo.Client
	config      *config.Config
}

func MakeRequestHandlers(mongoClient *mongo.Client, config *config.Config) PlayerServiceRequestHandlers {
	return PlayerServiceRequestHandlers{mongoClient: mongoClient, config: config}
}

func (p PlayerServiceRequestHandlers) HandleCreatePlayer(player api.Player) error {

	_, err := p.mongoClient.Database(p.config.MongoDatabase).Collection(playerCollectionPrefix).InsertOne(context.Background(), player)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			logrus.WithField("player", player.Name).Warningf("failed to create player because that player name already exists")
			return api.ErrConflict
		}
		logrus.WithField("player", player.Name).Errorf("failed to save to mongo %s", err)
		return err
	}

	return nil
}

func (p PlayerServiceRequestHandlers) HandleUpdatePlayer(player api.Player) error {
	opts := options.Update().SetUpsert(false)
	filter := bson.D{{Key: "_id", Value: player.Name}}
	update := bson.D{{Key: "$set", Value: player}}

	result, err := p.mongoClient.Database(p.config.MongoDatabase).Collection(playerCollectionPrefix).UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		logrus.WithField("player", player.Name).Errorf("failed to save to mongo %s", err)
		return err
	}

	if result.MatchedCount == 0 {
		logrus.WithField("player", player.Name).Warningf("failed to find a player with the name %s to update", player.Name)
		return api.NewErrNotFound("player with name " + player.Name)
	}

	return nil
}

func (p PlayerServiceRequestHandlers) HandleDeletePlayer(playerName string) error {
	filter := bson.D{{Key: "_id", Value: playerName}}
	result, err := p.mongoClient.Database(p.config.MongoDatabase).Collection(playerCollectionPrefix).DeleteOne(context.Background(), filter)
	if err != nil {
		logrus.WithField("player", playerName).Errorf("failed to delete player %s, %s", playerName, err)
		return err
	}

	if result.DeletedCount == 0 {
		logrus.WithField("player", playerName).Errorf("failed to delete player %s, a matching resource was not found", playerName)
		return api.NewErrNotFound("player with name " + playerName)
	}

	return nil
}
