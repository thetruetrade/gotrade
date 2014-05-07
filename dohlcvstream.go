package gotrade

import (
	"math"
	"sync"
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
type DataSelectionFunc func(dataItem DOHLCV) float64

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

type DOHLCVTickReceiver interface {
	ReceiveDOHLCVTick(tickData DOHLCV, streamBarIndex int)
}

type TickReceiver interface {
	ReceiveTick(tickData float64, streamBarIndex int)
}

type DataStreamHolder interface {
	MinValue() float64
	MaxValue() float64
}

type DateStreamHolder interface {
	MinDate() time.Time
	MaxDate() time.Time
}

type DOHLCVStream struct {
	Data           []DOHLCV
	subscribers    []DOHLCVTickReceiver
	streamBarIndex int
	minValue       float64
	maxValue       float64
}

func NewDOHLCVStream() *DOHLCVStream {
	s := DOHLCVStream{streamBarIndex: 0}
	s.minValue = math.MaxFloat64
	s.maxValue = math.SmallestNonzeroFloat64
	return &s
}

func (p *DOHLCVStream) RecieveTick(tickData DOHLCV) {
	p.streamBarIndex++
	p.Data = append(p.Data, tickData)

	if p.minValue > tickData.L() {
		p.minValue = tickData.L()
	}

	if p.maxValue < tickData.H() {
		p.maxValue = tickData.H()
	}

	var waitGroup sync.WaitGroup

	// notify all the subscribers and wait
	for subscriberIndex := range p.subscribers {
		waitGroup.Add(1)
		var subscriber DOHLCVTickReceiver = p.subscribers[subscriberIndex]
		go func(subscriber DOHLCVTickReceiver) {
			defer waitGroup.Done()
			subscriber.ReceiveDOHLCVTick(tickData, p.streamBarIndex)

		}(subscriber)
	}

	waitGroup.Wait()
}

func (p *DOHLCVStream) MinDate() time.Time {
	// do some checks here, return an error object too
	return p.Data[0].D()
}

func (p *DOHLCVStream) MaxDate() time.Time {
	// do some checks here, return an error object too
	return p.Data[len(p.Data)-1].D()
}

func (p *DOHLCVStream) MinValue() float64 {
	return p.minValue
}

func (p *DOHLCVStream) MaxValue() float64 {
	return p.maxValue
}

func (p *DOHLCVStream) AddTickSubscription(subscriber DOHLCVTickReceiver) {
	p.subscribers = append(p.subscribers, subscriber)
}

func (p *DOHLCVStream) RemoveTickSubscription(subscriber DOHLCVTickReceiver) {

}
