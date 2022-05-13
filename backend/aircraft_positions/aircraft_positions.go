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
	icao24     string
	ch         chan string
	stopReader chan int
	stopWriter chan int
	ws         *websocket.Conn
}

func (cl *WSClient) Send(msg string) {
	cl.ch <- msg
}

func (cl *WSClient) WSReadLoop() {
	for {
		if _, _, err := cl.ws.NextReader(); err != nil {
			cl.ws.Close()
			hub.UnregisterClient(cl)
			break
		}
	}
}

// func (cl *WSClient) WSWriteLoop() {
// 	for {
// 		if _, err := cl.ws.NextWriter(1); err != nil {
// 			cl.ws.Close()
// 			hub.UnregisterClient(cl)
// 			break
// 		}
// 	}
// }

type RefreshableTimer struct {
	duration       time.Duration
	stop           chan int
	refresh        chan int
	timeoutHandler func()
}

func (state *RefreshableTimer) StartTimer() {
	go func() {
		for {
			select {
			case <-state.refresh:
			case <-state.stop:
				return
			case <-time.After(state.duration):
				state.timeoutHandler()
				return
			}
		}
	}()
}

func (state *RefreshableTimer) StopTimer() {
	trySend(state.stop)
}

func (state *RefreshableTimer) RefreshTimer() {
	trySend(state.refresh)
}

type PubSubListener struct {
	sync.RWMutex
	Icao24     string
	SubID      string
	cache      string
	wsClients  []*WSClient
	stopListen func()
	timer      RefreshableTimer
}

func makePubSubListener(icao24 string, subID string, stopListen func()) *PubSubListener {
	ps := &PubSubListener{
		Icao24:     icao24,
		SubID:      subID,
		cache:      "",
		stopListen: stopListen,
		timer: RefreshableTimer{
			duration:       time.Second * 30,
			stop:           make(chan int, 5),
			refresh:        make(chan int, 5),
			timeoutHandler: stopListen,
		},
	}
	return ps
}

func (state *PubSubListener) AddWSClient(cli *WSClient) {
	state.Lock()
	defer state.Unlock()
	state.wsClients = append(state.wsClients, cli)
}

func (state *PubSubListener) RemoveWSClient(cli *WSClient) {
	state.Lock()
	defer state.Unlock()
	idx := slices.IndexFunc(state.wsClients, func(c *WSClient) bool { return *c == *cli })
	if idx >= 0 {
		length := len(state.wsClients) - 1
		state.wsClients[idx] = state.wsClients[length]
		state.wsClients = state.wsClients[:length]
	} else {
		fmt.Println("No client found in pbListener WSClients list")
	}
}

func (state *PubSubListener) SendAll(msg string) {
	state.RLock()
	defer state.RUnlock()
	for i := 0; i < len(state.wsClients); i++ {
		state.wsClients[i].Send(msg)
	}
}

func (state *PubSubListener) GetCache() string {
	state.RLock()
	defer state.RUnlock()
	return state.cache
}

func (state *PubSubListener) SetCache(val string) {
	state.Lock()
	defer state.Unlock()
	state.cache = val
}

func (state *PubSubListener) StartPubSubListener(ctx context.Context) {
	filter := `"attributes.icao24 = "` + state.Icao24 + `"`
	sub, err := gcp.CreateSubscription(state.SubID+"_"+state.Icao24, posTopic, filter)
	checkErr(err)

	state.timer.StartTimer()

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		state.timer.RefreshTimer()

		strMsg := string(msg.Data)

		state.SendAll(strMsg)

		state.SetCache(strMsg)

		msg.Ack()
	})

	if err != nil {
		fmt.Println("Error in Pub/Sub Receive: ", err)
		state.timer.StopTimer()
	}
}

func (state *PubSubListener) StopPubSubListener() {
	state.timer.StopTimer()
	state.stopListen()
}

type Hub struct {
	sync.RWMutex
	icaos map[string]*PubSubListener
	SubID string
}

func (h *Hub) RegisterClient(cl *WSClient) {
	h.Lock()
	defer h.Unlock()

	if h.icaos[cl.icao24] == nil {
		cancellableCtx, cancel := context.WithCancel(context.Background())

		psListener := makePubSubListener(cl.icao24, h.SubID, cancel)

		h.icaos[cl.icao24] = psListener

		h.icaos[cl.icao24].AddWSClient(cl)

		go psListener.StartPubSubListener(cancellableCtx)

	} else {
		h.icaos[cl.icao24].AddWSClient(cl)

		firstVal := h.icaos[cl.icao24].GetCache()
		if firstVal != "" {
			cl.Send(firstVal)
		}
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
		psListener.StopPubSubListener()

		h.icaos[cl.icao24] = nil
	}
}

var hub = Hub{}

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

	defer ws.Close()

	cl := &WSClient{
		ws:         ws,
		icao24:     icao24,
		ch:         make(chan string, 100),
		stopReader: make(chan int, 5),
		stopWriter: make(chan int, 5),
	}

	hub.RegisterClient(cl)

	go cl.WSReadLoop()

	var msg string

	for {
		msg = <-cl.ch
		err = ws.WriteMessage(websocket.TextMessage, []byte(msg))

		if err != nil {
			fmt.Println(err)
			hub.UnregisterClient(cl)
			break
		}
	}

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

func tryRecv(ch chan int) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func pubsubHandler(topicID string, subscriptionID string) {

	topic, err := gcp.GetTopic(topicID)
	checkErr(err)

	sub, err := gcp.GetSubscription(subscriptionID, topic)
	checkErr(err)

	err = sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		messageHandler(subscriptionID, ctx, msg)
	})
	checkErr(err)
}

func messageHandler(subscriptionID string, ctx context.Context, msg *pubsub.Message) {
	defer func() {
		if r := recover(); r != nil {
			msg.Ack()
		}
	}()

	msg.Ack()

	// aircraftList.Write(string(msg.Data))
}

var upgrader = websocket.Upgrader{}
