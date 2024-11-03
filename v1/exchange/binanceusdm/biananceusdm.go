package binanceusdm

import (
	"context"
	"encoding/json"
	"fmt"
	"go-cryptomktlib/v1/base"
	"go-cryptomktlib/v1/data"
	"net/http"
)

const TEST_URL = "https://testnet.binancefuture.com"
const TEST_WS_URL = "wss://fstream.binancefuture.com"
const WS_URL = "wss://fstream.binance.com"

type Binanceusdm struct {
	key        string
	secret     string
	httpClient *base.Client
	wsClient   *base.WsClient
	verbose    bool
	handleChan map[string]chan<- interface{}
}

func NewBinanceusdm(key, secret string, verbose bool) *Binanceusdm {
	return &Binanceusdm{key: key, secret: secret, verbose: verbose, handleChan: make(map[string]chan<- interface{})}
}

func (exchange *Binanceusdm) httpConnect() error {
	if exchange.httpClient != nil {
		return nil
	}

	httpClient, err := base.NewClient(exchange.verbose)
	if err != nil {
		return err
	}

	exchange.httpClient = httpClient
	return nil
}

func (exchange *Binanceusdm) sign(r *base.Request) error {
	fullUrl := fmt.Sprintf("%s/%s", TEST_URL, r.Endpoint)
	r.EncodeParams()
	queryString := r.ParamString
	if r.Method == http.MethodGet || r.Method == http.MethodDelete {
		fullUrl += "?" + queryString
	}
	r.Url = fullUrl
	return nil
}

func (exchange *Binanceusdm) fetch(ctx context.Context, r *base.Request) ([]byte, error) {
	err := exchange.httpConnect()
	if err != nil {
		panic(err)
	}
	err = exchange.sign(r)
	if err != nil {
		panic(err)
	}
	respByte, getErr := exchange.httpClient.Do(ctx, r)
	if getErr != nil {
		panic(getErr)
	}
	return respByte, nil
}

func (exchange *Binanceusdm) KeepAlive(ctx context.Context) (FapiV1TimeResponse, error) {
	request := base.NewRequest(http.MethodGet, "fapi/v1/time", true)

	respByte, err := exchange.fetch(ctx, request)
	if err != nil {
		panic(err)
	}

	jsonResponse := FapiV1TimeResponse{}
	err = json.Unmarshal(respByte, &jsonResponse)
	if err != nil {
		panic(err)
	}
	return jsonResponse, nil
}

func (exchange *Binanceusdm) FetchMarket(ctx context.Context) (*[]data.Instrument, *FapiV1ExchangeInfoResponse, error) {
	request := base.NewRequest(http.MethodGet, "fapi/v1/exchangeInfo", true)

	respByte, err := exchange.fetch(ctx, request)
	if err != nil {
		panic(err)
	}

	jsonResponse := FapiV1ExchangeInfoResponse{}
	err = json.Unmarshal(respByte, &jsonResponse)
	if err != nil {
		panic(err)
	}

	return jsonResponse.ToInstruments(), &jsonResponse, nil
}

func (exchange *Binanceusdm) FetchOrderBook(ctx context.Context, symbol string) (data.OrderBook, *FapiV1Depth, error) {
	request := base.NewRequest(http.MethodGet, "fapi/v1/depth", true)
	request.SetParam("symbol", symbol)

	respByte, err := exchange.fetch(ctx, request)
	if err != nil {
		panic(err)
	}

	jsonResponse := FapiV1Depth{}
	err = json.Unmarshal(respByte, &jsonResponse)
	if err != nil {
		panic(err)
	}
	return data.ToOrderBook(&jsonResponse), &jsonResponse, nil
}

func (exchange *Binanceusdm) handleMessage(message []byte) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal(message, &jsonMap)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", jsonMap)
	eventType, ok := jsonMap["e"]
	if ok {
		if eventType == "depthUpdate" {
			ch, ok := exchange.handleChan["partialDiff"]
			if ok {
				jsonResponse := PartialBookDepth{}
				err = json.Unmarshal(message, &jsonResponse)
				ch <- data.ToOrderBook(&jsonResponse)
			}
		}
	}
}

func (exchange *Binanceusdm) stream() string {
	return "0"
}

func (exchange *Binanceusdm) connect() error {
	if exchange.wsClient != nil {
		return nil
	}
	w := base.NewWsClient(WS_URL+"/ws/"+exchange.stream(), exchange.handleMessage)
	exchange.wsClient = w
	return nil
}

func (exchange *Binanceusdm) WatchOrderBook(ch chan<- interface{}) error {
	_ = exchange.connect()

	request := Request{
		Method: "SUBSCRIBE",
		Id:     1,
		Params: []string{"btcusdt@depth"},
	}
	s, _ := json.Marshal(request)
	exchange.wsClient.SendMessage(s)

	exchange.handleChan["partialDiff"] = ch
	return nil
}
