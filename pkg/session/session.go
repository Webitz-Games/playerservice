package session

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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

	return session.SessionID, nil
}

func (s *SessionService) GetSession(sessionID string) (Session, error) {
	var session Session
	filter := bson.D{{Key: "_id", Value: sessionID}}
	result := s.mongoClient.Database(s.config.MongoDatabase).Collection(sessionCollectionPrefix).FindOne(context.Background(), filter)
	err := result.Decode(&session)
	if err != nil {
		return Session{}, err
	}

	if !s.ValidSession(session.ExpirationTime) {
		return Session{}, errors.New("session is expired")
	}

	return session, nil
}

func (s *SessionService) ValidSession(sessionExpiration time.Time) bool {
	if sessionExpiration.Before(time.Now()) {
		return false
	}
	return true
}
