package main

import (
	srvPkg "github.com/catalyst-zero/middleware-server"
	v1Pkg "github.com/zyndiecate/matic/fixture/simple/middleware/v1"
)

func main() {
	logger := srvPkg.NewLogger(srvPkg.LoggerOptions{Name: "middleware-example"})
	v1 := &v1Pkg.V1{Logger: logger}

	srv := srvPkg.NewServer("127.0.0.1", "8080")
	srv.SetLogger(logger)

	srv.Serve("GET", "/v1/hello-world", v1.MiddlewareOne, v1.MiddlewareTwo)

	srv.Listen()
}
