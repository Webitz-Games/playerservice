package api

import (
	"errors"
	"github.com/MakeNowJust/heredoc"
	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
	"net/http"
	"playerapi/pkg/appmessage"
	"playerapi/pkg/constants"
	"playerapi/pkg/response"
)

func addPlayerUpdateRoute(webservice *restful.WebService, handler UpdatePlayerHandler) {
	webservice.
		Route(webservice.PUT(playerRoutePath+"/{"+playerPathParameter+"}").
			Param(webservice.
				PathParameter(playerPathParameter, "PlayerID of the player").
				DataType("string").Required(true)).
			To(bindUpdatePlayerHandler(handler)).
			Operation("UpdatePlayer").
			Doc("Update a Player").
			Notes(heredoc.Doc(`
				Updates a Player
			`)).
			Reads(PlayerConfig{}).
			Returns(http.StatusCreated, http.StatusText(http.StatusCreated), Player{}).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), response.Error{}).
			Returns(http.StatusNotFound, http.StatusText(http.StatusForbidden), response.Error{}).
			Returns(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), response.Error{}))
}

func bindUpdatePlayerHandler(handler UpdatePlayerHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		action := constants.ActionUpdatePlayer
		additionalMessage := make(map[string]string)
		var playerConfig PlayerConfig

		playerID := req.PathParameter(playerPathParameter)

		err := req.ReadEntity(&playerConfig)
		if err != nil {
			errorCode := appmessage.EIDUnableToParseRequestBody
			errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
			logrus.Error(errorMessage)
			response.RespondError(req, resp, errorCode, http.StatusBadRequest, action, additionalMessage, err)
			return
		}

		err = playerConfig.Validate()
		if err != nil {
			errorCode := appmessage.EIDValidationError
			errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
			logrus.Error(errorMessage)
			response.RespondError(req, resp, errorCode, http.StatusBadRequest, action, additionalMessage, err)
			return
		}

		player := Player{
			PlayerID:     playerID,
			PlayerConfig: playerConfig,
		}
		err = handler.HandleUpdatePlayer(player)
		if err != nil {
			var notFoundErr *ErrNotFound
			switch {
			case errors.As(err, &notFoundErr):
				errorCode := appmessage.EIDUserNotFound
				errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
				logrus.Error(errorMessage)
				response.RespondError(req, resp, errorCode, http.StatusNotFound, action, additionalMessage, err)
				return
			default:
				errorCode := appmessage.EIDInternalServerError
				errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
				logrus.Error(errorMessage)
				response.RespondError(req, resp, errorCode, http.StatusInternalServerError, action, additionalMessage, err)
				return
			}
		}
		response.Write(req, resp, http.StatusOK, appmessage.EIDCreatePlayerSuccess, "updated player")
	}
}
