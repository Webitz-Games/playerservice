package api

import "github.com/emicklei/go-restful/v3"

const (
	playerRoutePath = "/players"
	playerID        = "id"
)

type PlayerLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PlayerLoginResponse struct {
	SessionID string
	PlayerID  string
}

type Player struct {
	PlayerID   string `bson:"_id"`
	PlayerName string `json:"name"`
	Email      string `json:"email"`
	Password   int    `json:"password"`
	Data       PlayerData
}

type PlayerCreateRequest struct {
	PlayerName string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type PlayerCreateResponse struct {
	PlayerID string `json:"playerID"`
}

type PlayerDataRequest struct {
	PlayerID  string `json:"playerID"`
	SessionID string `json:"sessionID"`
}

type PlayerDataResponse struct {
	Data PlayerData
}

type PlayerData struct {
	Data interface{} `json:"data"`
}

type PlayerSaveDataRequest struct {
	PlayerID  string      `json:"playerID"`
	SessionID string      `json:"sessionID"`
	Data      interface{} `json:"data"`
}

type PlayerRequestHandlers interface {
	CreatePlayerHandler
	UpdatePlayerHandler
	DeletePlayerHandler
	LoginPlayerHandler
	GetPlayerDataHandler
	SavePlayerDataHandler
}

func (d *PlayerSaveDataRequest) Validate() error {
	if d.Data == nil {
		return NewInvalidErr("data cannot be empty")
	}
	return nil
}

func (r *PlayerDataRequest) Validate() error {
	if r.SessionID == "" {
		return NewInvalidErr("sessionID cannot be empty")
	}
	if r.PlayerID == "" {
		return NewInvalidErr("playerID cannot be empty")
	}
	return nil
}

func (lr *PlayerLoginRequest) Validate() error {
	if lr.Email == "" {
		return NewInvalidErr("email cannot be empty")
	}
	if lr.Password == "" {
		return NewInvalidErr("player password cannot be empty")
	}
	return nil
}

func (p *Player) Validate() error {
	if p.PlayerID == "" {
		return NewInvalidErr("player id cannot be empty")
	}
	return nil
}

func (pc *PlayerCreateRequest) Validate() error {
	if pc.Email == "" {
		return NewInvalidErr("email cannot be empty")
	}
	if pc.Password == "" {
		return NewInvalidErr("player password cannot be empty")
	}
	if pc.PlayerName == "" {
		return NewInvalidErr("player name cannot be empty")
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
	HandleCreatePlayer(playerCreateRequest PlayerCreateRequest) (PlayerCreateResponse, error)
}

type LoginPlayerHandler interface {
	HandlePlayerLogin(loginRequest PlayerLoginRequest) (PlayerLoginResponse, error)
}

type GetPlayerDataHandler interface {
	GetPlayerData(playerDataRequest PlayerDataRequest) (PlayerDataResponse, error)
}

type SavePlayerDataHandler interface {
	SavePlayerData(data PlayerSaveDataRequest) error
}

func addPlayerRoutes(webservice *restful.WebService, handlers PlayerRequestHandlers) {
	addPlayerCreateRoute(webservice, handlers)
	addPlayerDeleteRoute(webservice, handlers)
	addPlayerUpdateRoute(webservice, handlers)
	addPlayerLoginRoute(webservice, handlers)
	addPlayerGetDataRoute(webservice, handlers)
	addPlayerSaveDataRoute(webservice, handlers)

}
