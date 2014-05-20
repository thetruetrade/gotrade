package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type AroonOsc struct {
	*baseIndicatorWithLookback

	//private variables
	aroon *AroonWithoutStorage
	Data  []float64
}

func NewAroonOsc(lookbackPeriod int) (indicator *AroonOsc, err error) {
	ind := AroonOsc{baseIndicatorWithLookback: newBaseIndicatorWithLookback(lookbackPeriod + 1)}

	ind.aroon, err = NewAroonWithoutStorage(lookbackPeriod,
		func(dataItemAroonUp float64, dataItemAroonDown float64, streamBarIndex int) {
			ind.dataLength++

			result := dataItemAroonUp - dataItemAroonDown
			if ind.validFromBar == -1 {
				ind.validFromBar = streamBarIndex
			}

			if result > ind.maxValue {
				ind.maxValue = result
			}

			if result < ind.minValue {
				ind.minValue = result
			}

			ind.Data = append(ind.Data, result)
		})
	return &ind, nil
}

func (ind *AroonOsc) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.aroon.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
