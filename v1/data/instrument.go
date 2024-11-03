package data

type InstrumentType int

const (
	Spot InstrumentType = iota
	Perpetual
	Future
	Option
)

var instrumentRepresentation = map[InstrumentType]string{
	Spot:      "spot",
	Perpetual: "perpetual",
	Future:    "future",
	Option:    "option",
}

func (instType InstrumentType) String() string {
	return instrumentRepresentation[instType]
}

type Instrument struct {
	Symbol         string
	BaseCurrency   string
	QuoteCurrency  string
	TickSize       string
	StepSize       string
	InstrumentType InstrumentType
}
