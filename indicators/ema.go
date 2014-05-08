// Exponential Moving Average (EMA)
package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type baseEMA struct {
	*baseIndicatorWithLookback

	// private variables
	periodTotal          float64
	periodCounter        int
	multiplier           float64
	previousEMA          float64
	valueAvailableAction ValueAvailableAction
}

func newBaseEMA(lookbackPeriod int) *baseEMA {
	newEMA := baseEMA{baseIndicatorWithLookback: newBaseIndicatorWithLookback(lookbackPeriod),
		periodCounter: lookbackPeriod * -1,
		multiplier:    float64(2.0 / float64(lookbackPeriod+1.0))}

	return &newEMA
}

// An Exponential Moving Average Indicator
type EMA struct {
	*baseEMA

	// public variables
	Data []float64
}

// NewEMA returns a new Exponential Moving Average (EMA) configured with the
// specified lookbackPeriod
func NewEMA(lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *EMA, err error) {
	newEMA := EMA{baseEMA: newBaseEMA(lookbackPeriod)}
	newEMA.selectData = selectData
	newEMA.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newEMA.Data = append(newEMA.Data, dataItem)
	}

	return &newEMA, nil
}

func NewEMAForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *EMA, err error) {
	newEma, err := NewEMA(lookbackPeriod, selectData)
	priceStream.AddTickSubscription(newEma)
	return newEma, err
}

func (ema *baseEMA) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ema.selectData(tickData)
	ema.ReceiveTick(selectedData, streamBarIndex)
}

func (ema *baseEMA) ReceiveTick(tickData float64, streamBarIndex int) {
	ema.periodCounter += 1
	if ema.periodCounter < 0 {
		ema.periodTotal += tickData
	} else if ema.periodCounter == 0 {
		ema.dataLength += 1

		if ema.validFromBar == -1 {
			ema.validFromBar = streamBarIndex
		}

		ema.periodTotal += tickData
		result := ema.periodTotal / float64(ema.LookbackPeriod)
		ema.previousEMA = result

		if result > ema.maxValue {
			ema.maxValue = result
		}

		if result < ema.minValue {
			ema.minValue = result
		}

		ema.valueAvailableAction(ema.previousEMA, streamBarIndex)

	} else if ema.periodCounter > 0 {
		ema.dataLength += 1

		result := (tickData-ema.previousEMA)*ema.multiplier + ema.previousEMA
		ema.previousEMA = result

		if result > ema.maxValue {
			ema.maxValue = result
		}

		if result < ema.minValue {
			ema.minValue = result
		}

		ema.valueAvailableAction(ema.previousEMA, streamBarIndex)
	}
}
