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

func addPlayerDeleteRoute(webservice *restful.WebService, handler DeletePlayerHandler) {
	webservice.
		Route(webservice.DELETE(playerRoutePath+"/{"+playerPathParameter+"}").
			Param(webservice.
				PathParameter(playerPathParameter, "Name of the player").
				DataType("string").Required(true)).
			To(bindDeletePlayerHandler(handler)).
			Operation("DeletePlayer").
			Doc("Delete a Player").
			Notes(heredoc.Doc(`
				Deletes a Player
			`)).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), nil).
			Returns(http.StatusNotFound, http.StatusText(http.StatusNotFound), response.Error{}).
			Returns(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), response.Error{}))
}

func bindDeletePlayerHandler(handler DeletePlayerHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		action := constants.ActionDeletePlayer
		additionalMessage := make(map[string]string)

		playerName := req.PathParameter(playerPathParameter)

		err := handler.HandleDeletePlayer(playerName)
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

		response.Write(req, resp, http.StatusNoContent, appmessage.EIDGetPlayerSuccess, nil)
	}
}
