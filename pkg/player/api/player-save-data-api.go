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

func addPlayerSaveDataRoute(webservice *restful.WebService, handler SavePlayerDataHandler) {
	webservice.
		Route(webservice.POST(playerRoutePath+"/player/data").
			To(bindPlayerSaveDataRoute(handler)).
			Operation("PlayerData").
			Doc("Get Data of a Player").
			Notes(heredoc.Doc(`
			Get data associated with the user
		`)).
			Reads(PlayerSaveDataRequest{}).
			Returns(http.StatusCreated, http.StatusText(http.StatusCreated), nil).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), response.Error{}))
}

func bindPlayerSaveDataRoute(handler SavePlayerDataHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		action := constants.ActionSavePlayerData
		additionalMessage := make(map[string]string)

		var saveRequest PlayerSaveDataRequest
		err := req.ReadEntity(&saveRequest)
		if err != nil {
			errorCode := appmessage.EIDUnableToParseRequestBody
			errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
			logrus.Error(errorMessage)
			response.RespondError(req, resp, errorCode, http.StatusBadRequest, action, additionalMessage, err)
			return
		}

		err = saveRequest.Validate()
		if err != nil {
			errorCode := appmessage.EIDValidationError
			errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
			logrus.Error(errorMessage)
			response.RespondError(req, resp, errorCode, http.StatusBadRequest, action, additionalMessage, err)
			return
		}

		err = handler.SavePlayerData(saveRequest)
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
		response.Write(req, resp, http.StatusOK, appmessage.EIDGetPlayerSuccess, "player data saved")

	}
}
