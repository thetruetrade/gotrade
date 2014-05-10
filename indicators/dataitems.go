package indicators

type MACDData interface {
	// MACD
	M() float64
	// Signal
	S() float64
	// Histogram
	H() float64
}

type MACDDataItem struct {
	macd      float64
	signal    float64
	histogram float64
}

func (data *MACDDataItem) M() float64 {
	return data.macd
}

func (data *MACDDataItem) S() float64 {
	return data.signal
}

func (data *MACDDataItem) H() float64 {
	return data.histogram
}

func NewMACDDataItem(macd float64, signal float64, histogram float64) *MACDDataItem {
	return &MACDDataItem{macd: macd, signal: signal, histogram: histogram}
}

type BollingerBand interface {
	// Upper bollinger band
	U() float64
	// Middle bollinger band
	M() float64
	// Lower bollinger band
	L() float64
}

type BollingerBandDataItem struct {
	upperBand  float64
	middleBand float64
	lowerBand  float64
}

func (bb *BollingerBandDataItem) U() float64 {
	return bb.upperBand
}

func (bb *BollingerBandDataItem) L() float64 {
	return bb.lowerBand
}

func (bb *BollingerBandDataItem) M() float64 {
	return bb.middleBand
}

func NewBollingerBandDataItem(upperBand float64, middleBand float64, lowerBand float64) *BollingerBandDataItem {
	return &BollingerBandDataItem{upperBand: upperBand, middleBand: middleBand, lowerBand: lowerBand}
}
