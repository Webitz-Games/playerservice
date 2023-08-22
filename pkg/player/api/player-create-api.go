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

func addPlayerCreateRoute(webservice *restful.WebService, handler CreatePlayerHandler) {
	webservice.
		Route(webservice.POST(playerRoutePath).
			To(bindCreatePlayerHandler(handler)).
			Operation("CreatePlayer").
			Doc("Create a Player").
			Notes(heredoc.Doc(`
				Creates a new Player
			`)).
			Reads(PlayerCreateRequest{}).
			Returns(http.StatusCreated, http.StatusText(http.StatusCreated), nil).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), response.Error{}))
}

func bindCreatePlayerHandler(handler CreatePlayerHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		action := constants.ActionCreatePlayer
		additionalMessage := make(map[string]string)

		var newPlayerCreateRequest PlayerCreateRequest
		err := req.ReadEntity(&newPlayerCreateRequest)
		if err != nil {
			errorCode := appmessage.EIDUnableToParseRequestBody
			errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
			logrus.Error(errorMessage)
			response.RespondError(req, resp, errorCode, http.StatusBadRequest, action, additionalMessage, err)
			return
		}

		err = newPlayerCreateRequest.Validate()
		if err != nil {
			errorCode := appmessage.EIDValidationError
			errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
			logrus.Error(errorMessage)
			response.RespondError(req, resp, errorCode, http.StatusBadRequest, action, additionalMessage, err)
			return
		}

		playerResponse, err := handler.HandleCreatePlayer(newPlayerCreateRequest)
		if err != nil {
			var conflictErr *ErrResourceConflict
			switch {
			case errors.As(err, &conflictErr):
				errorCode := appmessage.EIDValidationError
				errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
				logrus.Error(errorMessage)
				response.RespondError(req, resp, errorCode, http.StatusConflict, action, additionalMessage, err)
				return
			default:
				errorCode := appmessage.EIDInternalServerError
				errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
				logrus.Error(errorMessage)
				response.RespondError(req, resp, errorCode, http.StatusInternalServerError, action, additionalMessage, err)
				return
			}
		}

		response.Write(req, resp, http.StatusCreated, appmessage.EIDCreatePlayerSuccess, playerResponse)
	}
}
