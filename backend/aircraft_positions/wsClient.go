package main

import (
	"context"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	icao24       string
	cliId        string
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
		Log.Debug.Println("Failed to notify WSWriteLoop of error: ", msg)
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
				cl.clientCancel()
				cl.ws.Close()
				hub.UnregisterClient(cl)
				Log.Error.Println("Error writing to ws with icao24=", cl.icao24, ", err: ", err)
				return
			}

		case msg = <-cl.chErr:
			cl.ws.WriteMessage(websocket.CloseMessage, []byte(msg))
			cl.clientCancel()
			cl.ws.Close()
			Log.Debug.Println("Closing WS with icao24=", cl.icao24, " due some error (positions unavailable).")
			return

		case <-cl.clientCtx.Done():
			cl.ws.WriteMessage(websocket.CloseMessage, []byte("Closing the websocket"))
			cl.ws.Close()
			Log.Debug.Println("Closing WS with icao24=", cl.icao24, ": operation cancelled.")
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
			Log.Debug.Println("Error reading from websocket with icao24=", cl.icao24, " likely closed by someone, err: ", err)
			cl.clientCancel()
			cl.ws.Close()
			hub.UnregisterClient(cl)
			return
		}

		if msgType == websocket.CloseMessage {
			cl.clientCancel()
			cl.ws.Close()
			hub.UnregisterClient(cl)
			Log.Debug.Println("Closing websocket with icao24=", cl.icao24, " due to client close request.")
			return
		}
	}
}
