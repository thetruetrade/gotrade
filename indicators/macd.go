// Moving Average Convergence and Divergence (MACD)
package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// MACD Line: (12-day EMA - 26-day EMA)

// Signal Line: 9-day EMA of MACD Line

// MACD Histogram: MACD Line - Signal Line

// A Moving Average Convergence-Divergence (MACD) Indicator
type MACD struct {
	*baseIndicatorWithLookback

	// private variables
	valueAvailableAction ValueAvailableActionMACD
	fastTimePeriod       int
	slowTimePeriod       int
	signalTimePeriod     int
	emaFast              *EMA
	emaSlow              *EMA
	emaSignal            *EMA
	currentFastEMA       float64
	currentSlowEMA       float64
	currentMACD          float64
	emaSlowSkip          int

	// public variables
	MACD      []float64
	Signal    []float64
	Histogram []float64
}

// NewMACD returns a new Moving Average Convergence-Divergence (MACD) Indicator configured with the
// specified timePeriod. The MACD results are stored in the DATA field.
func NewMACD(fastTimePeriod int, slowTimePeriod int, signalTimePeriod int, selectData gotrade.DataSelectionFunc) (indicator *MACD, err error) {
	newMACD := MACD{baseIndicatorWithLookback: newBaseIndicatorWithLookback(slowTimePeriod + signalTimePeriod - 2),
		fastTimePeriod:   fastTimePeriod,
		slowTimePeriod:   slowTimePeriod,
		signalTimePeriod: signalTimePeriod}

	// shift the fast ema up so that it has valid data at the same time as the slow emas
	newMACD.emaSlowSkip = slowTimePeriod - fastTimePeriod
	newMACD.emaFast, _ = NewEMA(fastTimePeriod, selectData)

	newMACD.emaFast.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newMACD.currentFastEMA = dataItem
	}

	newMACD.emaSlow, _ = NewEMA(slowTimePeriod, selectData)

	newMACD.emaSlow.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newMACD.currentSlowEMA = dataItem

		newMACD.currentMACD = newMACD.currentFastEMA - newMACD.currentSlowEMA

		newMACD.emaSignal.ReceiveTick(newMACD.currentMACD, streamBarIndex)
	}

	newMACD.emaSignal, _ = NewEMA(signalTimePeriod, selectData)

	newMACD.emaSignal.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newMACD.dataLength += 1
		if newMACD.validFromBar == -1 {
			newMACD.validFromBar = streamBarIndex
		}

		// MACD Line: (12-day EMA - 26-day EMA)

		// Signal Line: 9-day EMA of MACD Line

		// MACD Histogram: MACD Line - Signal Line

		macd := newMACD.currentFastEMA - newMACD.currentSlowEMA
		signal := dataItem
		histogram := macd - signal

		// MAX

		if macd > newMACD.maxValue {
			newMACD.maxValue = macd
		}

		if signal > newMACD.maxValue {
			newMACD.maxValue = signal
		}

		if histogram > newMACD.maxValue {
			newMACD.maxValue = histogram
		}

		// MIN

		if macd < newMACD.minValue {
			newMACD.minValue = macd
		}

		if signal < newMACD.minValue {
			newMACD.minValue = signal
		}

		if histogram < newMACD.minValue {
			newMACD.minValue = histogram
		}
		newMACD.valueAvailableAction(macd, signal, histogram, streamBarIndex)
	}

	newMACD.selectData = selectData
	newMACD.valueAvailableAction = func(dataItemMACD float64, dataItemSignal float64, dataItemHistogram float64, streamBarIndex int) {
		newMACD.MACD = append(newMACD.MACD, dataItemMACD)
		newMACD.Signal = append(newMACD.Signal, dataItemSignal)
		newMACD.Histogram = append(newMACD.Histogram, dataItemHistogram)
	}
	return &newMACD, nil
}

func NewMACDForStream(priceStream *gotrade.DOHLCVStream, fastTimePeriod int, slowTimePeriod int, signalTimePeriod int, selectData gotrade.DataSelectionFunc) (indicator *MACD, err error) {
	newMACD, err := NewMACD(fastTimePeriod, slowTimePeriod, signalTimePeriod, selectData)
	priceStream.AddTickSubscription(newMACD)
	return newMACD, err
}

func (ind *MACD) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *MACD) ReceiveTick(tickData float64, streamBarIndex int) {
	if streamBarIndex > ind.emaSlowSkip {
		ind.emaFast.ReceiveTick(tickData, streamBarIndex)
	}
	ind.emaSlow.ReceiveTick(tickData, streamBarIndex)
}
