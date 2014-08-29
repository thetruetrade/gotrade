package indicators

// Tema(X) = (2 * EMA(X, CLOSE)) - (EMA(X, EMA(X, CLOSE)))

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Tripple Exponential Moving Average Indicator (Tema), no storage, for use in other indicators
type TemaWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	ema1        *EmaWithoutStorage
	ema2        *EmaWithoutStorage
	ema3        *EmaWithoutStorage
	currentEMA  float64
	currentEMA2 float64
	timePeriod  int
}

func NewTemaWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *TemaWithoutStorage, err error) {

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

	lookback := 3 * (timePeriod - 1)
	ind := TemaWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
		timePeriod:                   timePeriod,
	}

	ind.ema1, err = NewEmaWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentEMA = dataItem
		ind.ema2.ReceiveTick(dataItem, streamBarIndex)
	})

	ind.ema2, _ = NewEmaWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentEMA2 = dataItem
		ind.ema3.ReceiveTick(dataItem, streamBarIndex)
	})

	ind.ema3, _ = NewEmaWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {

		//TEMA = (3*EMA – 3*EMA(EMA)) + EMA(EMA(EMA))
		result := (3*ind.currentEMA - 3*ind.currentEMA2) + dataItem

		ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
	})

	return &ind, err
}

// A Tripple Exponential Moving Average Indicator (Tema)
type Tema struct {
	*TemaWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewTema creates a Tripple Exponential Moving Average Indicator (Tema) for online usage
func NewTema(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Tema, err error) {
	ind := Tema{selectData: selectData}
	ind.TemaWithoutStorage, err = NewTemaWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})
	return &ind, err
}

// NewDefaultTema creates a Tripple Exponential Moving Average Indicator (Tema) for online usage with default parameters
//	- timePeriod: 30
func NewDefaultTema() (indicator *Tema, err error) {
	timePeriod := 30
	return NewTema(timePeriod, gotrade.UseClosePrice)
}

// NewTemaWithSrcLen creates a Tripple Exponential Moving Average Indicator (Tema) for offline usage
func NewTemaWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Tema, err error) {
	ind, err := NewTema(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultTemaWithSrcLen creates a Tripple Exponential Moving Average Indicator (Tema) for offline usage with default parameters
func NewDefaultTemaWithSrcLen(sourceLength uint) (indicator *Tema, err error) {
	ind, err := NewDefaultTema()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewTemaForStream creates a Tripple Exponential Moving Average Indicator (Tema) for online usage with a source data stream
func NewTemaForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Tema, err error) {
	ind, err := NewTema(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultTemaForStream creates a Tripple Exponential Moving Average Indicator (Tema) for online usage with a source data stream
func NewDefaultTemaForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Tema, err error) {
	ind, err := NewDefaultTema()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewTemaForStreamWithSrcLen creates a Tripple Exponential Moving Average Indicator (Tema) for offline usage with a source data stream
func NewTemaForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Tema, err error) {
	ind, err := NewTemaWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultTemaForStreamWithSrcLen creates a Tripple Exponential Moving Average Indicator (Tema) for offline usage with a source data stream
func NewDefaultTemaForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Tema, err error) {
	ind, err := NewDefaultTemaWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *Tema) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *TemaWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.ema1.ReceiveTick(tickData, streamBarIndex)
}
