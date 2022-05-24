package main

import (
	gcp "aircraft_positions/gcp"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"sync"
	"time"
)

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

	filter := "attributes.icao24=\"" + state.icao24 + "\""

	Log.Debug.Println("Creating new pubsub subscription: ", state.subId)

	sub, err := gcp.CreateSubscription(state.subId, posTopic, filter)
	defer gcp.DeleteSubscription(state.subId, posTopic)
	defer Log.Debug.Println("Deleting pubsub subscription: ", state.subId)
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
		Log.Error.Println("Error in pubsub Receive: ", err)
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

	idx := slices.IndexFunc(state.wsClients, func(c *WSClient) bool { return c.cliId == cli.cliId })

	if idx >= 0 {
		length := len(state.wsClients) - 1
		state.wsClients[idx] = state.wsClients[length]
		state.wsClients = state.wsClients[:length]
	} else {
		Log.Error.Println("Failed to remove WSClient from pbListener with icao=", cli.icao24, ": Client not found in list.")
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

func NewPubSubListener(icao24 string, ctx context.Context, stopListen func()) *PubSubListener {
	timeoutCtx, timeoutCancel := context.WithCancel(context.Background())

	/* SubIDs are required to:
	start with a letter,
	be longer than 3 and shorted than 255 characters,
	not start with "goog"
	contain only Letters [A-Za-z], numbers [0-9], dashes -, underscores _, periods ., tildes ~, plus signs +, and percent signs %
	*/
	subID := "sub_" + icao24 + "_" + uuid.New().String()

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
