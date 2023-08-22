package api

import "github.com/emicklei/go-restful/v3"

const (
	playerRoutePath     = "/players"
	playerPathParameter = "player"
)

type Player struct {
	PlayerID string `bson:"_id"`
	PlayerConfig
}

type PlayerConfig struct {
	PlayerName string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Token      string `json:"token"`
}

type PlayerResponse struct {
	PlayerID string `json:"playerID"`
}

type PlayerRequestHandlers interface {
	CreatePlayerHandler
	UpdatePlayerHandler
	DeletePlayerHandler
	LoginPlayerHandler
}

//func (p Player) Validate() error {
//	if p.PlayerID == "" {
//		return NewInvalidErr("player name cannot be empty")
//	}
//	return p.PlayerConfig.Validate()
//}

func (pc *PlayerConfig) Validate() error {
	if pc.Email == "" {
		return NewInvalidErr("email cannot be empty")
	}
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
	HandleCreatePlayer(player Player) (PlayerResponse, error)
}

type LoginPlayerHandler interface {
	HandlePlayerLogin(player PlayerConfig) error
}

func addPlayerRoutes(webservice *restful.WebService, handlers PlayerRequestHandlers) {
	addPlayerCreateRoute(webservice, handlers)
	addPlayerDeleteRoute(webservice, handlers)
	addPlayerUpdateRoute(webservice, handlers)
	addPlayerLoginRoute(webservice, handlers)
}
