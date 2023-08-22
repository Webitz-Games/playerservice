package requesthandlers

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"playerapi/pkg/config"
	"playerapi/pkg/player/api"
	"playerapi/pkg/utils"
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

func (p PlayerServiceRequestHandlers) HandleCreatePlayer(player api.Player) (api.PlayerResponse, error) {

	var result api.Player
	var playerResponse api.PlayerResponse
	filter := bson.D{{Key: "playerconfig.email", Value: player.Email}}
	singleResult := p.mongoClient.Database(p.config.MongoDatabase).Collection(playerCollectionPrefix).FindOne(context.Background(), filter)
	err := singleResult.Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			player.PlayerID = utils.GenerateUUID()
			_, err = p.mongoClient.Database(p.config.MongoDatabase).Collection(playerCollectionPrefix).InsertOne(context.Background(), player)
			logrus.WithField("id", player.PlayerID).Info("player created")
			playerResponse.PlayerID = player.PlayerID
			return playerResponse, nil
		}
	}
	logrus.WithField("email", player.Email).Warning("failed to create player because an email already exists")
	return playerResponse, api.ErrConflict
}

func (p PlayerServiceRequestHandlers) HandleUpdatePlayer(player api.Player) error {
	opts := options.Update().SetUpsert(false)
	filter := bson.D{{Key: "_id", Value: player.PlayerID}}
	update := bson.D{{Key: "$set", Value: player}}

	result, err := p.mongoClient.Database(p.config.MongoDatabase).Collection(playerCollectionPrefix).UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		logrus.WithField("id", player.PlayerID).Errorf("failed to save to mongo %s", err)
		return err
	}

	if result.MatchedCount == 0 {
		logrus.WithField("player", player.PlayerName).Warningf("failed to find a player with the id %s to update", player.PlayerID)
		return api.NewErrNotFound("player with id " + player.PlayerID)
	}

	return nil
}

func (p PlayerServiceRequestHandlers) HandleDeletePlayer(playerID string) error {
	filter := bson.D{{Key: "_id", Value: playerID}}
	result, err := p.mongoClient.Database(p.config.MongoDatabase).Collection(playerCollectionPrefix).DeleteOne(context.Background(), filter)
	if err != nil {
		logrus.WithField("id", playerID).Errorf("failed to delete player %s, %s", playerID, err)
		return err
	}

	if result.DeletedCount == 0 {
		logrus.WithField("id", playerID).Errorf("failed to delete player %s, a matching resource was not found", playerID)
		return api.NewErrNotFound("player with id " + playerID)
	}

	return nil
}

func (p PlayerServiceRequestHandlers) HandlePlayerLogin(playerConfig api.PlayerConfig) error {

	return nil
}

func (p PlayerServiceRequestHandlers) GetPlayer(playerName string) {

}
