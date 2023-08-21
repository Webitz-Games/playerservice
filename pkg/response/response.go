package response

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
	"net/http"
	"playerapi/pkg/constants"
	"strings"
)

type Error struct {
	ErrorCode    int    `x-nullable:"false"`
	ErrorMessage string `x-nullable:"false"`
}

type ErrorV1 struct {
	ErrorCode        int               `json:"errorCode" x-nullable:"false"`
	ErrorMessage     string            `json:"errorMessage" x-nullable:"false"`
	MessageVariables map[string]string `json:"messageVariables" x-nullable:"false"`
}

// ConstructErrorMessage construct error message with this following format
// "unable to {action}: {reason}, user ID: {user_id}, {other details}: {other details} , . . ."
func ConstructErrorMessage(action, reason string, additionalMessage map[string]string) string {
	var additionalMsg string
	if len(additionalMessage) == 0 {
		additionalMsg = ""
	} else {
		for k, v := range additionalMessage {
			additionalMsg = additionalMsg + ", " + k + ": " + v
		}
	}
	userAction := action
	if strings.Contains(action, constants.Separator) {
		userAction = strings.ReplaceAll(userAction, constants.Separator, " ")
	}

	return fmt.Sprintf("unable to %s: %s, %s", userAction, reason, additionalMsg)
}

func RespondError(request *restful.Request, response *restful.Response, errorCode, httpStatus int, action string, additionalMessage map[string]string, err error) {
	msg := constants.ErrorCodeMapping[errorCode]
	if err != nil {
		if len(msg) > 0 {
			msg += ": "
		}
		msg += err.Error()
	}
	errorResponse := ConstructErrorResponse(errorCode, action, msg, additionalMessage)
	if e := response.WriteHeaderAndJson(httpStatus, errorResponse, restful.MIME_JSON); e != nil {
		logrus.Error(e)
	}
}

func ConstructErrorResponse(errorCode int, action, reason string, additionalMessage map[string]string) *ErrorV1 {
	return &ErrorV1{
		ErrorCode:        errorCode,
		ErrorMessage:     ConstructErrorMessage(action, reason, additionalMessage),
		MessageVariables: additionalMessage,
	}
}

func Write(request *restful.Request, response *restful.Response, httpStatusCode int, eventID int, entity interface{}) {
	err := response.WriteHeaderAndJson(httpStatusCode, entity, restful.MIME_JSON)
	if err != nil {
		WriteError(request, response, http.StatusInternalServerError, err,
			&Error{
				ErrorCode:    constants.UnableToWriteResponse,
				ErrorMessage: "unable to write response: " + err.Error(),
			})
		return
	}
	logrus.Error(err)
}

func WriteError(
	request *restful.Request,
	response *restful.Response,
	httpStatusCode int,
	eventErr error,
	errorResponse *Error,
) {
	resp, err := json.Marshal(errorResponse)
	if err != nil {
		logrus.WithFields(logrus.Fields{"status_code": httpStatusCode}).
			Error("unable to marshal error response: ", err)
	}
	err = response.WriteErrorString(httpStatusCode, string(resp))
	if err != nil {
		logrus.WithFields(logrus.Fields{"status_code": httpStatusCode, "error_response": errorResponse}).
			Error("unable to write error response: ", err)
	}
	if httpStatusCode >= 500 {
		logrus.Error(request, errorResponse.ErrorCode, errorResponse.ErrorMessage+":", eventErr)
		return
	}
	logrus.Error(request, errorResponse.ErrorCode, errorResponse.ErrorMessage+":", eventErr)
}
