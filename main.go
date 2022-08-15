package main

import (
	"github.com/raj-ptl/go-status-check/server"
	"github.com/raj-ptl/go-status-check/status"
)

func main() {
	status.InitializeMap()
	server.ServeRequests()
}
