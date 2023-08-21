package swagger

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"net/http"
)

func CreateSwagger(basePath string, serviceName string, version string) *restful.WebService {
	swaggerConfig := restfulspec.Config{
		WebServices: restful.DefaultContainer.RegisteredWebServices(),
		APIPath:     "/apidocs/api.json",
		PostBuildSwaggerObjectHandler: func(s *spec.Swagger) {
			s.Info = &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       serviceName,
					Description: "Webitz Player Service",
					Version:     version,
				},
			}
		},
	}

	http.Handle(basePath+"/apidocs/",
		http.StripPrefix(basePath+"/apidocs/", http.FileServer(http.Dir("player/docs/swagger.json"))))

	return restfulspec.NewOpenAPIService(swaggerConfig)
}
