// Exponential Moving Average (EMA)
package indicators

import (
	"github.com/thetruetrade/gotrade"
	"math"
)

type emaBase struct {
	Indicator

	// public variables
	LookbackPeriod int

	// private variables
	periodTotal   float64
	periodCounter int
	multiplier    float64
	previousEMA   float64
}

// An Exponential Moving Average Indicator
type EMA struct {
	emaBase

	// public variables
	Data []float64
}

// NewEMA returns a new Exponential Moving Average (EMA) configured with the
// specified lookbackPeriod
func NewEMA(lookbackPeriod int, transformData gotrade.DataTransformationFunc) (indicator *EMA, err error) {
	newEMA := EMA{emaBase: emaBase{LookbackPeriod: lookbackPeriod,
		periodCounter: lookbackPeriod * -1,
		multiplier:    float64(2.0 / float64(lookbackPeriod+1.0)),
		Indicator:     Indicator{validFromBar: -1, transformData: transformData, minValue: math.MaxFloat64, maxValue: math.SmallestNonzeroFloat64}}}

	newEMA.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newEMA.Data = append(newEMA.Data, dataItem)
	}

	return &newEMA, nil
}

func NewEMAForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, transformData gotrade.DataTransformationFunc) (indicator *EMA, err error) {
	newEma, err := NewEMA(lookbackPeriod, transformData)
	priceStream.AddSubscription(newEma)
	return newEma, err
}

func (ema *emaBase) RecieveOrderedTick(dataItem gotrade.DOHLCV, streamBarIndex int) {
	ema.periodCounter += 1
	ema.dataLength += 1
	var transformedData = ema.transformData(dataItem)

	if transformedData > ema.maxValue {
		ema.maxValue = transformedData
	}

	if transformedData < ema.minValue {
		ema.minValue = transformedData
	}

	if ema.periodCounter < 0 {
		ema.periodTotal += transformedData
	} else if ema.periodCounter == 0 {
		ema.periodTotal += transformedData
		ema.previousEMA = ema.periodTotal / float64(ema.LookbackPeriod)
		ema.valueAvailableAction(ema.previousEMA, streamBarIndex)

	} else if ema.periodCounter > 0 {
		result := (transformedData-ema.previousEMA)*ema.multiplier + ema.previousEMA
		ema.previousEMA = result
		ema.valueAvailableAction(ema.previousEMA, streamBarIndex)
	}
}
