package v1

import (
	"fmt"
	"net/http"

	apiSchemaPkg "github.com/catalyst-zero/api-schema"
	srvPkg "github.com/catalyst-zero/middleware-server"
)

func (v1 *V1) HelloWorldOne(res http.ResponseWriter, req *http.Request, ctx *srvPkg.Context) error {
	v1.Logger.Debug("hello world one")
	return ctx.Next()
}

func (v1 *V1) HelloWorldTwo(res http.ResponseWriter, req *http.Request, ctx *srvPkg.Context) error {
	v1.Logger.Debug("hello world two")
	return ctx.Next()
}

func Foo(res http.ResponseWriter, req *http.Request, ctx *srvPkg.Context) error {
	fmt.Printf("%#v\n", "foo")
	return ctx.Next()
}

func (v1 *V1) HelloWorldThree(res http.ResponseWriter, req *http.Request, ctx *srvPkg.Context) error {
	v1.Logger.Debug("hello world three")
	return ctx.Response.Json(apiSchemaPkg.StatusData("hello world"), http.StatusOK)
}

func (v1 *V1) Hello(res http.ResponseWriter, req *http.Request, ctx *srvPkg.Context) error {
	v1.Logger.Debug("hello")
	return ctx.Response.Json(apiSchemaPkg.StatusData("hello"), http.StatusOK)
}

func (v1 *V1) World(res http.ResponseWriter, req *http.Request, ctx *srvPkg.Context) error {
	v1.Logger.Debug("world")
	return ctx.Response.Json(apiSchemaPkg.StatusData("world"), http.StatusOK)
}
