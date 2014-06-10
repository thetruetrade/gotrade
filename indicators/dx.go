// Directional Movement Index (DX)
package indicators

// DX = ( (+DI)-(-DI) ) / ( (+DI) + (-DI) )

import (
	"github.com/thetruetrade/gotrade"
	"math"
)

type DXWithoutStorage struct {
	*baseIndicator
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableAction
	periodCounter        int
	lookbackCounter      int
	minusDI              *MinusDI
	plusDI               *PlusDI
	currentPlusDI        float64
	currentMinusDI       float64
}

func NewDXWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableAction) (indicator *DXWithoutStorage, err error) {
	var lookback int = 2
	if timePeriod > 1 {
		lookback = timePeriod
	}

	newDX := DXWithoutStorage{baseIndicator: newBaseIndicator(timePeriod),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               lookback * -1,
		lookbackCounter:             -2,
		currentPlusDI:               0.0,
		currentMinusDI:              0.0}

	newDX.valueAvailableAction = valueAvailableAction

	newDX.minusDI, _ = NewMinusDI(timePeriod)

	newDX.minusDI.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newDX.currentMinusDI = dataItem
	}

	newDX.plusDI, _ = NewPlusDI(timePeriod)

	newDX.plusDI.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newDX.currentPlusDI = dataItem

		var result float64
		tmp := newDX.currentMinusDI + newDX.currentPlusDI
		if tmp != 0.0 {
			result = 100.0 * (math.Abs(newDX.currentMinusDI-newDX.currentPlusDI) / tmp)
		} else {
			result = 0.0
		}

		newDX.dataLength += 1

		if newDX.validFromBar == -1 {
			newDX.validFromBar = streamBarIndex
		}

		if result > newDX.maxValue {
			newDX.maxValue = result
		}

		if result < newDX.minValue {
			newDX.minValue = result
		}
		newDX.valueAvailableAction(result, streamBarIndex)

	}

	return &newDX, nil
}

// A Directional Movement Indicator
type DX struct {
	*DXWithoutStorage

	// public variables
	Data []float64
}

// NewDX returns a new Directional Movement Indicator (DX) configured with the
// specified timePeriod. The DX results are stored in the DATA field.
func NewDX(timePeriod int) (indicator *DX, err error) {

	newDX := DX{}
	newDX.DXWithoutStorage, err = NewDXWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			newDX.Data = append(newDX.Data, dataItem)
		})

	return &newDX, err
}

func NewDXForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *DX, err error) {
	newDX, err := NewDX(timePeriod)
	priceStream.AddTickSubscription(newDX)
	return newDX, err
}

func (ind *DXWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.minusDI.ReceiveDOHLCVTick(tickData, streamBarIndex)
	ind.plusDI.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
