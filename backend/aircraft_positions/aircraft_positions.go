package main

import (
	gcp "aircraft_positions/gcp"
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

const env_credJson = "GOOGLE_APPLICATION_CREDENTIALS"

const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"
const env_topicID = "GOOGLE_PUBSUB_AIRCRAFT_POSITIONS_TOPIC_ID"
const env_subID = "GOOGLE_PUBSUB_AIRCRAFT_POSITIONS_SUBSCRIBER_ID"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

var posTopic *pubsub.Topic

type WSClient struct {
	icao24       string
	ch           chan string
	chErr        chan string
	clientCtx    context.Context
	clientCancel func()
	ws           *websocket.Conn
}

func (cl *WSClient) Send(msg string) {
	cl.ch <- msg
}

func (cl *WSClient) SendErrAndClose(msg string) {
	select {
	case cl.chErr <- msg:
	default:
		fmt.Printf("Failed sending error on wsClient")
		cl.clientCancel()
	}
}

func (cl *WSClient) WSWriteLoop() {
	defer cl.clientCancel()
	defer cl.ws.Close()

	var msg string
	var err error

	for {
		select {
		case msg = <-cl.ch:
			err = cl.ws.WriteMessage(websocket.TextMessage, []byte(msg))

			if err != nil {
				fmt.Println(err)
				cl.clientCancel()
				cl.ws.Close()
				hub.UnregisterClient(cl)
				return
			}

		case msg = <-cl.chErr:
			cl.ws.WriteMessage(websocket.CloseMessage, []byte(msg))
			cl.clientCancel()
			cl.ws.Close()
			return

		case <-cl.clientCtx.Done():
			cl.ws.WriteMessage(websocket.CloseMessage, []byte("Closing the websocket"))
			cl.ws.Close()
			return
		}
	}
}

func (cl *WSClient) WSReadLoop() {
	defer cl.clientCancel()

	for {
		select {
		case <-cl.clientCtx.Done():
			return
		default:
		}

		msgType, _, err := cl.ws.NextReader()

		if err != nil {
			cl.clientCancel()
			cl.ws.Close()
			hub.UnregisterClient(cl)
			return
		}

		if msgType == websocket.CloseMessage {
			cl.clientCancel()
			cl.ws.Close()
			hub.UnregisterClient(cl)
			return
		}
	}
}

type RefreshableTimer struct {
	duration       time.Duration
	refresh        chan int
	timeoutHandler func()
	timerCtx       context.Context
	timerCancel    func()
}

func (state *RefreshableTimer) StartTimer() {
	go func() {
		defer state.timerCancel()

		for {
			select {
			case <-state.refresh:
			case <-state.timerCtx.Done():
				return
			case <-time.After(state.duration):
				state.timeoutHandler()
				return
			}
		}
	}()
}

func (state *RefreshableTimer) StopTimer() {
	state.timerCancel()
}

func (state *RefreshableTimer) RefreshTimer() {
	trySend(state.refresh)
}

func (state *RefreshableTimer) SetTimeoutHandler(f func()) {
	state.timeoutHandler = f
}

type PubSubListener struct {
	icao24         string
	subId          string
	cache          string
	cacheLock      sync.RWMutex
	wsClients      []*WSClient
	clientsLock    sync.RWMutex
	listenerCtx    context.Context
	listenerCancel func()
	timer          RefreshableTimer
}

func (state *PubSubListener) StartPubSubListener(ctx context.Context) {
	defer state.listenerCancel()

	filter := `"attributes.icao24 = "` + state.icao24 + `"`
	sub, err := gcp.CreateSubscription(state.subId+"_"+state.icao24, posTopic, filter)
	checkErr(err)

	state.timer.StartTimer()

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		defer msg.Ack()

		state.timer.RefreshTimer()

		strMsg := string(msg.Data)

		state.SendAll(strMsg)

		state.SetCache(strMsg)
	})

	if err != nil {
		fmt.Println("Error in Pub/Sub Receive: ", err)
		state.timer.StopTimer()
		hub.UnregisterFailedListener(state)
	}
}

func (state *PubSubListener) AddWSClient(cli *WSClient) {
	state.clientsLock.Lock()
	defer state.clientsLock.Unlock()
	state.wsClients = append(state.wsClients, cli)
}

func (state *PubSubListener) RemoveWSClient(cli *WSClient) {
	state.clientsLock.Lock()
	defer state.clientsLock.Unlock()

	idx := slices.IndexFunc(state.wsClients, func(c *WSClient) bool { return c == cli })

	if idx >= 0 {
		length := len(state.wsClients) - 1
		state.wsClients[idx] = state.wsClients[length]
		state.wsClients = state.wsClients[:length]
	} else {
		fmt.Println("No client found in pbListener WSClients list")
	}
}

func (state *PubSubListener) SendAll(msg string) {
	state.clientsLock.RLock()
	defer state.clientsLock.RUnlock()
	for i := 0; i < len(state.wsClients); i++ {
		state.wsClients[i].Send(msg)
	}
}

func (state *PubSubListener) SendAllErr(msg string) {
	state.clientsLock.RLock()
	defer state.clientsLock.RUnlock()
	for i := 0; i < len(state.wsClients); i++ {
		state.wsClients[i].SendErrAndClose(msg)
	}
}

func (state *PubSubListener) GetCache() string {
	state.cacheLock.RLock()
	defer state.cacheLock.RUnlock()
	return state.cache
}

func (state *PubSubListener) SetCache(val string) {
	state.cacheLock.Lock()
	defer state.cacheLock.Unlock()
	state.cache = val
}

func NewPubSubListener(icao24 string, subID string, ctx context.Context, stopListen func()) *PubSubListener {
	timeoutCtx, timeoutCancel := context.WithCancel(context.Background())

	ls := PubSubListener{
		icao24:         icao24,
		subId:          subID,
		cache:          "",
		listenerCtx:    ctx,
		listenerCancel: stopListen,
		timer: RefreshableTimer{
			duration:    time.Second * 30,
			refresh:     make(chan int, 5),
			timerCtx:    timeoutCtx,
			timerCancel: timeoutCancel,
		},
	}

	return &ls
}

type Hub struct {
	sync.RWMutex
	icaos map[string]*PubSubListener
	SubID string
}

func (h *Hub) RegisterClient(cl *WSClient) {
	h.Lock()
	defer h.Unlock()

	if existingListener, ok := h.icaos[cl.icao24]; ok {
		existingListener.AddWSClient(cl)

		firstVal := existingListener.GetCache()
		if firstVal != "" {
			cl.Send(firstVal)
		}
	} else {
		cancellableCtx, listenerCancel := context.WithCancel(context.Background())

		newListener := NewPubSubListener(cl.icao24, h.SubID, cancellableCtx, listenerCancel)
		newListener.timer.SetTimeoutHandler(h.MakeListenerTimeoutHandler(newListener))

		h.icaos[cl.icao24] = newListener

		h.icaos[cl.icao24].AddWSClient(cl)

		go newListener.StartPubSubListener(cancellableCtx)
	}
}

func (h *Hub) UnregisterClient(cl *WSClient) {
	h.Lock()
	defer h.Unlock()

	if h.icaos[cl.icao24] == nil {
		return
	}

	psListener := h.icaos[cl.icao24]
	psListener.RemoveWSClient(cl)

	length := len(psListener.wsClients)
	if length <= 0 {
		psListener.listenerCancel()
		psListener.timer.StopTimer()

		delete(h.icaos, cl.icao24)
	}
}

func (h *Hub) UnregisterFailedListener(ls *PubSubListener) {
	h.Lock()
	defer h.Unlock()

	fmt.Println("UnregisterFailedListener: ", ls.icao24)

	// Notify all clients of the error and stop the writer/listeners goroutines
	ls.SendAllErr("Error: server was unable to receive positions")

	// Delete the PubSubListener from the map of active listeners
	delete(h.icaos, ls.icao24)
}

func (h *Hub) MakeListenerTimeoutHandler(ls *PubSubListener) func() {
	return func() {
		ls.listenerCancel()
		h.UnregisterFailedListener(ls)
	}
}

var hub = Hub{}

var upgrader = websocket.Upgrader{}

func main() {

	var credJson = mustGetenv(env_credJson)
	var projectID = mustGetenv(env_projectID)
	var topicID = mustGetenv(env_topicID)
	var subID = mustGetenv(env_subID)

	err := gcp.Initialize(credJson, projectID)
	checkErr(err)
	defer gcp.PubsubClient.Close()

	posTopic, err = gcp.GetTopic(topicID)
	checkErr(err)

	hub.icaos = make(map[string]*PubSubListener)
	hub.SubID = subID

	router := gin.New()

	router.SetTrustedProxies(nil)
	router.GET("/airspace/aircraft/:icao24/position", httpRequestHandler)

	router.Run()
}

func httpRequestHandler(c *gin.Context) {
	icao24 := c.Param("icao24")

	if !validIcao24(icao24) {
		c.String(http.StatusNotAcceptable, "invalid icao24")
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// TODO check what best to do here
		fmt.Println(err)
		return
	}

	cancellableCtx, cancel := context.WithCancel(context.Background())

	cl := &WSClient{
		ws:           ws,
		icao24:       icao24,
		ch:           make(chan string, 100),
		chErr:        make(chan string, 10),
		clientCtx:    cancellableCtx,
		clientCancel: cancel,
	}

	hub.RegisterClient(cl)

	go cl.WSReadLoop()
	go cl.WSReadLoop()
}

func validIcao24(icao24 string) bool {
	return len(icao24) == 6
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("Environment variable not set: " + k)
	}
	return v
}

func trySend(ch chan int) {
	select {
	case ch <- 1:
	default:
	}
}
