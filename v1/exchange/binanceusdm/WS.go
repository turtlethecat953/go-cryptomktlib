package binanceusdm

type Request struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
}

type PartialBookDepth struct {
	E               string     `json:"e"`
	EventTime       int64      `json:"E"`
	TransactionTime int64      `json:"T"`
	S               string     `json:"s"`
	B               [][]string `json:"b"`
	A               [][]string `json:"a"`
}

func (p *PartialBookDepth) GetBids() *[][]string {
	return &p.B
}

func (p *PartialBookDepth) GetAsks() *[][]string {
	return &p.A
}

func (p *PartialBookDepth) GetExchangeTs() int64 {
	return p.EventTime
}

func (p *PartialBookDepth) GetTransactionTs() int64 {
	return p.TransactionTime
}
