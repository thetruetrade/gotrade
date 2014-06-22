package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

type LinearRegWithoutStorage struct {
	*baseIndicator
	*baseIndicatorWithTimePeriod

	// private variables
	periodCounter        int
	periodHistory        *list.List
	sumX                 float64
	sumXSquare           float64
	divisor              float64
	valueAvailableAction ValueAvailableActionLinearReg
}

func NewLinearRegWithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableActionLinearReg) (indicator *LinearRegWithoutStorage, err error) {
	newVar := LinearRegWithoutStorage{baseIndicator: newBaseIndicator(timePeriod - 1),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               (timePeriod) * -1,
		periodHistory:               list.New()}

	//timePeriodF := float64(timePeriod)
	//timePeriodFMinusOne := timePeriodF - 1.0
	var sumX float64 = 14.0 * 13.0 * 0.5
	newVar.sumX = sumX
	//newVar.sumX = timePeriodF * timePeriodFMinusOne * 0.5
	var sumXSquare float64 = 14.0 * 13.0 * (2.0*14 - 1.0) / 6.0
	newVar.sumXSquare = sumXSquare
	//newVar.sumXSquare = timePeriodF * timePeriodFMinusOne * (2.0*timePeriodF - 1.0) / 6.0
	var divisor float64 = sumX*sumX - 14.0*sumXSquare
	//newVar.divisor = //newVar.sumX*newVar.sumX - timePeriodF*newVar.sumXSquare
	newVar.divisor = divisor

	newVar.selectData = selectData
	newVar.valueAvailableAction = valueAvailableAction

	return &newVar, nil
}

type LinearReg struct {
	*LinearRegWithoutStorage

	// public variables
	Data []float64
}

func NewLinearReg(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinearReg, err error) {
	newVar := LinearReg{}
	newVar.LinearRegWithoutStorage, err = NewLinearRegWithoutStorage(timePeriod, selectData,
		func(dataItem float64, slope float64, intercept float64, streamBarIndex int) {
			newVar.Data = append(newVar.Data, dataItem)
		})

	return &newVar, err
}

func NewLinearRegForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinearReg, err error) {
	newVar, err := NewLinearReg(timePeriod, selectData)
	priceStream.AddTickSubscription(newVar)
	return newVar, err
}

func (ind *LinearRegWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
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

		if result > ind.maxValue {
			ind.maxValue = result
		}

		if result < ind.minValue {
			ind.minValue = result
		}
		ind.valueAvailableAction(result, m, b, streamBarIndex)
	}

	ind.periodHistory.PushBack(tickData)

	if ind.periodHistory.Len() >= ind.GetTimePeriod() {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)
	}

}
