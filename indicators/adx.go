// Average Directional Movement Index (ADX)
package indicators

// ADX = ( (+DI)-(-DI) ) / ( (+DI) + (-DI) )

import (
	"github.com/thetruetrade/gotrade"
)

type ADXWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableAction
	periodCounter        int
	dx                   *DX
	currentDX            float64
	sumDX                float64
	previousADX          float64
}

func NewADXWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableAction) (indicator *ADXWithoutStorage, err error) {

	newADX := ADXWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds((2 * timePeriod) - 1),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               timePeriod * -1,
		currentDX:                   0.0,
		sumDX:                       0.0,
		previousADX:                 0.0}

	newADX.valueAvailableAction = valueAvailableAction

	newADX.dx, _ = NewDX(timePeriod)

	newADX.dx.valueAvailableAction = func(dataItem float64, streamBarIndex int) {

		newADX.currentDX = dataItem

		newADX.periodCounter += 1
		if newADX.periodCounter < 0 {
			newADX.sumDX += newADX.currentDX
		} else if newADX.periodCounter == 0 {
			newADX.dataLength += 1

			if newADX.validFromBar == -1 {
				newADX.validFromBar = streamBarIndex
			}

			newADX.sumDX += newADX.currentDX
			result := newADX.sumDX / float64(newADX.GetTimePeriod())
			if result > newADX.maxValue {
				newADX.maxValue = result
			}

			if result < newADX.minValue {
				newADX.minValue = result
			}
			newADX.valueAvailableAction(result, streamBarIndex)
			newADX.previousADX = result

		} else {

			newADX.dataLength += 1

			result := (newADX.previousADX*float64(newADX.GetTimePeriod()-1) + newADX.currentDX) / float64(newADX.GetTimePeriod())
			if result > newADX.maxValue {
				newADX.maxValue = result
			}

			if result < newADX.minValue {
				newADX.minValue = result
			}
			newADX.valueAvailableAction(result, streamBarIndex)
			newADX.previousADX = result
		}

	}

	return &newADX, nil
}

// A Directional Movement Indicator
type ADX struct {
	*ADXWithoutStorage

	// public variables
	Data []float64
}

// NewADX returns a new Directional Movement Indicator (ADX) configured with the
// specified timePeriod. The ADX results are stored in the DATA field.
func NewADX(timePeriod int) (indicator *ADX, err error) {

	newADX := ADX{}
	newADX.ADXWithoutStorage, err = NewADXWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			newADX.Data = append(newADX.Data, dataItem)
		})

	return &newADX, err
}

func NewADXForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *ADX, err error) {
	newADX, err := NewADX(timePeriod)
	priceStream.AddTickSubscription(newADX)
	return newADX, err
}

func (ind *ADXWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.dx.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
