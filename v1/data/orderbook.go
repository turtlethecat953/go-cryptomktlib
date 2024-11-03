package data

import "strconv"

type Level struct {
	Price    float64
	Quantity float64
}

type Book struct {
	Levels []Level
}

// OrderBook All Ts are in milliseconds
type OrderBook struct {
	Bids          Book
	Asks          Book
	ExchangeTs    int64
	TransactionTs int64
}

type IOrderBook interface {
	GetBids() *[][]string
	GetAsks() *[][]string
	GetExchangeTs() int64
	GetTransactionTs() int64
}

func ToOrderBook(rawBook IOrderBook) OrderBook {
	orderBook := OrderBook{
		ExchangeTs:    rawBook.GetExchangeTs(),
		TransactionTs: rawBook.GetTransactionTs(),
	}

	bids := Book{
		Levels: []Level{},
	}
	for _, bid := range *rawBook.GetBids() {
		price, err := strconv.ParseFloat(bid[0], 64)
		if err != nil {
			panic(err)
		}
		quantity, err := strconv.ParseFloat(bid[1], 64)
		if err != nil {
			panic(err)
		}
		bids.Levels = append(bids.Levels,
			Level{
				Price:    price,
				Quantity: quantity})
	}

	asks := Book{
		Levels: []Level{},
	}
	for _, ask := range *rawBook.GetAsks() {
		price, err := strconv.ParseFloat(ask[0], 64)
		if err != nil {
			panic(err)
		}
		quantity, err := strconv.ParseFloat(ask[1], 64)
		if err != nil {
			panic(err)
		}
		asks.Levels = append(asks.Levels,
			Level{
				Price:    price,
				Quantity: quantity})
	}

	orderBook.Bids = bids
	orderBook.Asks = asks
	return orderBook
}
