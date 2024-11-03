package base

type exchange interface {
	FetchMarket()
	FetchOrderBook()
	FetchOrders()
	FetchMyTrades()
	CreateOrder()
	CancelOrder()
	CancelAllOrder()
}
