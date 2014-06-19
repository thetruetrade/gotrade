/* Package GoTrade implements

*/
package gotrade

import (
	"math"
	"sync"
	"time"
)

type DOHLCVStreamTickReceiver interface {
	ReceiveTick(tickData DOHLCV)
}

type DataStreamHolder interface {
	MinValue() float64
	MaxValue() float64
}

type DateStreamHolder interface {
	MinDate() time.Time
	MaxDate() time.Time
}

type interDayBarType int

type intraDayBarType int

const (
	MinuteBar intraDayBarType = iota
	DailyBar  interDayBarType = iota
	WeeklyBar
	MonthlyBar
)

type DOHLCVStream struct {
	Data           []DOHLCV
	subscribers    []DOHLCVTickReceiver
	streamBarIndex int
	minValue       float64
	maxValue       float64
}

type InterDayDOHLCVStream struct {
	*DOHLCVStream
	streamBarType interDayBarType
}

func NewInterDayDOHLCVStream(streamBarType interDayBarType) *InterDayDOHLCVStream {
	s := InterDayDOHLCVStream{DOHLCVStream: &DOHLCVStream{streamBarIndex: 0,
		minValue: math.MaxFloat64,
		maxValue: math.SmallestNonzeroFloat64},
		streamBarType: streamBarType}
	return &s
}

func NewDailyDOHLCVStream() *InterDayDOHLCVStream {
	return NewInterDayDOHLCVStream(DailyBar)
}

func NewWeeklyDOHLCVStream() *InterDayDOHLCVStream {
	return NewInterDayDOHLCVStream(WeeklyBar)
}

func NewMonthlyDOHLCVStream() *InterDayDOHLCVStream {
	return NewInterDayDOHLCVStream(MonthlyBar)
}

func (p *DOHLCVStream) ReceiveTick(tickData DOHLCV) {
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

type IntraDayDOHLCVStream struct {
	*DOHLCVStream
	intraDayBarInterval int
}

func NewIntraDayDOHLCVStream(barIntervalInMins int) *IntraDayDOHLCVStream {
	s := IntraDayDOHLCVStream{DOHLCVStream: &DOHLCVStream{streamBarIndex: 0,
		minValue: math.MaxFloat64,
		maxValue: math.SmallestNonzeroFloat64},
		intraDayBarInterval: barIntervalInMins}
	return &s
}
