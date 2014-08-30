package gotrade

import (
	"time"
)

type DOHLCV interface {
	D() time.Time
	O() float64
	H() float64
	L() float64
	C() float64
	V() float64
}

type OHLCV interface {
	O() float64
	H() float64
	L() float64
	C() float64
	V() float64
}

type OHLC interface {
	O() float64
	H() float64
	L() float64
	C() float64
}

type DOHLCVDataItem struct {
	date        time.Time
	openPrice   float64
	highPrice   float64
	lowPrice    float64
	closePrice  float64
	volumePrice float64
}

func NewDOHLCVDataItem(date time.Time, openPrice float64, highPrice float64, lowPrice float64, closePrice float64, volume float64) *DOHLCVDataItem {
	return &DOHLCVDataItem{date, openPrice, highPrice, lowPrice, closePrice, volume}
}

func (di *DOHLCVDataItem) D() time.Time {
	return di.date
}

func (di *DOHLCVDataItem) O() float64 {
	return di.openPrice
}

func (di *DOHLCVDataItem) H() float64 {
	return di.highPrice
}

func (di *DOHLCVDataItem) L() float64 {
	return di.lowPrice
}

func (di *DOHLCVDataItem) C() float64 {
	return di.closePrice
}

func (di *DOHLCVDataItem) V() float64 {
	return di.volumePrice
}

// A function that selects which data property to use from a DOHLCV data structure
type DOHLCVDataSelectionFunc func(dataItem DOHLCV) float64

// Close price DOHLCV data selector
func UseClosePrice(dataItem DOHLCV) float64 {
	return dataItem.C()
}

// Open price DOHLCV data selector
func UseOpenPrice(dataItem DOHLCV) float64 {
	return dataItem.O()
}

// High price DOHLCV data selector
func UseHighPrice(dataItem DOHLCV) float64 {
	return dataItem.H()
}

// Low price DOHLCV data selector
func UseLowPrice(dataItem DOHLCV) float64 {
	return dataItem.L()
}

// Volume DOHLCV data selector
func UseVolume(dataItem DOHLCV) float64 {
	return dataItem.V()
}

// Consumer of DOHLCV Ticks
type DOHLCVTickReceiver interface {
	ReceiveDOHLCVTick(tickData DOHLCV, streamBarIndex int)
}

// Consumer of a float tick
type TickReceiver interface {
	ReceiveTick(tickData float64, streamBarIndex int)
}
