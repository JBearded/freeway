package server

import (
	"encoding/json"
	"freeway/model"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type router struct {
	method  string
	handler func(http.ResponseWriter, *http.Request)
}

// Routers http接口配置
var Routers map[string]*router = map[string]*router{
	"/websocket/info": {
		method:  http.MethodGet,
		handler: websocketInfo,
	},
	"/websocket/push": {
		method:  http.MethodGet,
		handler: websocketPush,
	},
	"/api": {
		method:  http.MethodGet,
		handler: api,
	},
}

func websocketInfo(w http.ResponseWriter, r *http.Request) {
	appID := r.URL.Query().Get("appId")
	manager := GetGlobalWsManager()
	result := make(map[string]string)
	for connID := range manager.apps[appID].conns {
		client := manager.apps[appID].conns[connID]
		result[connID] = client.String()
	}
	data, _ := json.Marshal(result)
	w.Write(data)
}

func websocketPush(w http.ResponseWriter, r *http.Request) {
	appID := r.URL.Query().Get("appId")
	connID := r.URL.Query().Get("connId")
	message := r.URL.Query().Get("message")
	manager := GetGlobalWsManager()
	client := manager.apps[appID].conns[connID]
	if client == nil {
		w.Write([]byte("nil"))
	} else {
		client.push(message)
		w.Write([]byte("ok"))
	}
}

func api(w http.ResponseWriter, r *http.Request) {
	body, readBodyErr := ioutil.ReadAll(r.Body)
	if readBodyErr != nil {
		log.Error("read request body error", readBodyErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var ob model.Student
	json.Unmarshal(body, &ob)
	ob.ID = 2
	result, _ := json.Marshal(ob)
	w.Write(result)
}
