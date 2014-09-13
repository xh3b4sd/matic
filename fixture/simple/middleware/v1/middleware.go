package v1

import (
	"net/http"

	apiSchemaPkg "github.com/catalyst-zero/api-schema"
	srvPkg "github.com/catalyst-zero/middleware-server"
)

func (v1 *V1) MiddlewareOne(res http.ResponseWriter, req *http.Request, ctx *srvPkg.Context) error {
	v1.Logger.Debug("middleware one")
	return ctx.Next()
}

func (v1 *V1) MiddlewareTwo(res http.ResponseWriter, req *http.Request, ctx *srvPkg.Context) error {
	v1.Logger.Debug("middleware two")
	return ctx.Response.Json(apiSchemaPkg.StatusData("hello world"), http.StatusOK)
}
