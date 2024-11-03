package base

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

type WsClient struct {
	baseUrl string
	conn    *websocket.Conn
}

func NewWsClient(baseUrl string, handleMessage func([]byte)) *WsClient {
	ws := &WsClient{baseUrl: baseUrl}

	Dialer := websocket.Dialer{
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: false,
	}
	fmt.Printf("Base URL %s", baseUrl)
	c, _, err := Dialer.Dial(baseUrl, nil)
	if err != nil {
		panic(err)
	}
	c.SetReadLimit(655350)
	ws.conn = c
	//doneC := make(chan struct{})
	//stopC := make(chan struct{})
	go func() {
		// This function will exit either on error from
		// websocket.Conn.ReadMessage or when the stopC channel is
		// closed by the client.
		defer c.Close()
		// Wait for the stopC channel to be closed.  We do that in a
		// separate goroutine because ReadMessage is a blocking
		// operation.
		silent := false
		//go func() {
		//	select {
		//	case <-stopC:
		//		silent = true
		//	case <-doneC:
		//	}
		//	c.Close()
		//}()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if !silent {
					fmt.Printf("Error %v", err)
				}
				return
			}
			handleMessage(message)
		}
	}()
	return ws
}

func (w *WsClient) SendMessage(data []byte) error {
	return w.conn.WriteMessage(websocket.TextMessage, data)
}
