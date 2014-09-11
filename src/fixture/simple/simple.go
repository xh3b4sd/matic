package main

import (
	srvPkg "github.com/catalyst-zero/middleware-server"
)

func main() {
	logger := srvPkg.NewLogger(srvPkg.LoggerOptions{Name: "middleware-example"})
	v1 := &V1{Logger: logger}

	srv := srvPkg.NewServer("127.0.0.1", "8080")
	srv.SetLogger(logger)

	srv.Serve("GET", "/v1/hello-world", v1.middlewareOne, v1.middlewareTwo)

	srv.Listen()
}
