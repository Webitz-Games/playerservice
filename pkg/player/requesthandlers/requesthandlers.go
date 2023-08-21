package requesthandlers

import "go.mongodb.org/mongo-driver/mongo"

type PlayerServiceRequestHandlers struct {
	mongoClient *mongo.Client
}

func MakeRequestHandlers(mongoClient *mongo.Client) PlayerServiceRequestHandlers {
	return PlayerServiceRequestHandlers{mongoClient: mongoClient}
}

func (p PlayerServiceRequestHandlers) HandleCreatePlayer() error {
	//TODO implement me
	panic("implement me")
}

func (p PlayerServiceRequestHandlers) HandleUpdatePlayer() error {
	//TODO implement me
	panic("implement me")
}

func (p PlayerServiceRequestHandlers) HandleDeletePlayer() error {
	//TODO implement me
	panic("implement me")
}
