package main

import (
	"context"
	"sync"
)

type Hub struct {
	sync.RWMutex
	icaos map[string]*PubSubListener
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

		newListener := NewPubSubListener(cl.icao24, cancellableCtx, listenerCancel)
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

	// Notify all clients of the error and stop the writer/listeners goroutines
	ls.SendAllErr("Error: server was unable to receive positions")

	// Delete the PubSubListener from the map of active listeners
	delete(h.icaos, ls.icao24)
}

func (h *Hub) MakeListenerTimeoutHandler(ls *PubSubListener) func() {
	return func() {
		Log.Debug.Println("Pubsub subscription timed out, subID: ", ls.subId)
		ls.listenerCancel()
		h.UnregisterFailedListener(ls)
	}
}
