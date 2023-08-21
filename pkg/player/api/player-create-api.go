package api

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
	"net/http"
	"playerapi/pkg/appmessage"
	"playerapi/pkg/constants"
	"playerapi/pkg/response"
)

func addPlayerCreateRoute(webservice *restful.WebService, handler CreatePlayerHandler) {
	webservice.
		Route(webservice.POST(playerRoutePath).
			To(bindCreatePlayerHandler(handler)).
			Operation("CreatePlayer").
			Doc("Create a Player").
			Notes(heredoc.Doc(`
				Creates a new Player
			`)).
			Reads(Player{}).
			Returns(http.StatusCreated, http.StatusText(http.StatusCreated), nil).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), response.Error{}))
}

func bindCreatePlayerHandler(handler CreatePlayerHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		action := constants.ActionCreatePlayer
		additionalMessage := make(map[string]string)

		var newPlayer Player
		err := req.ReadEntity(&newPlayer)
		if err != nil {
			errorCode := appmessage.EIDUnableToParseRequestBody
			errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
			logrus.Error(errorMessage)
			response.RespondError(req, resp, errorCode, http.StatusBadRequest, action, additionalMessage, err)
			return
		}

		err = newPlayer.Validate()
		if err != nil {
			errorCode := appmessage.EIDValidationError
			errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
			logrus.Error(errorMessage)
			response.RespondError(req, resp, errorCode, http.StatusBadRequest, action, additionalMessage, err)
			return
		}

		handler.HandleCreatePlayer()

		response.Write(req, resp, http.StatusCreated, appmessage.EIDCreatePlayerSuccess, "created player")
	}
}
