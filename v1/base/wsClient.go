package base

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"sync"
	"time"
)

type WsClient struct {
	conn    *websocket.Conn
	logger  *log.Logger
	closed  bool
	verbose bool
}

func NewWsClient(
	ctx context.Context,
	wg *sync.WaitGroup,
	baseUrl string,
	handleMessage func([]byte),
	onClose func(),
	verbose bool) *WsClient {
	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

	ws := &WsClient{
		logger:  logger,
		verbose: verbose,
		closed:  false}

	Dialer := websocket.Dialer{
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: false,
	}
	ws.log("Connecting to %s", baseUrl)
	c, _, err := Dialer.Dial(baseUrl, nil)
	if err != nil {
		panic(err)
	}
	c.SetReadLimit(655350)
	ws.conn = c
	wg.Add(1)
	go ws.readLoop(ctx, wg, handleMessage, onClose)
	return ws
}

func (w *WsClient) log(format string, v ...interface{}) {
	if w.verbose {
		w.logger.Printf(format, v...)
	}
}

func (w *WsClient) close() {
	if w.conn == nil {
		w.log("Trying to close a closed connection...")
	}
	err := w.conn.Close()
	if err != nil {
		fmt.Printf("Error During Closing Websocket Connection %v", err)
	}
	w.log("Closed Websocket Client")
	w.closed = true
}

func (w *WsClient) readLoop(ctx context.Context, wg *sync.WaitGroup, handleMessage func([]byte), onClose func()) {
	defer func() {
		onClose()
		w.close()
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			w.log("Stopping ReadLoop Gracefully...")
			return
		default:
			if ctx.Err() != nil {
				return
			}
			_, message, err := w.conn.ReadMessage()
			if err != nil {
				w.log("Closing ReadLoop due to Irrecoverable Error %v", err)
				return
			}
			if w.closed {
				return
			}
			handleMessage(message)
		}
	}
}

func (w *WsClient) SendMessage(data []byte) error {
	return w.conn.WriteMessage(websocket.TextMessage, data)
}
