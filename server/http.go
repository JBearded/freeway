package server

import (
	"freeway/common"
	"freeway/config"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/petermattis/goid"
	log "github.com/sirupsen/logrus"
)

type myHTTPHandler struct {
	handlerMethod func(http.ResponseWriter, *http.Request)
}

func (h *myHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handlerMethod(w, r)
}

// Start 启动服务
func Start(port string, handlerMethod func(http.ResponseWriter, *http.Request)) {
	httpHandler := &myHTTPHandler{
		handlerMethod: handlerMethod,
	}
	server := &http.Server{
		Addr:           ":" + port,
		Handler:        httpHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	defer server.Close()
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("start server error, port:"+port+" ", err)
	}
}

//RunHTTPServer 运行http服务
func RunHTTPServer(w http.ResponseWriter, r *http.Request) {
	goID := goid.Get()
	log.Info("http request go id:", goID)
	url := r.URL.Path
	pathReg, _ := regexp.Compile("(?:\\/([^?#]*))?")
	path := string(pathReg.Find([]byte(url)))
	values := strings.Split(path, ".")
	var api, suffix string
	if len(values) > 0 {
		api = values[0]
	}
	if len(values) > 1 {
		suffix = values[1]
	}
	log.Info("url:", url, " ", "request:", api, " ", "suffix:", suffix)
	rt := Routers[api]
	if rt == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not Found"))
		return
	}
	requestMethod := r.Method
	if requestMethod != rt.method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Request Method Not Allowed"))
		return
	}
	if contentType := common.GetContentType(suffix); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	//跨域支持
	crossReg, _ := regexp.Compile(config.Get().Server.HTTP.AllowOrigin)
	origin := r.Header.Get("Origin")
	referer := r.Header.Get("Referer")
	var crossURL string = referer
	if origin != "" {
		crossURL = origin
	}
	if crossReg.Match([]byte(crossURL)) {
		w.Header().Add("Access-Control-Allow-Headers", "X-Requested-With")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST")
		w.Header().Add("Access-Control-Allow-Origin", crossURL)
	}
	rt.handler(w, r)
}
