package server

import (
	"net/http"
	"time"

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
	server.ListenAndServe()
}

//RunHTTPServer 运行http服务
func RunHTTPServer(w http.ResponseWriter, r *http.Request) {
	app := globalWsManager.apps["123"]
	for key := range app.conns {
		client := app.conns[key]
		log.Info(client.connID)
	}
	appID := r.URL.Query().Get("appId")
	connID := r.URL.Query().Get("connId")
	log.Info("appId:" + appID + " connId:" + connID)
	w.Write([]byte(appID))
}
