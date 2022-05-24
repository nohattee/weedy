package main

import (
	http "weedy/internal/controller/http/v1"
	"weedy/pkg/httpserver"
)

func main() {
	server := httpserver.NewHttpServer("", "8080", "")

	server.AddController(&http.UserController{})
	server.Run()
}
