package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

type LinearRegWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	periodCounter        int
	periodHistory        *list.List
	sumX                 float64
	sumXSquare           float64
	divisor              float64
	valueAvailableAction ValueAvailableActionLinearReg
}

func NewLinearRegWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionLinearReg) (indicator *LinearRegWithoutStorage, err error) {
	newVar := LinearRegWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod - 1),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               (timePeriod) * -1,
		periodHistory:               list.New()}

	timePeriodF := float64(timePeriod)
	timePeriodFMinusOne := timePeriodF - 1.0
	newVar.sumX = timePeriodF * timePeriodFMinusOne * 0.5
	newVar.sumXSquare = timePeriodF * timePeriodFMinusOne * (2.0*timePeriodF - 1.0) / 6.0
	newVar.divisor = newVar.sumX*newVar.sumX - timePeriodF*newVar.sumXSquare

	newVar.valueAvailableAction = valueAvailableAction

	return &newVar, nil
}

type LinearReg struct {
	*LinearRegWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

func NewLinearReg(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinearReg, err error) {
	newVar := LinearReg{selectData: selectData}
	newVar.LinearRegWithoutStorage, err = NewLinearRegWithoutStorage(timePeriod,
		func(dataItem float64, slope float64, intercept float64, streamBarIndex int) {
			newVar.Data = append(newVar.Data, dataItem)

			if dataItem > newVar.LinearRegWithoutStorage.maxValue {
				newVar.LinearRegWithoutStorage.maxValue = dataItem
			}

			if dataItem < newVar.LinearRegWithoutStorage.minValue {
				newVar.LinearRegWithoutStorage.minValue = dataItem
			}
		})

	return &newVar, err
}

func NewLinearRegForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinearReg, err error) {
	newVar, err := NewLinearReg(timePeriod, selectData)
	priceStream.AddTickSubscription(newVar)
	return newVar, err
}

func (ind *LinearReg) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *LinearRegWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1

	if ind.periodCounter >= 0 {
		sumXY := 0.0
		sumY := 0.0
		i := ind.GetTimePeriod()
		var value float64 = 0.0
		for e := ind.periodHistory.Front(); e != nil; e = e.Next() {
			i--
			value = e.Value.(float64)
			sumY += value
			sumXY += (float64(i) * value)
		}
		sumY += tickData
		timePeriod := float64(ind.GetTimePeriod())
		m := (timePeriod*sumXY - ind.sumX*sumY) / ind.divisor
		b := (sumY - m*ind.sumX) / timePeriod
		result := b + m*float64(timePeriod-1.0)

		ind.dataLength += 1
		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}

		ind.valueAvailableAction(result, m, b, streamBarIndex)
	}

	ind.periodHistory.PushBack(tickData)

	if ind.periodHistory.Len() >= ind.GetTimePeriod() {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)
	}

}
