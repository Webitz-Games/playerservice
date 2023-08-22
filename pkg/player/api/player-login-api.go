package api

import "github.com/emicklei/go-restful/v3"

func addPlayerLoginRoute(webservice *restful.WebService, handler LoginPlayerHandler) {

}

func bindLoginPlayerHandler(handler LoginPlayerHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {

	}
}
