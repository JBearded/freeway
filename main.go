package main

import (
	"freeway/common"
	"freeway/config"
	"freeway/logger"
	"freeway/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	profileName := os.Getenv("PROFILE")
	profile := common.ParseProfile(profileName)
	err := config.Init(profile)
	if err != nil {
		log.Error(err)
		return
	}
	logger.Init(profile)
	go server.Start(config.Get().Server.Websocket.Port, func(w http.ResponseWriter, r *http.Request) {
		server.RunWebsocketServer(w, r)
	})
	go server.Start(config.Get().Server.HTTP.Port, func(w http.ResponseWriter, r *http.Request) {
		server.RunHTTPServer(w, r)
	})

	ChanShutdown := make(chan os.Signal)
	signal.Ignore(syscall.SIGHUP)
	<-ChanShutdown
}
