package server

import (
	"freeway/config"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type wsClient struct {
	appID      string
	connID     string
	conn       *websocket.Conn
	readBuffer chan []byte
	closed     atomic.Value
}

func (client *wsClient) receive() {
	conn := client.conn
	//TODO 释放所有相关数据
	defer func() {
		client.closed.Store(true)
		delete(globalWsManager.apps[client.appID].conns, client.connID)
		close(client.readBuffer)
		conn.Close()
	}()
	for {
		messageType, p, err := conn.ReadMessage()
		if messageType == websocket.CloseMessage {
			log.Info("close connection", client.connID)
			return
		}
		if err != nil {
			log.Error("readMessage error", err)
			return
		}
		client.readBuffer <- p
	}
}
func (client *wsClient) read() {
	ticker := time.NewTicker(60 * time.Second)
	//TODO 释放所有相关数据
	defer func() {
		client.conn.Close()
		ticker.Stop()
	}()
	for {
		log.Info("read...")
		select {
		case message, ok := <-client.readBuffer:
			if ok {
				client.conn.WriteMessage(websocket.TextMessage, message)
			}
		case <-ticker.C:
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error("ping error", err)
				return
			}
		}
	}
}

type wsApp struct {
	appID         string
	conns         map[string]*wsClient
	brocastBuffer chan []byte
}

var globalWsManager = &wsManager{
	apps: make(map[string]*wsApp),
}

type wsManager struct {
	apps     map[string]*wsApp
	appLocks sync.Map
}

func (manager *wsManager) add(appID string, con *websocket.Conn) {
	uuid, err := uuid.NewUUID()
	log.Info(uuid)
	if err != nil {
		log.Error("uuid error", err)
		return
	}
	newConnID := uuid.String()
	client := &wsClient{
		appID:      appID,
		connID:     newConnID,
		conn:       con,
		readBuffer: make(chan []byte, 1024),
	}
	sm := &sync.Mutex{}
	value, _ := manager.appLocks.LoadOrStore(appID, sm)
	appLock := value.(*sync.Mutex)
	appLock.Lock()
	if manager.apps[appID] == nil {
		newWsApp := &wsApp{
			appID: appID,
			conns: make(map[string]*wsClient),
		}
		manager.apps[appID] = newWsApp
	}
	manager.apps[appID].conns[newConnID] = client
	appLock.Unlock()

	go client.receive()
	go client.read()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		match, err := regexp.Match(config.Get().Server.Websocket.AllowOrigin, []byte(origin))
		log.Info(origin, match, err)
		return (err == nil && match)
	},
}

// Count 计数器
var Count int64 = 0

//RunWebsocketServer 运行websocket
func RunWebsocketServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("websocket upgrade error", err)
	}
	atomic.AddInt64(&Count, 1)
	appID := r.URL.Query().Get("appId")
	log.Info(appID + " count : " + strconv.FormatInt(Count, 10))
	if len(appID) <= 0 {
		log.Error("empty appId")
		return
	}
	globalWsManager.add(appID, conn)
}
