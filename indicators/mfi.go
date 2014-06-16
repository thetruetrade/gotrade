// Average True Range (MFI)
package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

// A plus DM Indicator
type MFIWithoutStorage struct {
	*baseIndicator
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableAction
	periodCounter        int
	typicalPrice         *TypicalPriceWithoutStorage
	positiveMoneyFlow    float64
	negativeMoneyFlow    float64
	positiveHistory      *list.List
	negativeHistory      *list.List
	previousTypicalPrice float64
	currentVolume        float64
}

// NewMFIWithoutStorage returns a new Money Flow Index (MFI) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The MFI results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewMFIWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableAction) (indicator *MFIWithoutStorage, err error) {
	newMFI := MFIWithoutStorage{baseIndicator: newBaseIndicator(timePeriod),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               (timePeriod * -1) - 1,
		positiveHistory:             list.New(),
		negativeHistory:             list.New(),
		positiveMoneyFlow:           0.0,
		negativeMoneyFlow:           0.0,
		currentVolume:               0.0,
		previousTypicalPrice:        0.0}
	newMFI.typicalPrice, err = NewTypicalPriceWithoutStorage(func(dataItem float64, streamBarIndex int) {
		newMFI.periodCounter += 1

		if newMFI.periodCounter > (newMFI.GetTimePeriod() * -1) {
			moneyFlow := dataItem * newMFI.currentVolume

			if newMFI.periodCounter <= 0 {
				if dataItem > newMFI.previousTypicalPrice {
					newMFI.positiveMoneyFlow += moneyFlow
					newMFI.positiveHistory.PushBack(moneyFlow)
					newMFI.negativeHistory.PushBack(0.0)
				} else if dataItem < newMFI.previousTypicalPrice {
					newMFI.negativeMoneyFlow += moneyFlow
					newMFI.positiveHistory.PushBack(0.0)
					newMFI.negativeHistory.PushBack(moneyFlow)
				} else {
					newMFI.positiveHistory.PushBack(0.0)
					newMFI.negativeHistory.PushBack(0.0)
				}
			}

			if newMFI.periodCounter == 0 {

				result := 100.0 * (newMFI.positiveMoneyFlow / (newMFI.positiveMoneyFlow + newMFI.negativeMoneyFlow))

				newMFI.dataLength += 1
				if newMFI.validFromBar == -1 {
					newMFI.validFromBar = streamBarIndex
				}

				if result > newMFI.maxValue {
					newMFI.maxValue = result
				}

				if result < newMFI.minValue {
					newMFI.minValue = result
				}

				newMFI.valueAvailableAction(result, streamBarIndex)
			}
			if newMFI.periodCounter > 0 {
				firstPositive := newMFI.positiveHistory.Front().Value.(float64)
				newMFI.positiveMoneyFlow -= firstPositive

				firstNegative := newMFI.negativeHistory.Front().Value.(float64)
				newMFI.negativeMoneyFlow -= firstNegative

				if dataItem > newMFI.previousTypicalPrice {
					newMFI.positiveMoneyFlow += moneyFlow
					newMFI.positiveHistory.PushBack(moneyFlow)
					newMFI.negativeHistory.PushBack(0.0)
				} else if dataItem < newMFI.previousTypicalPrice {
					newMFI.negativeMoneyFlow += moneyFlow
					newMFI.positiveHistory.PushBack(0.0)
					newMFI.negativeHistory.PushBack(moneyFlow)
				} else {
					newMFI.positiveHistory.PushBack(0.0)
					newMFI.negativeHistory.PushBack(0.0)
				}

				result := 100.0 * (newMFI.positiveMoneyFlow / (newMFI.positiveMoneyFlow + newMFI.negativeMoneyFlow))

				newMFI.dataLength += 1

				if result > newMFI.maxValue {
					newMFI.maxValue = result
				}

				if result < newMFI.minValue {
					newMFI.minValue = result
				}

				newMFI.valueAvailableAction(result, streamBarIndex)
			}

		}
		newMFI.previousTypicalPrice = dataItem

		if newMFI.positiveHistory.Len() > newMFI.GetTimePeriod() {
			first := newMFI.positiveHistory.Front()
			newMFI.positiveHistory.Remove(first)
		}

		if newMFI.negativeHistory.Len() > newMFI.GetTimePeriod() {
			first := newMFI.negativeHistory.Front()
			newMFI.negativeHistory.Remove(first)
		}
	})
	newMFI.valueAvailableAction = valueAvailableAction

	return &newMFI, nil
}

// An Average True Range Indicator
type MFI struct {
	*MFIWithoutStorage

	// public variables
	Data []float64
}

// NewMFI returns a new Money Flow Index (MFI) configured with the
// specified timePeriod. The MFI results are stored in the Data field.
func NewMFI(timePeriod int) (indicator *MFI, err error) {
	newMFI := MFI{}
	newMFI.MFIWithoutStorage, err = NewMFIWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		newMFI.Data = append(newMFI.Data, dataItem)
	})

	return &newMFI, err
}

func NewMFIForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *MFI, err error) {
	newMFI, err := NewMFI(timePeriod)
	priceStream.AddTickSubscription(newMFI)
	return newMFI, err
}

func (ind *MFIWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.currentVolume = tickData.V()
	ind.typicalPrice.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
