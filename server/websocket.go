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

const pingErrorLimitTimes int32 = 2

type wsClient struct {
	appID          string
	connID         string
	addr           string
	conn           *websocket.Conn
	closed         atomic.Value
	pingErrorTimes int32
}

func (client *wsClient) String() string {
	return client.appID + " " + client.connID + " " + client.addr + " " + strconv.FormatBool(client.closed.Load().(bool))
}

func (client *wsClient) pingLoop() {
	ticker := time.NewTicker(config.Get().Server.Websocket.PingPeriodSeconds * time.Second)
	//释放链接和数据
	defer func() {
		ticker.Stop()
		client.closed.Store(true)
		delete(globalWsManager.apps[client.appID].conns, client.connID)
		client.conn.Close()
	}()
	for {
		if client.closed.Load().(bool) {
			return
		}
		<-ticker.C
		if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			atomic.AddInt32(&client.pingErrorTimes, 1)
			//连续n次失败之后，直接return
			if client.pingErrorTimes >= pingErrorLimitTimes {
				return
			}
		} else {
			atomic.StoreInt32(&client.pingErrorTimes, 0)
		}
	}
}

func (client *wsClient) receive() {
	conn := client.conn
	//释放链接和数据
	defer func() {
		client.closed.Store(true)
		delete(globalWsManager.apps[client.appID].conns, client.connID)
		conn.Close()
	}()
	for {
		messageType, p, err := conn.ReadMessage()
		if messageType == websocket.CloseMessage {
			log.Info("closeMessageType, ", client.String())
			return
		}
		if err != nil {
			log.Error("readMessage error, ", client.String())
			return
		}
		go client.handleMessage(messageType, p)
	}
}

func (client *wsClient) handleMessage(messageType int, message []byte) {
	//TODO 处理消息
	client.conn.WriteMessage(messageType, message)
}

func (client *wsClient) push(message string) {
	client.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

type wsApp struct {
	appID         string
	conns         map[string]*wsClient
	brocastBuffer chan []byte
}

type wsManager struct {
	apps     map[string]*wsApp
	appLocks sync.Map
}

var globalWsManager = &wsManager{
	apps: make(map[string]*wsApp),
}

// GetGlobalWsManager 获取全局webscoket管理变量
func GetGlobalWsManager() *wsManager {
	return globalWsManager
}

func (manager *wsManager) brocast(message string) {
	//TODO
}

func (manager *wsManager) add(appID string, con *websocket.Conn, remoteAddr string) {
	uuid, err := uuid.NewUUID()
	log.Info(uuid)
	if err != nil {
		log.Error("uuid error", err)
		return
	}
	newConnID := uuid.String()
	client := &wsClient{
		appID:  appID,
		connID: newConnID,
		conn:   con,
		addr:   remoteAddr,
	}
	client.closed.Store(false)
	atomic.AddInt64(&Count, 1)
	log.Info(client.String() + " count : " + strconv.FormatInt(Count, 10))
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

	go client.pingLoop()
	go client.receive()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  config.Get().Server.Websocket.ReadBufferSize,
	WriteBufferSize: config.Get().Server.Websocket.WriteBufferSize,
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
	addr := r.RemoteAddr
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("websocket upgrade error", err)
	}
	appID := r.URL.Query().Get("appId")
	if len(appID) <= 0 {
		log.Error("empty appId")
		return
	}
	globalWsManager.add(appID, conn, addr)
}
