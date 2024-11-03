package binanceusdm

import (
	"go-cryptomktlib/v1/data"
)

type FapiV1TimeResponse struct {
	ServerTime int `json:"serverTime"`
}

type FapiV1ExchangeInfoResponse struct {
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	RateLimits      []struct {
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
		RateLimitType string `json:"rateLimitType"`
	} `json:"rateLimits"`
	ServerTime int64 `json:"serverTime"`
	Assets     []struct {
		Asset             string `json:"asset"`
		MarginAvailable   bool   `json:"marginAvailable"`
		AutoAssetExchange string `json:"autoAssetExchange"`
	} `json:"assets"`
	Symbols []struct {
		Symbol                string   `json:"symbol"`
		Pair                  string   `json:"pair"`
		ContractType          string   `json:"contractType"`
		DeliveryDate          int64    `json:"deliveryDate"`
		OnboardDate           int64    `json:"onboardDate"`
		Status                string   `json:"status"`
		MaintMarginPercent    string   `json:"maintMarginPercent"`
		RequiredMarginPercent string   `json:"requiredMarginPercent"`
		BaseAsset             string   `json:"baseAsset"`
		QuoteAsset            string   `json:"quoteAsset"`
		MarginAsset           string   `json:"marginAsset"`
		PricePrecision        int      `json:"pricePrecision"`
		QuantityPrecision     int      `json:"quantityPrecision"`
		BaseAssetPrecision    int      `json:"baseAssetPrecision"`
		QuotePrecision        int      `json:"quotePrecision"`
		UnderlyingType        string   `json:"underlyingType"`
		UnderlyingSubType     []string `json:"underlyingSubType"`
		SettlePlan            int      `json:"settlePlan"`
		TriggerProtect        string   `json:"triggerProtect"`
		Filters               []struct {
			FilterType        string `json:"filterType"`
			MaxPrice          string `json:"maxPrice,omitempty"`
			MinPrice          string `json:"minPrice,omitempty"`
			TickSize          string `json:"tickSize,omitempty"`
			MaxQty            string `json:"maxQty,omitempty"`
			MinQty            string `json:"minQty,omitempty"`
			StepSize          string `json:"stepSize,omitempty"`
			Limit             int    `json:"limit,omitempty"`
			Notional          string `json:"notional,omitempty"`
			MultiplierUp      string `json:"multiplierUp,omitempty"`
			MultiplierDown    string `json:"multiplierDown,omitempty"`
			MultiplierDecimal string `json:"multiplierDecimal,omitempty"`
		} `json:"filters"`
		OrderType       []string `json:"OrderType"`
		TimeInForce     []string `json:"timeInForce"`
		LiquidationFee  string   `json:"liquidationFee"`
		MarketTakeBound string   `json:"marketTakeBound"`
	} `json:"symbols"`
	Timezone string `json:"timezone"`
}

func (exchangeInfoResponse *FapiV1ExchangeInfoResponse) ToInstruments() *[]data.Instrument {
	var instruments []data.Instrument
	for _, exchangeInstrument := range exchangeInfoResponse.Symbols {
		tickSize := "0"
		stepSize := "0"
		for _, filter := range exchangeInstrument.Filters {
			if filter.FilterType == "PRICE_FILTER" {
				tickSize = filter.TickSize
			}
			if filter.FilterType == "LOT_SIZE" {
				stepSize = filter.StepSize
			}
		}

		instruments = append(instruments, data.Instrument{
			Symbol:         exchangeInstrument.Symbol,
			BaseCurrency:   exchangeInstrument.BaseAsset,
			QuoteCurrency:  exchangeInstrument.QuoteAsset,
			TickSize:       tickSize,
			StepSize:       stepSize,
			InstrumentType: data.Perpetual})
	}
	return &instruments
}

type FapiV1Depth struct {
	LastUpdateId int        `json:"lastUpdateId"`
	E            int64      `json:"E"`
	T            int64      `json:"T"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

func (response *FapiV1Depth) GetBids() *[][]string {
	return &response.Bids
}

func (response *FapiV1Depth) GetAsks() *[][]string {
	return &response.Asks
}

func (response *FapiV1Depth) GetExchangeTs() int64 {
	return response.E
}

func (response *FapiV1Depth) GetTransactionTs() int64 {
	return response.T
}
