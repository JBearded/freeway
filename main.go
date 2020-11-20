package main

import (
	"fmt"
	"freeway/common"
	"freeway/config"
	"freeway/gossip"
	"freeway/logger"
	"freeway/server"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"github.com/hashicorp/memberlist"
	log "github.com/sirupsen/logrus"
)

var gossipMemberlist memberlist.Memberlist

func main() {
	profileName := os.Getenv("PROFILE")
	if strings.Trim(profileName, " ") == "" {
		log.Fatal("environment PROFILE empty")
	}
	profile := common.ParseProfile(profileName)
	if profile == common.NoneProfile {
		log.Fatal("incorrect environment PROFILE, [develop,test,production]")
	}
	err := config.Init(profile)
	if err != nil {
		log.Fatal("init profile resource error ", err)
	}
	logger.Init(profile)
	go server.Start(config.Get().Server.Websocket.Port, func(w http.ResponseWriter, r *http.Request) {
		server.RunWebsocketServer(w, r)
	})
	go server.Start(config.Get().Server.HTTP.Port, func(w http.ResponseWriter, r *http.Request) {
		server.RunHTTPServer(w, r)
	})

	myDelegate := new(gossip.MyDelegate)
	uuidValue, _ := uuid.NewUUID()
	gossipConfig := memberlist.DefaultLocalConfig()
	gossipConfig.Name = uuidValue.String()
	gossipConfig.BindPort = 7788
	gossipConfig.AdvertisePort = gossipConfig.BindPort
	gossipConfig.Delegate = myDelegate

	gossipMemberlist, err := memberlist.Create(gossipConfig)
	if err != nil {
		log.Error("Failed to create memberlist: " + err.Error())
	}
	localNode := gossipMemberlist.LocalNode()

	// Join an existing cluster by specifying at least one known member.
	_, joinErr := gossipMemberlist.Join([]string{
		fmt.Sprintf("%s:%d", localNode.Addr.To4(), localNode.Port),
	})
	if joinErr != nil {
		log.Error("Failed to join cluster: " + err.Error())
	}

	ChanShutdown := make(chan os.Signal)
	signal.Ignore(syscall.SIGHUP)
	<-ChanShutdown
}
