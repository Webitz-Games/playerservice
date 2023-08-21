package requesthandlers

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
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
		logrus.WithField("player", player.Name).Errorf("failed to save to mongo %s", err.Error())
		return err
	}

	return nil
}

func (p PlayerServiceRequestHandlers) HandleUpdatePlayer() error {
	//TODO implement me
	panic("implement me")
}

func (p PlayerServiceRequestHandlers) HandleDeletePlayer() error {
	//TODO implement me
	panic("implement me")
}
