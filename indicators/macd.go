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
	*baseIndicator

	// private variables
	valueAvailableAction ValueAvailableActionMACD
	fastLookbackPeriod   int
	slowLookbackPeriod   int
	signalLookbackPeriod int
	emaFast              *EMA
	emaSlow              *EMA
	emaSignal            *EMA
	currentFastEMA       float64
	currentSlowEMA       float64
	currentMACD          float64
	emaSlowSkip          int

	// public variables
	Data []MACDData
}

// NewMACD returns a new Moving Average Convergence-Divergence (MACD) Indicator configured with the
// specified lookbackPeriods. The MACD results are stored in the DATA field.
func NewMACD(fastLookbackPeriod int, slowLookbackPeriod int, signalLookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *MACD, err error) {
	newMACD := MACD{baseIndicator: newBaseIndicator(),
		fastLookbackPeriod:   fastLookbackPeriod,
		slowLookbackPeriod:   slowLookbackPeriod,
		signalLookbackPeriod: signalLookbackPeriod}

	// shift the fast ema up so that it has valid data at the same time as the slow emas
	newMACD.emaSlowSkip = slowLookbackPeriod - fastLookbackPeriod
	newMACD.emaFast, _ = NewEMA(fastLookbackPeriod, selectData)

	newMACD.emaFast.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newMACD.currentFastEMA = dataItem
	}

	newMACD.emaSlow, _ = NewEMA(slowLookbackPeriod, selectData)

	newMACD.emaSlow.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newMACD.currentSlowEMA = dataItem

		newMACD.currentMACD = newMACD.currentFastEMA - newMACD.currentSlowEMA

		newMACD.emaSignal.ReceiveTick(newMACD.currentMACD, streamBarIndex)
	}

	newMACD.emaSignal, _ = NewEMA(signalLookbackPeriod, selectData)

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
		newMACD.valueAvailableAction(NewMACDDataItem(macd, signal, histogram), streamBarIndex)
	}

	newMACD.selectData = selectData
	newMACD.valueAvailableAction = func(dataItem MACDData, streamBarIndex int) {
		newMACD.Data = append(newMACD.Data, dataItem)
	}
	return &newMACD, nil
}

func NewMACDForStream(priceStream *gotrade.DOHLCVStream, fastLookbackPeriod int, slowLookbackPeriod int, signalLookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *MACD, err error) {
	newMACD, err := NewMACD(fastLookbackPeriod, slowLookbackPeriod, signalLookbackPeriod, selectData)
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
