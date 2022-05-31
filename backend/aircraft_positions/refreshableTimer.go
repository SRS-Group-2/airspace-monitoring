package main

import (
	"context"
	"time"
)

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
				Log.Debug.Println("Stopping timer, operation cancelled.")
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
