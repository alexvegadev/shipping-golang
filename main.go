package main

import (
	"shipping/service"
)

func main() {
	server := new(service.Server)
	server.RunServer()
}
