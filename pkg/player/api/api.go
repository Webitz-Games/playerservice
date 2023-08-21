package api

import (
	"github.com/emicklei/go-restful/v3"
	"playerapi/pkg/swagger"
)

type RequestHandlers interface {
	PlayerRequestHandlers
}

func RegisterRoutes(basePath string, serviceName string, handlers RequestHandlers) {
	restful.DefaultContainer.Add(AddRoutes(serviceName, handlers))
	restful.DefaultContainer.Add(swagger.CreateSwagger(basePath, serviceName, "0.0.1"))
}

func AddRoutes(serviceName string, handlers RequestHandlers) *restful.WebService {
	webservice := new(restful.WebService)
	webservice.
		Path("webitz/player/api").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc(serviceName).
		ApiVersion("0.0.1")

	addPlayerRoutes(webservice, handlers)

	return webservice
}
