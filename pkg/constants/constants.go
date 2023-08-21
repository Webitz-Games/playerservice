package constants

const (
	Separator             = "_"
	ActionCreatePlayer    = "create_player"
	UnableToWriteResponse = 99

	InternalServerError         = 20000
	UnauthorizedAccess          = 20001
	ValidationError             = 20002
	ForbiddenAccess             = 20003
	TooManyRequests             = 20007
	UserNotFound                = 20008
	InsufficientPermissions     = 20013
	InvalidAudience             = 20014
	InsufficientScope           = 20015
	UnableToParseRequestBody    = 20019
	InvalidPaginationParameters = 20021
	TokenIsNotUserToken         = 20022
)

var ErrorCodeMapping = map[int]string{
	// Global Error Codes
	InternalServerError:         "internal server error",
	UnauthorizedAccess:          "unauthorized access",
	ValidationError:             "validation error",
	ForbiddenAccess:             "forbidden access",
	TooManyRequests:             "too many requests",
	UserNotFound:                "user not found",
	InsufficientPermissions:     "insufficient permissions",
	InvalidAudience:             "invalid audience",
	InsufficientScope:           "insufficient scope",
	UnableToParseRequestBody:    "unable to parse request body",
	InvalidPaginationParameters: "invalid pagination parameter",
	TokenIsNotUserToken:         "token is not user token",
}
