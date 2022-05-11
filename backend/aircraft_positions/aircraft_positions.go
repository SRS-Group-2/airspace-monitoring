package main

import (
	gcp "aircraft_positions/gcp"
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

const env_credJson = "GOOGLE_APPLICATION_CREDENTIALS"

const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"
const env_topicID = "GOOGLE_PUBSUB_AIRCRAFT_LIST_TOPIC_ID"
const env_subID = "GOOGLE_PUBSUB_AIRCRAFT_LIST_SUBSCRIBER_ID"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

var posTopic *pubsub.Topic

type Hub struct {
	sync.RWMutex
	icaos map[string]*PubSubListener
	SubID string
}

type PubSubListener struct {
	sync.RWMutex
	Icao24    string
	SubID     string
	wsClients []WSClient
	cache     string
	//ctx       context.Context
	Cancel func()
}

func (h *Hub) RegisterClient(icao24 string) chan string {
	h.Lock()
	defer h.Unlock()

	ch := make(chan string, 100)

	if h.icaos[icao24] == nil {
		psListener := &PubSubListener{}
		psListener.Icao24 = icao24
		psListener.SubID = h.SubID

		var ctx context.Context
		ctx, psListener.Cancel = context.WithCancel(context.Background())

		h.icaos[icao24] = psListener

		h.icaos[icao24].AddWSClient(ch)

		go startPubSubListener(ctx, psListener)

	} else {
		h.icaos[icao24].AddWSClient(ch)
	}
	return ch
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
		h.icaos[icao24] = nil
	}
}

type WSClient chan string

func (c *WSClient) Send(msg string) {
	*c <- msg
}

func (pbl *PubSubListener) AddWSClient(cli WSClient) {
	pbl.Lock()
	defer pbl.Unlock()
	pbl.wsClients = append(pbl.wsClients, cli)
}

func (pbl *PubSubListener) RemoveWSClient(cli WSClient) {
	pbl.Lock()
	defer pbl.Unlock()
	idx := slices.IndexFunc(pbl.wsClients, func(c WSClient) bool { return c == cli })
	if idx >= 0 {
		length := len(pbl.wsClients) - 1
		pbl.wsClients[idx] = pbl.wsClients[length]
		pbl.wsClients = append(pbl.wsClients)
		pbl.wsClients = pbl.wsClients[:length]
	} else {
		fmt.Println("No client found in pbListener WSClients list")
	}
}

func (pbl *PubSubListener) Clients() []WSClient {
	pbl.RLock()
	defer pbl.RUnlock()
	return pbl.wsClients
}

func (pbl *PubSubListener) SendAll(msg string) {
	pbl.RLock()
	defer pbl.RUnlock()
	for i := 0; i < len(pbl.wsClients); i++ {
		pbl.wsClients[i].Send(msg)
	}
}

func startPubSubListener(ctx context.Context, state *PubSubListener) {
	filter := `"attributes.icao24 = "` + state.Icao24 + `"`
	sub, err := gcp.CreateSubscription(state.SubID+"_"+state.Icao24, posTopic, filter)
	checkErr(err)

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		state.SendAll(string(msg.Data))
	})

	if err != nil {
		fmt.Println("Error in Pub/Sub Receive: ", err)
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
	router.GET("/airspace/aircraft/position/:icao24", httpRequestHandler)

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
			// TODO: fix it
			hub.UnregisterClient(ch)
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
