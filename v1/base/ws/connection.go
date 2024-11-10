package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Connection struct {
	conn               *websocket.Conn
	isConnected        atomic.Bool
	exchangeName       string
	url                string
	shutDownCh         <-chan struct{}
	readMessageErrorCh chan<- error
	writeLock          sync.Mutex
	wg                 *sync.WaitGroup
	verbose            bool
}

func NewConnection(exchangeName string, url string, verbose bool,
	shutDownCh <-chan struct{}, readMessageErrorCh chan<- error, wg *sync.WaitGroup) *Connection {
	w := Connection{
		exchangeName:       exchangeName,
		url:                url,
		verbose:            verbose,
		shutDownCh:         shutDownCh,
		readMessageErrorCh: readMessageErrorCh,
		wg:                 wg}
	w.setConnected(false)
	return &w
}

func (c *Connection) setConnected(state bool) {
	c.isConnected.Store(state)
}

func (c *Connection) IsConnected() bool {
	return c.isConnected.Load()
}

func (c *Connection) Connect(dialer *websocket.Dialer) error {
	var err error
	var conStatus *http.Response
	c.conn, conStatus, err = dialer.DialContext(context.Background(), c.url, http.Header{})
	if err != nil {
		if conStatus != nil {
			_ = conStatus.Body.Close()
			return fmt.Errorf("%s websocket Connection: %v %v %v Error: %c", c.exchangeName, c.url, conStatus, conStatus.StatusCode, err)
		}
		return fmt.Errorf("%s websocket Connection: %v Error: %c", c.exchangeName, c.url, err)
	}
	_ = conStatus.Body.Close()
	if c.verbose {
		log.Printf("%v Websocket connected to %s\n", c.exchangeName, c.url)
	}
	return nil
}

func (c *Connection) ReadMessage() ([]byte, error) {
	_, resp, err := c.conn.ReadMessage()
	if err != nil {
		c.setConnected(false)
		select {
		case c.readMessageErrorCh <- errConnectionFault:
		default:
			log.Printf("%s failed to relay websocket error: %v\n", c.exchangeName, err)
		}
		return nil, err
	}
	return resp, nil
}

func (c *Connection) Shutdown() error {
	if c == nil || c.conn == nil {
		return nil
	}
	c.setConnected(false)
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	return c.conn.Close()
}

func (c *Connection) SendJSONMessage(data any) error {
	return c.writeToConn(func() error {
		if c.verbose {
			if msg, err := json.Marshal(data); err == nil { // WriteJSON will error for us anyway
				log.Printf("%v %v: Sending message: %v", c.exchangeName, c.url, string(msg))
			}
		}
		return c.conn.WriteJSON(data)
	})
}

// SendRawMessage sends a message over the Connection without JSON encoding it
func (c *Connection) SendRawMessage(messageType int, message []byte) error {
	return c.writeToConn(func() error {
		if c.verbose {
			log.Printf("%v %v: Sending message: %v", c.exchangeName, c.url, string(message))
		}
		return c.conn.WriteMessage(messageType, message)
	})
}

func (c *Connection) writeToConn(writeConn func() error) error {
	if !c.IsConnected() {
		return fmt.Errorf("%v websocket Connection: cannot send message as Connection is disconnected", c.exchangeName)
	}

	// TODO Add Rate Limiting Logic

	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	return writeConn()
}

func (c *Connection) SetupPingHandler(handler PingHandler) {
	if handler.UseDefaultHandler {
		c.conn.SetPingHandler(func(msg string) error {
			err := c.conn.WriteControl(handler.MessageType, []byte(msg), time.Now().Add(handler.Delay))
			if err == websocket.ErrCloseSent {
				return nil
			} else if e, ok := err.(net.Error); ok && e.Timeout() {
				return nil
			}
			return err
		})
		return
	}
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(handler.Delay)
		for {
			select {
			case <-c.shutDownCh:
				ticker.Stop()
				return
			case <-ticker.C:
				err := c.SendRawMessage(handler.MessageType, handler.Message)
				if err != nil {
					log.Printf("%v websocket Connection: ping handler failed to send message [%s]: %v", c.exchangeName, handler.Message, err)
					return
				}
			}
		}
	}()
}
