package api

import "github.com/emicklei/go-restful/v3"

const (
	playerRoutePath     = "/players"
	playerPathParameter = "player"
)

type Player struct {
	Name string `json:"name" bson:"_id"`
	PlayerConfig
}

type PlayerConfig struct {
	PlayerName string `json:"player_name"`
	Password   string `json:"password"`
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
	return p.PlayerConfig.Validate()
}

func (pc *PlayerConfig) Validate() error {
	if pc.PlayerName == "" {
		return NewInvalidErr("player name cannot be empty")
	}
	if pc.Password == "" {
		return NewInvalidErr("player password cannot be empty")
	}
	return nil
}

type DeletePlayerHandler interface {
	HandleDeletePlayer(playerName string) error
}

type UpdatePlayerHandler interface {
	HandleUpdatePlayer(player Player) error
}

type CreatePlayerHandler interface {
	HandleCreatePlayer(player Player) error
}

func addPlayerRoutes(webservice *restful.WebService, handlers PlayerRequestHandlers) {
	addPlayerCreateRoute(webservice, handlers)
	addPlayerDeleteRoute(webservice, handlers)
	addPlayerUpdateRoute(webservice, handlers)

}
