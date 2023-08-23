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

func addPlayerLoginRoute(webservice *restful.WebService, handler LoginPlayerHandler) {
	webservice.
		Route(webservice.POST(playerRoutePath+"/player/login").
			To(bindLoginPlayerHandler(handler)).
			Operation("LoginPlayer").
			Doc("Login a player").
			Notes(heredoc.Doc(`
			Login a player
		`)).
			Reads(PlayerLoginRequest{}).
			Returns(http.StatusCreated, http.StatusText(http.StatusCreated), nil).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), response.Error{}))
}

func bindLoginPlayerHandler(handler LoginPlayerHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		action := constants.ActionLoginPlayer
		additionalMessage := make(map[string]string)

		var loginRequest PlayerLoginRequest
		err := req.ReadEntity(&loginRequest)
		if err != nil {
			errorCode := appmessage.EIDUnableToParseRequestBody
			errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
			logrus.Error(errorMessage)
			response.RespondError(req, resp, errorCode, http.StatusBadRequest, action, additionalMessage, err)
			return
		}

		err = loginRequest.Validate()
		if err != nil {
			errorCode := appmessage.EIDValidationError
			errorMessage := response.ConstructErrorMessage(action, constants.ErrorCodeMapping[errorCode], additionalMessage)
			logrus.Error(errorMessage)
			response.RespondError(req, resp, errorCode, http.StatusBadRequest, action, additionalMessage, err)
			return
		}

		//TODO check errors here on what we return
		loginResponse, err := handler.HandlePlayerLogin(loginRequest)
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
		response.Write(req, resp, http.StatusOK, appmessage.EIDPlayerLoginSuccess, loginResponse)

	}
}
