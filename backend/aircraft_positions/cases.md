Fail cases:

- pubSubListener:
  - 1. Error on listen function:
    - Stop timeout goroutine
    - Send close msg to all clients
    - Stop all client readers
    - Stop all client writers
    - Remove icao and listener from hub map.
  - 2. Message timeout:
    - Stop Listener
    - Send close msg to all clients
    - Stop all client readers
    - Stop all client writers
    - Remove icao and listener from hub map.
  - 3. Be cancelled by hub due to no clients left:
    - Stop timeout goroutine
- WS Reader:
  - 1. close msg from client:
    - Stop WS writer
    - close WS
    - remove client from listener
    - Potentially stop listener
    - Potentially stop listener timeout goroutine
    - Potentially remove icao and listener from hub map
  - 2. error reading msg from client:
    - Stop WS writer
    - close WS
    - remove client from listener
    - Potentially stop listener
    - Potentially stop listener timeout goroutine
    - Potentially remove icao and listener from hub map
- WS Writer:
  - 1. error writing msg from client:
    - Stop WS reader
    - close WS
    - remove client from listener
    - Potentially stop listener
    - Potentially stop listener timeout goroutine
    - Potentially remove icao and listener from hub map

listenerCanceller()
- done() listener context

timer.StopTimer():
- done() on timer context

timeOutHandler():
- Stop the 

cl.SendErrAndClose():
- Send close msg to client
- Cancel client so that Reader is stopped.

hub.UnregisterFailedListener(ls):
- Send close msg to all clients
- Stop all client writers
- Stop all client readers
- Remove icao and listener from hub map.


hub.UnregisterClient(cl):
    - Returns if no listener is found for that icao
    - Returns if no matching client is found for that listener
    - Cancels the timer and the listener if the icao is out of clients.
    - Removes the listener reference from the icao

TODO list:
    - test locally pubsub emulator
      - create topic 
      - create msg every 15 seconds
      - msg have to be the given format
      - potentially enable on the project
    - find a way to test websocket 
