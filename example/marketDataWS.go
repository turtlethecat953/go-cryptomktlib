package example

import (
	"context"
	"fmt"
	"go-cryptomktlib/v1/data"
	"go-cryptomktlib/v1/exchange/binanceusdm"
	"sync"
	"time"
)

func handleOrderBook(ch <-chan interface{}, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("Bye Bye")
		wg.Done()
	}()
	for val := range ch {
		ob, ok := val.(data.OrderBook)
		if !ok {
			continue
		}
		fmt.Printf("Handling OB %+v\n", ob)
	}
}

func main() {
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	binanceUsdm := binanceusdm.NewBinanceusdm(ctx, wg, "", "", true)
	orderBookCh := binanceUsdm.WatchOrderBook()
	wg.Add(1)
	go handleOrderBook(orderBookCh, wg)
	time.Sleep(1 * time.Second)
	cancel()
	wg.Wait()
}
