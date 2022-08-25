package main

import (
	"github.com/raj-ptl/go-status-check/constants"
	"github.com/raj-ptl/go-status-check/server"
	"github.com/raj-ptl/go-status-check/status"
)

func main() {
	go status.PollUpdateAllSites(constants.POLLING_RATE)
	server.ServeRequests()
}
