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

type WSClient chan string

func (c *WSClient) Send(msg string) {
	*c <- msg
}

type PubSubListener struct {
	sync.RWMutex
	Icao24    string
	SubID     string
	wsClients []WSClient
	cache     string
	//ctx       context.Context
	Cancel    func()
	Timeout   time.Duration
	stopCh    chan int
	refreshCh chan int
}

func (state *PubSubListener) AddWSClient(cli WSClient) {
	state.Lock()
	defer state.Unlock()
	state.wsClients = append(state.wsClients, cli)
}

func (state *PubSubListener) RemoveWSClient(cli WSClient) {
	state.Lock()
	defer state.Unlock()
	idx := slices.IndexFunc(state.wsClients, func(c WSClient) bool { return c == cli })
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

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		tryNotifyCh(state.refreshCh)

		strMsg := string(msg.Data)

		state.SendAll(strMsg)

		state.SetCache(strMsg)

		msg.Ack()
	})

	if err != nil {
		fmt.Println("Error in Pub/Sub Receive: ", err)
		state.stopCh <- 1
	}
}

type Hub struct {
	sync.RWMutex
	icaos map[string]*PubSubListener
	SubID string
}

func refreshableTimeout(refresh chan int, stop chan int, timeoutHandler func(), timeout time.Duration) {
	for {
		select {
		case <-refresh:
		case <-stop:
			return
		case <-time.After(timeout):
			timeoutHandler()
			return
		}
	}
}

func (h *Hub) RegisterClient(icao24 string) chan string {
	h.Lock()
	defer h.Unlock()

	ch := make(chan string, 100)

	if h.icaos[icao24] == nil {
		psListener := &PubSubListener{}
		psListener.Icao24 = icao24
		psListener.SubID = h.SubID
		psListener.Timeout = time.Second * 30
		psListener.stopCh = make(chan int, 5)
		psListener.refreshCh = make(chan int, 5)
		psListener.cache = ""

		var ctx context.Context
		ctx, psListener.Cancel = context.WithCancel(context.Background())

		h.icaos[icao24] = psListener

		h.icaos[icao24].AddWSClient(ch)

		go psListener.StartPubSubListener(ctx)

		go refreshableTimeout(psListener.refreshCh, psListener.stopCh, psListener.Cancel, psListener.Timeout)

	} else {
		h.icaos[icao24].AddWSClient(ch)

		firstVal := h.icaos[icao24].GetCache()
		if firstVal != "" {
			ch <- firstVal
		}
	}
	return ch
}

func tryNotifyCh(ch chan int) {
	select {
	case ch <- 1:
	default:
	}
}

func (h *Hub) UnregisterClient(icao24 string, ch chan string) {
	h.Lock()
	defer h.Unlock()

	if h.icaos[icao24] == nil {
		return
	}

	psListener := h.icaos[icao24]
	psListener.RemoveWSClient(ch)

	length := len(psListener.wsClients)
	if length <= 0 {
		psListener.Cancel()

		tryNotifyCh(psListener.stopCh)

		psListener.stopCh <- 1
		h.icaos[icao24] = nil
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

	ch := hub.RegisterClient(icao24)

	var msg string

	for {
		msg = <-ch
		err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			fmt.Println(err)
			hub.UnregisterClient(icao24, ch)
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
