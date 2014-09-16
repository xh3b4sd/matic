package v1

import (
	srvPkg "github.com/catalyst-zero/middleware-server"
	logPkg "github.com/op/go-logging"
)

type V1 struct {
	Logger *logPkg.Logger
}

func (v1 *V1) SetupHelloRoute(srv *srvPkg.Server) {
	srv.Serve("GET", "/v1/hello", v1.Hello)
}

func SetupWorldRoute(v1 *V1, srv *srvPkg.Server) {
	srv.Serve("GET", "/v1/world", Foo, v1.World)
}
