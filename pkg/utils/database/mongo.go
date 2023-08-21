package database

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"playerapi/pkg/config"
	"time"
)

func NewMongoClient(config *config.Config) (*mongo.Client, error) {
	mongoURL := config.MongoURL
	opts := options.Client().
		ApplyURI(mongoURL)

	if config.MongoUserName != "" && config.MongoPassword != "" {
		credential := options.Credential{
			AuthSource: config.MongoDatabase,
			Username:   config.MongoUserName,
			Password:   config.MongoPassword,
		}
		opts = opts.SetAuth(credential)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.MongoContextTimeout)*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		logrus.Error("unable to connect to MongoDB: ", err)
		return nil, err
	}

	return mongoClient, nil
}
