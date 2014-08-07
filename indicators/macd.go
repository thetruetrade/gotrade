// Moving Average Convergence and Divergence (Macd)
package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Moving Average Convergence-Divergence (Macd) Indicator
type Macd struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionMacd
	fastTimePeriod       int
	slowTimePeriod       int
	signalTimePeriod     int
	emaFast              *EmaWithoutStorage
	emaSlow              *EmaWithoutStorage
	emaSignal            *EmaWithoutStorage
	currentFastEma       float64
	currentSlowEma       float64
	currentMacd          float64
	emaSlowSkip          int
	selectData           gotrade.DataSelectionFunc

	// public variables
	Macd      []float64
	Signal    []float64
	Histogram []float64
}

// NewMacd creates a Moving Average Convergence Divergence Indicator (Macd) for online usage
func NewMacd(fastTimePeriod int, slowTimePeriod int, signalTimePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Macd, err error) {

	// the minimum fastTimePeriod for this indicator is 2
	if fastTimePeriod < 2 {
		return nil, errors.New("fastTimePeriod is less than the minimum (2)")
	}

	// check the maximum fastTimePeriod
	if fastTimePeriod > MaximumLookbackPeriod {
		return nil, errors.New("fastTimePeriod is greater than the maximum (100000)")
	}

	// the minimum slowTimePeriod for this indicator is 2
	if slowTimePeriod < 2 {
		return nil, errors.New("slowTimePeriod is less than the minimum (2)")
	}

	// check the maximum slowTimePeriod
	if slowTimePeriod > MaximumLookbackPeriod {
		return nil, errors.New("slowTimePeriod is greater than the maximum (100000)")
	}

	// the minimum signalTimePeriod for this indicator is 2
	if signalTimePeriod < 1 {
		return nil, errors.New("signalTimePeriod is less than the minimum (2)")
	}

	// check the maximum slowTimePeriod
	if signalTimePeriod > MaximumLookbackPeriod {
		return nil, errors.New("signalTimePeriod is greater than the maximum (100000)")
	}

	lookback := slowTimePeriod + signalTimePeriod - 2
	ind := Macd{
		baseIndicator:    newBaseIndicator(lookback),
		baseFloatBounds:  newBaseFloatBounds(),
		fastTimePeriod:   fastTimePeriod,
		slowTimePeriod:   slowTimePeriod,
		signalTimePeriod: signalTimePeriod,
	}

	// shift the fast ema up so that it has valid data at the same time as the slow emas
	ind.emaSlowSkip = slowTimePeriod - fastTimePeriod
	ind.emaFast, err = NewEmaWithoutStorage(fastTimePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentFastEma = dataItem
	})

	ind.emaSlow, err = NewEmaWithoutStorage(slowTimePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentSlowEma = dataItem

		ind.currentMacd = ind.currentFastEma - ind.currentSlowEma

		ind.emaSignal.ReceiveTick(ind.currentMacd, streamBarIndex)
	})

	ind.emaSignal, err = NewEmaWithoutStorage(signalTimePeriod, func(dataItem float64, streamBarIndex int) {

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}

		// Macd Line: (12-day EmaWithoutStorage - 26-day EmaWithoutStorage)

		// Signal Line: 9-day EmaWithoutStorage of Macd Line

		// Macd Histogram: Macd Line - Signal Line

		macd := ind.currentFastEma - ind.currentSlowEma
		signal := dataItem
		histogram := macd - signal

		// update the maximum result value
		if macd > ind.maxValue {
			ind.maxValue = macd
		}

		if signal > ind.maxValue {
			ind.maxValue = signal
		}

		if histogram > ind.maxValue {
			ind.maxValue = histogram
		}

		// update the minimum result value
		if macd < ind.minValue {
			ind.minValue = macd
		}

		if signal < ind.minValue {
			ind.minValue = signal
		}

		if histogram < ind.minValue {
			ind.minValue = histogram
		}

		// notify of a new result value though the value available action
		ind.valueAvailableAction(macd, signal, histogram, streamBarIndex)
	})

	ind.selectData = selectData
	ind.valueAvailableAction = func(dataItemMacd float64, dataItemSignal float64, dataItemHistogram float64, streamBarIndex int) {
		ind.Macd = append(ind.Macd, dataItemMacd)
		ind.Signal = append(ind.Signal, dataItemSignal)
		ind.Histogram = append(ind.Histogram, dataItemHistogram)
	}
	return &ind, err
}

// NewDefaultMacd creates a Moving Average Convergence Divergence Indicator (Macd) for online usage with default parameters
//	fastTimePeriod - 12
//	slowTimePeriod - 26
//	signalTimePeriod - 9
func NewDefaultMacd() (indicator *Macd, err error) {
	fastTimePeriod := 12
	slowTimePeriod := 26
	signalTimePeriod := 9
	return NewMacd(fastTimePeriod, slowTimePeriod, signalTimePeriod, gotrade.UseClosePrice)
}

// NewMacdWithSrcLen creates a Moving Average Convergence Divergence Indicator (Macd) for offline usage
func NewMacdWithSrcLen(sourceLength int, fastTimePeriod int, slowTimePeriod int, signalTimePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Macd, err error) {
	ind, err := NewMacd(fastTimePeriod, slowTimePeriod, signalTimePeriod, selectData)
	ind.Macd = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	ind.Signal = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	ind.Histogram = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultMacdWithSrcLen creates a Moving Average Convergence Divergence Indicator (Macd) for offline usage with default parameters
func NewDefaultMacdWithSrcLen(sourceLength int) (indicator *Macd, err error) {
	ind, err := NewDefaultMacd()
	ind.Macd = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	ind.Signal = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	ind.Histogram = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewMacdForStream creates a Moving Average Convergence Divergence Indicator (Macd) for online usage with a source data stream
func NewMacdForStream(priceStream *gotrade.DOHLCVStream, fastTimePeriod int, slowTimePeriod int, signalTimePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Macd, err error) {
	ind, err := NewMacd(fastTimePeriod, slowTimePeriod, signalTimePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultMacdForStream creates a Moving Average Convergence Divergence Indicator (Macd) for online usage with a source data stream
func NewDefaultMacdForStream(priceStream *gotrade.DOHLCVStream) (indicator *Macd, err error) {
	ind, err := NewDefaultMacd()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewMacdForStreamWithSrcLen creates a Moving Average Convergence Divergence Indicator (Macd) for offline usage with a source data stream
func NewMacdForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream, fastTimePeriod int, slowTimePeriod int, signalTimePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Macd, err error) {
	ind, err := NewMacdWithSrcLen(sourceLength, fastTimePeriod, slowTimePeriod, signalTimePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultMacdForStreamWithSrcLen creates a Moving Average Convergence Divergence Indicator (Macd) for offline usage with a source data stream
func NewDefaultMacdForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *Macd, err error) {
	ind, err := NewDefaultMacdWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *Macd) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *Macd) ReceiveTick(tickData float64, streamBarIndex int) {
	if streamBarIndex > ind.emaSlowSkip {
		ind.emaFast.ReceiveTick(tickData, streamBarIndex)
	}
	ind.emaSlow.ReceiveTick(tickData, streamBarIndex)
}
