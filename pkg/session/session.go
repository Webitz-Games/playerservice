package session

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"playerapi/pkg/config"
	"playerapi/pkg/player/api"
	"playerapi/pkg/utils"
	"time"
)

const (
	sessionCollectionPrefix = "session_"
)

type Session struct {
	SessionID      string    `json:"sessionID" bson:"_id"`
	PlayerID       string    `json:"playerID"`
	ExpirationTime time.Time `json:"expirationTime"`
}

type SessionService struct {
	mongoClient *mongo.Client
	config      *config.Config
}

func New(mongo *mongo.Client, config *config.Config) SessionService {
	return SessionService{mongoClient: mongo, config: config}
}

func (s *SessionService) CreateSession(player api.Player) (string, error) {
	expiration := time.Now().Add(time.Hour * 24)
	session := Session{
		SessionID:      utils.GenerateUUID(),
		PlayerID:       player.PlayerID,
		ExpirationTime: expiration,
	}

	_, err := s.mongoClient.Database(s.config.MongoDatabase).Collection(sessionCollectionPrefix).InsertOne(context.Background(), session)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to create session for player %s", err))
	}

	player.SessionID = session.SessionID

	return session.SessionID, nil
}
