// Exponential Moving Average (EMA)
package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// An Exponential Moving Average Indicator
type EMAWithoutStorage struct {
	*baseIndicator
	*baseIndicatorWithLookback
	*baseIndicatorWithTimePeriod

	// private variables
	periodTotal          float64
	periodCounter        int
	multiplier           float64
	previousEMA          float64
	valueAvailableAction ValueAvailableAction
}

func NewEMAWithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableAction) (indicator *EMAWithoutStorage, err error) {
	newEMA := EMAWithoutStorage{baseIndicator: newBaseIndicator(),
		baseIndicatorWithLookback:   newBaseIndicatorWithLookback(timePeriod - 1),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               timePeriod * -1,
		multiplier:                  float64(2.0 / float64(timePeriod+1.0))}
	newEMA.selectData = selectData
	newEMA.valueAvailableAction = valueAvailableAction

	return &newEMA, err
}

// An Exponential Moving Average Indicator
type EMA struct {
	*EMAWithoutStorage

	// public variables
	Data []float64
}

// NewEMA returns a new Exponential Moving Average (EMA) configured with the
// specified timePeriod. The EMA results are stored in the Data field.
func NewEMA(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *EMA, err error) {
	newEMA := EMA{}
	newEMA.EMAWithoutStorage, err = NewEMAWithoutStorage(timePeriod, selectData,
		func(dataItem float64, streamBarIndex int) {
			newEMA.Data = append(newEMA.Data, dataItem)
		})

	return &newEMA, err
}

func NewEMAForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *EMA, err error) {
	newEma, err := NewEMA(timePeriod, selectData)
	priceStream.AddTickSubscription(newEma)
	return newEma, err
}

func (ema *EMAWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ema.selectData(tickData)
	ema.ReceiveTick(selectedData, streamBarIndex)
}

func (ema *EMAWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ema.periodCounter += 1
	if ema.periodCounter < 0 {
		ema.periodTotal += tickData
	} else if ema.periodCounter == 0 {
		ema.dataLength += 1

		if ema.validFromBar == -1 {
			ema.validFromBar = streamBarIndex
		}

		ema.periodTotal += tickData
		result := ema.periodTotal / float64(ema.GetTimePeriod())
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
