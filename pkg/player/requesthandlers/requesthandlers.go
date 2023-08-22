package requesthandlers

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"playerapi/pkg/config"
	"playerapi/pkg/player/api"
	"playerapi/pkg/session"
	"playerapi/pkg/utils"
)

const (
	playerCollectionPrefix = "player_"
)

type PlayerServiceRequestHandlers struct {
	mongoClient *mongo.Client
	config      *config.Config
	session     session.SessionService
}

func MakeRequestHandlers(mongoClient *mongo.Client, config *config.Config, session session.SessionService) PlayerServiceRequestHandlers {
	return PlayerServiceRequestHandlers{mongoClient: mongoClient, config: config, session: session}
}

func (p PlayerServiceRequestHandlers) HandleCreatePlayer(playerCreateRequest api.PlayerCreateRequest) (api.PlayerCreateResponse, error) {

	var newPlayer api.Player
	var playerResponse api.PlayerCreateResponse
	filter := bson.D{{Key: "email", Value: playerCreateRequest.Email}}
	singleResult := p.mongoClient.Database(p.config.MongoDatabase).Collection(playerCollectionPrefix).FindOne(context.Background(), filter)
	err := singleResult.Decode(&newPlayer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			newPlayer.PlayerID = utils.GenerateUUID()
			newPlayer.PlayerName = playerCreateRequest.PlayerName
			newPlayer.Email = playerCreateRequest.Email
			newPlayer.Password = playerCreateRequest.Password
			_, err = p.mongoClient.Database(p.config.MongoDatabase).Collection(playerCollectionPrefix).InsertOne(context.Background(), newPlayer)
			logrus.WithField("id", newPlayer.PlayerID).Info("player created")
			playerResponse.PlayerID = newPlayer.PlayerID
			return playerResponse, nil
		}
	}
	logrus.WithField("email", playerCreateRequest.Email).Warning("failed to create player because an email already exists")
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

func (p PlayerServiceRequestHandlers) HandlePlayerLogin(loginRequest api.PlayerLoginRequest) (api.PlayerLoginResponse, error) {
	var response api.PlayerLoginResponse
	player, err := p.getPlayer(loginRequest.Email)
	if err != nil {
		return response, err
	}

	if player.Password != loginRequest.Password {
		return response, api.NewInvalidErr("password did not match")
	}

	sessionID, err := p.session.CreateSession(player)
	if err != nil {
		return api.PlayerLoginResponse{}, err
	}

	player.SessionID = sessionID
	response.Player = player

	return response, nil
}

func (p PlayerServiceRequestHandlers) getPlayer(email string) (api.Player, error) {
	var result api.Player
	filter := bson.D{{"email", email}}
	singleResult := p.mongoClient.Database(p.config.MongoDatabase).Collection(playerCollectionPrefix).FindOne(context.Background(), filter)
	err := singleResult.Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.WithField("email", email).Warning("failed to find user email")
			return result, err
		}
	}

	return result, nil
}
