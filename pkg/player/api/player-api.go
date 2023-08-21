package api

import "github.com/emicklei/go-restful/v3"

const (
	playerRoutePath = "/player"
)

type Player struct {
	Name string `json:"name" bson:"_id"`
	PlayerConfig
}

type PlayerConfig struct {
	Name     string
	Password string
}

type PlayerRequestHandlers interface {
	CreatePlayerHandler
	UpdatePlayerHandler
	DeletePlayerHandler
}

func (p Player) Validate() error {
	if p.Name == "" {
		return NewInvalidErr("player name cannot be empty")
	}
	return nil
}

func (pc *PlayerConfig) Validate() error {
	if pc.Name == "" {
		return NewInvalidErr("player name cannot be empty")
	}
	if pc.Password == "" {
		return NewInvalidErr("player password cannot be empty")
	}
	return nil
}

type DeletePlayerHandler interface {
	HandleDeletePlayer() error
}

type UpdatePlayerHandler interface {
	HandleUpdatePlayer() error
}

type CreatePlayerHandler interface {
	HandleCreatePlayer() error
}

func addPlayerRoutes(webservice *restful.WebService, handlers PlayerRequestHandlers) {
	addPlayerCreateRoute(webservice, handlers)
}
