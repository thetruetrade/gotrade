package gotrade

import (
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

type DOHLCVStreamSubscriber interface {
	RecieveOrderedTick(dataItem DOHLCV, streamBarIndex int)
}

type DOHLCVStream struct {
	Data           []DOHLCV
	subscribers    []DOHLCVStreamSubscriber
	streamBarIndex int
}

func NewDOHLCVStream() *DOHLCVStream {
	return &DOHLCVStream{streamBarIndex: 0}
}

func (p *DOHLCVStream) RecieveOrderedTick(tickData DOHLCV) {
	p.streamBarIndex++
	p.Data = append(p.Data, tickData)
	var waitGroup sync.WaitGroup

	// notify all the subscribers and wait
	for subscriberIndex := range p.subscribers {
		waitGroup.Add(1)
		var subscriber DOHLCVStreamSubscriber = p.subscribers[subscriberIndex]
		go func(subscriber DOHLCVStreamSubscriber) {
			defer waitGroup.Done()
			subscriber.RecieveOrderedTick(tickData, p.streamBarIndex)

		}(subscriber)
	}

	waitGroup.Wait()
}

func (p *DOHLCVStream) RecieveTick(tickData DOHLCV) {

}

func (p *DOHLCVStream) MinDate() time.Time {
	// do some checks here, return an error object too
	return p.Data[0].D()
}

func (p *DOHLCVStream) MaxDate() time.Time {
	// do some checks here, return an error object too
	return p.Data[len(p.Data)-1].D()
}

func (p *DOHLCVStream) AddSubscription(subscriber DOHLCVStreamSubscriber) {
	p.subscribers = append(p.subscribers, subscriber)
}

func (p *DOHLCVStream) RemoveSubscription(subscriber DOHLCVStreamSubscriber) {

}
