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

func addPlayerGetDataRoute(webservice *restful.WebService, handler GetPlayerDataHandler) {
	webservice.
		Route(webservice.GET(playerRoutePath+"/player/data").
			To(bindGetPlayerDataHandler(handler)).
			Operation("PlayerData").
			Doc("Get Data of a Player").
			Notes(heredoc.Doc(`
			Get data associated with the user
		`)).
			Reads(PlayerDataGetRequest{}).
			Returns(http.StatusCreated, http.StatusText(http.StatusCreated), nil).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), response.Error{}))
}

func bindGetPlayerDataHandler(handler GetPlayerDataHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		action := constants.ActionGetPlayerData
		additionalMessage := make(map[string]string)

		var dataRequest PlayerDataGetRequest
		err := req.ReadEntity(&dataRequest)
		if err != nil {
			errorCode := appmessage.EIDUnableToParseRequestBody
			errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
			logrus.Error(errorMessage)
			response.RespondError(req, resp, errorCode, http.StatusBadRequest, action, additionalMessage, err)
			return
		}

		err = dataRequest.Validate()
		if err != nil {
			errorCode := appmessage.EIDValidationError
			errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
			logrus.Error(errorMessage)
			response.RespondError(req, resp, errorCode, http.StatusBadRequest, action, additionalMessage, err)
			return
		}

		playerDataResponse, err := handler.GetPlayerData(dataRequest)
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
		response.Write(req, resp, http.StatusOK, appmessage.EIDGetPlayerSuccess, playerDataResponse)

	}
}
