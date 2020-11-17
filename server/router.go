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
	"/hello": {
		method:  http.MethodGet,
		handler: hello,
	},
	"/api": {
		method:  http.MethodPost,
		handler: api,
	},
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
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
