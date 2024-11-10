package ws

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	uninitialisedState uint32 = iota
	disconnectedState
	connectingState
	connectedState
)

type Client struct {
	enabled                  atomic.Bool
	state                    atomic.Uint32
	verbose                  bool
	connectionMonitorRunning atomic.Bool
	trafficTimeout           time.Duration
	connectionMonitorDelay   time.Duration
	proxyAddr                string
	defaultURL               string
	defaultURLAuth           string
	runningURL               string
	runningURLAuth           string
	exchangeName             string
	m                        sync.Mutex

	DataHandler chan interface{}
	ToRoutine   chan interface{}

	ShutdownC         chan struct{}
	Wg                sync.WaitGroup
	ReadMessageErrors chan error

	// Standard stream Connection
	Conn *Connection
	// Authenticated stream Connection
	AuthConn Connection
}

func (w *Client) setState(s uint32) {
	w.state.Store(s)
}

func (w *Client) IsConnected() bool {
	return w.state.Load() == connectedState
}

// IsConnecting returns whether the websocket is connecting
func (w *Client) IsConnecting() bool {
	return w.state.Load() == connectingState
}

func (w *Client) Connect(dialer *websocket.Dialer, pingHandler PingHandler) error {
	w.m.Lock()
	defer w.m.Unlock()
	return w.connect(dialer, pingHandler)
}

func (w *Client) connect(dialer *websocket.Dialer, pingHandler PingHandler) error {
	if w.IsConnecting() {
		return fmt.Errorf("%v %w", w.exchangeName, errAlreadyReconnecting)
	}
	if w.IsConnected() {
		return fmt.Errorf("%v %w", w.exchangeName, errAlreadyConnected)
	}
	w.setState(connectingState)

	w.Conn = NewConnection(w.exchangeName, w.defaultURL, w.verbose, w.ShutdownC, w.ReadMessageErrors, &w.Wg)
	err := w.Conn.Connect(dialer)
	if err != nil {
		w.setState(disconnectedState)
		return fmt.Errorf("%v Error Connecting %w", w.exchangeName, err)
	}
	w.setState(connectedState)

	return nil
}

func (w *Client) shutdown() {
	//	TODO clean state and shutdown
}

func (w *Client) reconnect() {
	//	TODO either shutdown or just clean state and connect
}

func (w *Client) monitorConnection() {
	for {
		select {
		case <-w.ShutdownC:
			log.Println("Shutting down monitorConnection")
			return
		case err := <-w.ReadMessageErrors:
			if errors.Is(err, errConnectionFault) {
				log.Printf("%v Websocket Disconnected\n", w.exchangeName)
			}
			if w.IsConnected() {
				//	TODO Reconnect Logic
			}

		}
	}
}
