package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Money Flow Index Indicator (Mfi), no storage, for use in other indicators
type MfiWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	typicalPrice         *TypPriceWithoutStorage
	positiveMoneyFlow    float64
	negativeMoneyFlow    float64
	positiveHistory      *list.List
	negativeHistory      *list.List
	previousTypPrice     float64
	currentVolume        float64
	timePeriod           int
}

// NewMfiWithoutStorage creates a Money Flow Index Indicator (Mfi) without storage
func NewMfiWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *MfiWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for this indicator is 2
	if timePeriod < 2 {
		return nil, errors.New("timePeriod is less than the minimum (2)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	lookback := timePeriod
	ind := MfiWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        (timePeriod * -1) - 1,
		positiveHistory:      list.New(),
		negativeHistory:      list.New(),
		positiveMoneyFlow:    0.0,
		negativeMoneyFlow:    0.0,
		currentVolume:        0.0,
		previousTypPrice:     0.0,
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	ind.typicalPrice, err = NewTypPriceWithoutStorage(func(dataItem float64, streamBarIndex int) {
		ind.periodCounter += 1

		if ind.periodCounter > (ind.timePeriod * -1) {
			moneyFlow := dataItem * ind.currentVolume

			if ind.periodCounter <= 0 {
				if dataItem > ind.previousTypPrice {
					ind.positiveMoneyFlow += moneyFlow
					ind.positiveHistory.PushBack(moneyFlow)
					ind.negativeHistory.PushBack(0.0)
				} else if dataItem < ind.previousTypPrice {
					ind.negativeMoneyFlow += moneyFlow
					ind.positiveHistory.PushBack(0.0)
					ind.negativeHistory.PushBack(moneyFlow)
				} else {
					ind.positiveHistory.PushBack(0.0)
					ind.negativeHistory.PushBack(0.0)
				}
			}

			if ind.periodCounter == 0 {

				result := 100.0 * (ind.positiveMoneyFlow / (ind.positiveMoneyFlow + ind.negativeMoneyFlow))

				// increment the number of results this indicator can be expected to return
				ind.dataLength += 1
				if ind.validFromBar == -1 {
					// set the streamBarIndex from which this indicator returns valid results
					ind.validFromBar = streamBarIndex
				}

				// update the maximum result value
				if result > ind.maxValue {
					ind.maxValue = result
				}

				// update the minimum result value
				if result < ind.minValue {
					ind.minValue = result
				}

				// notify of a new result value though the value available action
				ind.valueAvailableAction(result, streamBarIndex)
			}
			if ind.periodCounter > 0 {
				firstPositive := ind.positiveHistory.Front().Value.(float64)
				ind.positiveMoneyFlow -= firstPositive

				firstNegative := ind.negativeHistory.Front().Value.(float64)
				ind.negativeMoneyFlow -= firstNegative

				if dataItem > ind.previousTypPrice {
					ind.positiveMoneyFlow += moneyFlow
					ind.positiveHistory.PushBack(moneyFlow)
					ind.negativeHistory.PushBack(0.0)
				} else if dataItem < ind.previousTypPrice {
					ind.negativeMoneyFlow += moneyFlow
					ind.positiveHistory.PushBack(0.0)
					ind.negativeHistory.PushBack(moneyFlow)
				} else {
					ind.positiveHistory.PushBack(0.0)
					ind.negativeHistory.PushBack(0.0)
				}

				result := 100.0 * (ind.positiveMoneyFlow / (ind.positiveMoneyFlow + ind.negativeMoneyFlow))

				// increment the number of results this indicator can be expected to return
				ind.dataLength += 1

				// update the maximum result value
				if result > ind.maxValue {
					ind.maxValue = result
				}

				// update the minimum result value
				if result < ind.minValue {
					ind.minValue = result
				}

				// notify of a new result value though the value available action
				ind.valueAvailableAction(result, streamBarIndex)
			}

		}
		ind.previousTypPrice = dataItem

		if ind.positiveHistory.Len() > ind.timePeriod {
			first := ind.positiveHistory.Front()
			ind.positiveHistory.Remove(first)
		}

		if ind.negativeHistory.Len() > ind.timePeriod {
			first := ind.negativeHistory.Front()
			ind.negativeHistory.Remove(first)
		}
	})

	return &ind, err
}

// A Money Flow Index Indicator (Mfi)
type Mfi struct {
	*MfiWithoutStorage

	// public variables
	Data []float64
}

// NewMfi creates a Money Flow Index Indicator (Mfi) for online usage
func NewMfi(timePeriod int) (indicator *Mfi, err error) {
	ind := Mfi{}
	ind.MfiWithoutStorage, err = NewMfiWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultMfi creates a Money Flow Index Indicator (Mfi) for online usage with default parameters
//	- timePeriod: 25
func NewDefaultMfi() (indicator *Mfi, err error) {
	timePeriod := 25
	return NewMfi(timePeriod)
}

// NewMfiWithSrcLen creates a Money Flow Index Indicator (Mfi) for offline usage
func NewMfiWithSrcLen(sourceLength uint, timePeriod int) (indicator *Mfi, err error) {
	ind, err := NewMfi(timePeriod)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultMfiWithSrcLen creates a Money Flow Index Indicator (Mfi) for offline usage with default parameters
func NewDefaultMfiWithSrcLen(sourceLength uint) (indicator *Mfi, err error) {
	ind, err := NewDefaultMfi()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewMfiForStream creates a Money Flow Index Indicator (Mfi) for online usage with a source data stream
func NewMfiForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *Mfi, err error) {
	ind, err := NewMfi(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultMfiForStream creates a Money Flow Index Indicator (Mfi) for online usage with a source data stream
func NewDefaultMfiForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Mfi, err error) {
	ind, err := NewDefaultMfi()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewMfiForStreamWithSrcLen creates a Money Flow Index Indicator (Mfi) for offline usage with a source data stream
func NewMfiForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *Mfi, err error) {
	ind, err := NewMfiWithSrcLen(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultMfiForStreamWithSrcLen creates a Money Flow Index Indicator (Mfi) for offline usage with a source data stream
func NewDefaultMfiForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Mfi, err error) {
	ind, err := NewDefaultMfiWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *MfiWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.currentVolume = tickData.V()
	ind.typicalPrice.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
