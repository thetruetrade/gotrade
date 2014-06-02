// Chainkin Oscillator (ChainkinOsc)
// this should be as simple as EMA(ADL,3) - EMA(ADL,10), however it seems the emas are intialised with the
// first adl value and not offset like the macd to conincide, they are both calculated from the 2nd bar and used before their
// lookback period is reached - so the emas are calcualted inline and not using the general EMAWithoutStorage
package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type ChainkinOscWithoutStorage struct {
	*baseIndicator

	// private variables
	fastTimePeriod       int
	slowTimePeriod       int
	valueAvailableAction ValueAvailableAction
	adl                  *ADLWithoutStorage
	emaFast              float64
	emaSlow              float64
	emaFastMultiplier    float64
	emaSlowMultiplier    float64
	periodCounter        int
	isInitialised        bool
}

func NewChainkinOscWithoutStorage(fastTimePeriod int, slowTimePeriod int, valueAvailableAction ValueAvailableAction) (indicator *ChainkinOscWithoutStorage, err error) {
	newChainkinOsc := ChainkinOscWithoutStorage{baseIndicator: newBaseIndicator(slowTimePeriod - 1),
		slowTimePeriod:    slowTimePeriod,
		fastTimePeriod:    fastTimePeriod,
		emaFastMultiplier: float64(2.0 / float64(fastTimePeriod+1.0)),
		emaSlowMultiplier: float64(2.0 / float64(slowTimePeriod+1.0)),
		periodCounter:     slowTimePeriod * -1,
		isInitialised:     false}

	newChainkinOsc.valueAvailableAction = valueAvailableAction

	newChainkinOsc.adl, err = NewADLWithoutStorage(func(dataItem float64, streamBarIndex int) {
		newChainkinOsc.periodCounter += 1

		if !newChainkinOsc.isInitialised {
			newChainkinOsc.emaFast = dataItem
			newChainkinOsc.emaSlow = dataItem
			newChainkinOsc.isInitialised = true
		}
		if newChainkinOsc.periodCounter < 0 {
			newChainkinOsc.emaFast = (dataItem-newChainkinOsc.emaFast)*newChainkinOsc.emaFastMultiplier + newChainkinOsc.emaFast
			newChainkinOsc.emaSlow = (dataItem-newChainkinOsc.emaSlow)*newChainkinOsc.emaSlowMultiplier + newChainkinOsc.emaSlow
		}

		if newChainkinOsc.periodCounter >= 0 {
			newChainkinOsc.dataLength += 1
			if newChainkinOsc.validFromBar == -1 {
				newChainkinOsc.validFromBar = streamBarIndex
			}

			newChainkinOsc.emaFast = (dataItem-newChainkinOsc.emaFast)*newChainkinOsc.emaFastMultiplier + newChainkinOsc.emaFast
			newChainkinOsc.emaSlow = (dataItem-newChainkinOsc.emaSlow)*newChainkinOsc.emaSlowMultiplier + newChainkinOsc.emaSlow
			chaikinOsc := newChainkinOsc.emaFast - newChainkinOsc.emaSlow

			if chaikinOsc > newChainkinOsc.maxValue {
				newChainkinOsc.maxValue = chaikinOsc
			}

			if chaikinOsc < newChainkinOsc.minValue {
				newChainkinOsc.minValue = chaikinOsc
			}
			newChainkinOsc.valueAvailableAction(chaikinOsc, streamBarIndex)
		}
	})

	return &newChainkinOsc, nil
}

// A Double Exponential Moving Average Indicator
type ChainkinOsc struct {
	*ChainkinOscWithoutStorage

	// public variables
	Data []float64
}

// NewChainkinOsc returns a new Double Exponential Moving Average (ChainkinOsc) configured with the
// specified timePeriod. The ChainkinOsc results are stored in the DATA field.
func NewChainkinOsc(fastTimePeriod int, slowTimePeriod int) (indicator *ChainkinOsc, err error) {

	newChainkinOsc := ChainkinOsc{}
	newChainkinOsc.ChainkinOscWithoutStorage, err = NewChainkinOscWithoutStorage(fastTimePeriod, slowTimePeriod,
		func(dataItem float64, streamBarIndex int) {
			newChainkinOsc.Data = append(newChainkinOsc.Data, dataItem)
		})

	return &newChainkinOsc, err
}

func NewChainkinOscForStream(priceStream *gotrade.DOHLCVStream, fastTimePeriod int, slowTimePeriod int) (indicator *ChainkinOsc, err error) {
	newChainkinOsc, err := NewChainkinOsc(fastTimePeriod, slowTimePeriod)
	priceStream.AddTickSubscription(newChainkinOsc)
	return newChainkinOsc, err
}

func (ind *ChainkinOscWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.adl.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
