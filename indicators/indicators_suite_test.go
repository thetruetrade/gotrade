package indicators_test

import (
	"encoding/csv"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	"github.com/thetruetrade/gotrade/feeds"
	"github.com/thetruetrade/gotrade/indicators"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	csvFeed          *feeds.CSVFileFeed
	sourceData       []float64        = []float64{5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0, 17.0, 18.0, 19.0, 20.0}
	sourceDOHLCVData []gotrade.DOHLCV = []gotrade.DOHLCV{gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 5.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 7.0, 8.0, 5.0, 6.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 8.0, 9.0, 6.0, 7.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 9.0, 10.0, 7.0, 8.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 10.0, 11.0, 8.0, 9.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 11.0, 12.0, 9.0, 10.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 12.0, 13.0, 10.0, 11.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 13.0, 14.0, 11.0, 12.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 14.0, 15.0, 12.0, 13.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 15.0, 16.0, 13.0, 14.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 16.0, 17.0, 14.0, 15.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 17.0, 18.0, 15.0, 16.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 18.0, 19.0, 16.0, 17.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 19.0, 20.0, 17.0, 18.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 20.0, 21.0, 18.0, 19.0, 0.0),
		gotrade.NewDOHLCVDataItem(time.Now(), 21.0, 22.0, 19.0, 20.0, 0.0)}
)

func TestIndicators(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Indicators Suite")
}

var _ = BeforeSuite(func() {
	csvFeed = feeds.NewCSVFileFeedWithDOHLCVFormat("../testdata/JSETOPI.2013.data",
		feeds.DashedYearDayMonthDateParserForLocation(time.Local))
})

var _ = AfterSuite(func() {
	csvFeed = nil
})

type IndicatorSharedSpec interface {
	GetIndicator() indicators.Indicator
	GetLength() int
}

type IndicatorWithFloatBoundsSharedSpec interface {
	GetIndicatorWithFloatBounds() indicators.IndicatorWithFloatBounds
	GetIndicator() indicators.Indicator
	GetLength() int
	GetMaximum() float64
	GetMinimum() float64
}

type IndicatorWithIntBoundsSharedSpec interface {
	GetIndicatorWithIntBounds() indicators.IndicatorWithIntBounds
	GetIndicator() indicators.Indicator
	GetLength() int
	GetMaximum() int64
	GetMinimum() int64
}

type IndicatorSharedSpecInputs struct {
	indicatorUnderTest indicators.Indicator
	sourceDataLength   int
}

func NewIndicatorSharedSpecInputs(indicatorUnderTest indicators.Indicator, sourceDataLength int) IndicatorSharedSpecInputs {
	ind := IndicatorSharedSpecInputs{indicatorUnderTest: indicatorUnderTest, sourceDataLength: sourceDataLength}
	return ind
}

func (spec IndicatorSharedSpecInputs) GetIndicator() indicators.Indicator {
	return spec.indicatorUnderTest
}

func (spec IndicatorSharedSpecInputs) GetLength() int {
	return spec.sourceDataLength
}

type IndicatorWithFloatBoundsSharedSpecInputs struct {
	*IndicatorSharedSpecInputs
	indicatorWithFloatBoundsUnderTest indicators.IndicatorWithFloatBounds
	getMaximum                        GetMaximumFloatFunc
	getMinimum                        GetMaximumFloatFunc
}

func NewIndicatorWithFloatBoundsSharedSpecInputs(indicatorUnderTest indicators.Indicator,
	sourceDataLength int,
	indicatorWithFloatBoundsUnderTest indicators.IndicatorWithFloatBounds,
	getMaximum GetMaximumFloatFunc,
	getMinimum GetMaximumFloatFunc) IndicatorWithFloatBoundsSharedSpecInputs {
	ind := IndicatorWithFloatBoundsSharedSpecInputs{IndicatorSharedSpecInputs: &IndicatorSharedSpecInputs{indicatorUnderTest: indicatorUnderTest, sourceDataLength: sourceDataLength},
		indicatorWithFloatBoundsUnderTest: indicatorWithFloatBoundsUnderTest,
		getMaximum:                        getMaximum,
		getMinimum:                        getMinimum}
	return ind
}

func (spec IndicatorWithFloatBoundsSharedSpecInputs) GetIndicatorWithFloatBounds() indicators.IndicatorWithFloatBounds {
	return spec.indicatorWithFloatBoundsUnderTest
}

func (spec IndicatorWithFloatBoundsSharedSpecInputs) GetMaximum() float64 {
	return spec.getMaximum()

}
func (spec IndicatorWithFloatBoundsSharedSpecInputs) GetMinimum() float64 {
	return spec.getMinimum()
}

type IndicatorWithIntBoundsSharedSpecInputs struct {
	*IndicatorSharedSpecInputs
	indicatorWithIntBoundsUnderTest indicators.IndicatorWithIntBounds
	getMaximum                      GetMaximumIntFunc
	getMinimum                      GetMaximumIntFunc
}

func NewIndicatorWithIntBoundsSharedSpecInputs(indicatorUnderTest indicators.Indicator,
	sourceDataLength int,
	indicatorWithIntBoundsUnderTest indicators.IndicatorWithIntBounds,
	getMaximum GetMaximumIntFunc,
	getMinimum GetMaximumIntFunc) IndicatorWithIntBoundsSharedSpecInputs {
	ind := IndicatorWithIntBoundsSharedSpecInputs{IndicatorSharedSpecInputs: &IndicatorSharedSpecInputs{indicatorUnderTest: indicatorUnderTest, sourceDataLength: sourceDataLength},
		indicatorWithIntBoundsUnderTest: indicatorWithIntBoundsUnderTest,
		getMaximum:                      getMaximum,
		getMinimum:                      getMinimum}
	return ind
}

func (spec IndicatorWithIntBoundsSharedSpecInputs) GetIndicatorWithIntBounds() indicators.IndicatorWithIntBounds {
	return spec.indicatorWithIntBoundsUnderTest
}

func (spec IndicatorWithIntBoundsSharedSpecInputs) GetMaximum() int64 {
	return spec.getMaximum()

}
func (spec IndicatorWithIntBoundsSharedSpecInputs) GetMinimum() int64 {
	return spec.getMinimum()
}

func ShouldBeAnIndicatorThatHasReceivedAllOfItsTicks(inputs IndicatorSharedSpec) {
	It("the indicator should be valid from some bar >= 1", func() {
		Expect(inputs.GetIndicator().ValidFromBar()).To(BeNumerically(">=", 1))
	})

	It("the indicator stream should have entries equal to the number of ticks less the lookback period", func() {
		Expect(inputs.GetIndicator().Length()).To(Equal(inputs.GetLength() - inputs.GetIndicator().GetLookbackPeriod()))
	})
}

func ShouldHaveFloatBoundsSetToMinMaxOfResults(inputs IndicatorWithFloatBoundsSharedSpec) {
	It("the indicator min should equal the result stream minimum", func() {
		Expect(inputs.GetIndicatorWithFloatBounds().MinValue()).To(Equal(inputs.GetMinimum()))
	})

	It("the indicator max should equal the result stream maximum", func() {
		Expect(inputs.GetIndicatorWithFloatBounds().MaxValue()).To(Equal(inputs.GetMaximum()))
	})
}

func ShouldHaveIntBoundsSetToMinMaxOfResults(inputs IndicatorWithIntBoundsSharedSpec) {
	It("the indicator min should equal the result stream minimum", func() {
		Expect(inputs.GetIndicatorWithIntBounds().MinValue()).To(Equal(inputs.GetMinimum()))
	})

	It("the indicator max should equal the result stream maximum", func() {
		Expect(inputs.GetIndicatorWithIntBounds().MaxValue()).To(Equal(inputs.GetMaximum()))
	})
}

func ShouldBeAnInitialisedIndicator(inputs IndicatorSharedSpec) {

	It("the indicator should not be valid from any bar yet", func() {
		Expect(inputs.GetIndicator().ValidFromBar()).To(Equal(-1))
	})

	It("the indicator stream should have no results", func() {
		Expect(inputs.GetIndicator().Length()).To(BeZero())
	})

	It("the indicator should have a valid lookback period", func() {
		Expect(inputs.GetIndicator().GetLookbackPeriod()).Should(BeNumerically(">=", indicators.MinimumLookbackPeriod))
		Expect(inputs.GetIndicator().GetLookbackPeriod()).Should(BeNumerically("<=", indicators.MaximumLookbackPeriod))
	})
	// nobounds
}

func ShouldNotHaveAnyFloatBoundsSetYet(inputs IndicatorWithFloatBoundsSharedSpec) {
	It("the indicator should have no minimum value set", func() {
		Expect(inputs.GetIndicatorWithFloatBounds().MinValue()).To(Equal(math.MaxFloat64))
	})

	It("the indicator should have no maximum value set", func() {
		Expect(inputs.GetIndicatorWithFloatBounds().MaxValue()).To(Equal(math.SmallestNonzeroFloat64))
	})
}

func ShouldNotHaveAnyIntBoundsSetYet(inputs IndicatorWithIntBoundsSharedSpec) {
	It("the indicator should have no minimum value set", func() {
		Expect(inputs.GetIndicatorWithIntBounds().MinValue()).To(Equal(int64(math.MaxInt64)))
	})

	It("the indicator should have no maximum value set", func() {
		Expect(inputs.GetIndicatorWithIntBounds().MaxValue()).To(Equal(int64(math.MinInt64)))
	})
}

func ShouldBeAnIndicatorThatHasReceivedFewerTicksThanItsLookbackPeriod(inputs IndicatorSharedSpec) {

	It("the indicator should not be valid from any bar yet", func() {
		Expect(inputs.GetIndicator().ValidFromBar()).To(Equal(-1))
	})

	It("the indicator stream should have no results", func() {
		Expect(inputs.GetIndicator().Length()).To(BeZero())
	})
}

func ShouldBeAnIndicatorThatHasReceivedTicksEqualToItsLookbackPeriod(inputs IndicatorSharedSpec) {
	It("the indicator stream should have a single entry", func() {
		Expect(inputs.GetIndicator().Length()).To(Equal(1))
	})

	It("the indicator should be valid from the lookback period", func() {
		Expect(inputs.GetIndicator().ValidFromBar()).To(Equal(inputs.GetIndicator().GetLookbackPeriod() + 1))
	})
}

func ShouldBeAnIndicatorThatHasReceivedMoreTicksThanItsLookbackPeriod(inputs IndicatorSharedSpec) {
	It("the indicator stream should have entries equal to the number of ticks less the lookback period", func() {
		Expect(inputs.GetIndicator().Length()).To(Equal(inputs.GetLength() - (inputs.GetIndicator().GetLookbackPeriod())))
	})
}

func LoadCSVPriceDataFromFile(fileName string) (results []float64, err error) {
	file, err := os.Open("../testdata/" + fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		priceValue, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		results = append(results, priceValue)
	}
	return results, nil
}

func LoadCSVIntPriceDataFromFile(fileName string) (results []int64, err error) {
	file, err := os.Open("../testdata/" + fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		tmp, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		priceValue := int64(tmp)
		results = append(results, priceValue)
	}
	return results, nil
}

func LoadCSVBollingerPriceDataFromFile(fileName string) (results []BollingerBand, err error) {
	file, err := os.Open("../testdata/" + fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		upperBandValue, err := strconv.ParseFloat(strings.TrimSpace(record[0]), 64)
		middleBandValue, err := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		lowerBandValue, err := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)

		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		results = append(results, NewBollingerBandDataItem(upperBandValue, middleBandValue, lowerBandValue))
	}
	return results, nil
}

func LoadCSVMACDPriceDataFromFile(fileName string) (results []MACDData, err error) {
	file, err := os.Open("../testdata/" + fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		macd, err := strconv.ParseFloat(strings.TrimSpace(record[0]), 64)
		signal, err := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		histogram, err := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)

		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		results = append(results, NewMACDDataItem(macd, signal, histogram))
	}
	return results, nil
}

func LoadCSVAroonPriceDataFromFile(fileName string) (results []AroonData, err error) {
	file, err := os.Open("../testdata/" + fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		up, err := strconv.ParseFloat(strings.TrimSpace(record[0]), 64)
		dwn, err := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)

		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		results = append(results, NewAroonDataItem(up, dwn))
	}
	return results, nil
}

func LoadCSVStochPriceDataFromFile(fileName string) (results []StochData, err error) {
	file, err := os.Open("../testdata/" + fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		k, err := strconv.ParseFloat(strings.TrimSpace(record[0]), 64)
		d, err := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)

		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		results = append(results, NewStochDataItem(k, d))
	}
	return results, nil
}

type GetMaximumFloatFunc func() float64

type GetMinimumFloatFunc func() float64

type GetMaximumIntFunc func() int64

type GetMinimumIntFunc func() int64

func GetFloatDataMax(floatArray []float64) float64 {
	max := math.SmallestNonzeroFloat64

	for i := range floatArray {
		if max < floatArray[i] {
			max = floatArray[i]
		}
	}

	return max
}

func GetIntDataMax(intArray []int64) int64 {
	max := int64(math.MinInt64)

	for i := range intArray {
		if max < intArray[i] {
			max = intArray[i]
		}
	}

	return max
}

func GetFloatDataMin(floatArray []float64) float64 {
	min := math.MaxFloat64

	for i := range floatArray {
		if min > floatArray[i] {
			min = floatArray[i]
		}
	}

	return min
}

func GetIntDataMin(intArray []int64) int64 {
	min := int64(math.MaxInt64)

	for i := range intArray {
		if min > intArray[i] {
			min = intArray[i]
		}
	}

	return min
}

func GetDataMaxDOHLCV(dohlcvArray []gotrade.DOHLCV, selectData gotrade.DataSelectionFunc) float64 {
	max := math.SmallestNonzeroFloat64

	for i := range dohlcvArray {
		var selectedData = selectData(dohlcvArray[i])
		if max < selectedData {
			max = selectedData
		}
	}

	return max
}

func GetDataMinDOHLCV(dohlcvArray []gotrade.DOHLCV, selectData gotrade.DataSelectionFunc) float64 {
	min := math.MaxFloat64

	for i := range dohlcvArray {
		var selectedData = selectData(dohlcvArray[i])
		if min > selectedData {
			min = selectedData
		}
	}

	return min
}

func GetDataMaxMACD(macd []float64, signal []float64, histogram []float64) float64 {
	max := math.SmallestNonzeroFloat64

	for i := range macd {
		macd := macd[i]
		signal := signal[i]
		histogram := histogram[i]

		if max < macd {
			max = macd
		}
		if max < signal {
			max = signal
		}
		if max < histogram {
			max = histogram
		}
	}

	return max
}

func GetDataMinMACD(macd []float64, signal []float64, histogram []float64) float64 {
	min := math.MaxFloat64

	for i := range macd {
		macd := macd[i]
		signal := signal[i]
		histogram := histogram[i]

		if min > macd {
			min = macd
		}
		if min > signal {
			min = signal
		}
		if min > histogram {
			min = histogram
		}
	}

	return min
}

func GetDataMaxStoch(slowK []float64, slowD []float64) float64 {
	max := math.SmallestNonzeroFloat64

	for i := range slowK {
		slowKVal := slowK[i]
		slowDVal := slowD[i]

		if max < slowKVal {
			max = slowKVal
		}
		if max < slowDVal {
			max = slowDVal
		}
	}

	return max
}

func GetDataMinStoch(slowK []float64, slowD []float64) float64 {
	min := math.MaxFloat64

	for i := range slowK {
		slowKVal := slowK[i]
		slowDVal := slowD[i]

		if min > slowKVal {
			min = slowKVal
		}
		if min > slowDVal {
			min = slowDVal
		}
	}

	return min
}

type MACDData interface {
	// MACD
	M() float64
	// Signal
	S() float64
	// Histogram
	H() float64
}

type MACDDataItem struct {
	macd      float64
	signal    float64
	histogram float64
}

func (data *MACDDataItem) M() float64 {
	return data.macd
}

func (data *MACDDataItem) S() float64 {
	return data.signal
}

func (data *MACDDataItem) H() float64 {
	return data.histogram
}

func NewMACDDataItem(macd float64, signal float64, histogram float64) *MACDDataItem {
	return &MACDDataItem{macd: macd, signal: signal, histogram: histogram}
}

type BollingerBand interface {
	// Upper bollinger band
	U() float64
	// Middle bollinger band
	M() float64
	// Lower bollinger band
	L() float64
}

type BollingerBandDataItem struct {
	upperBand  float64
	middleBand float64
	lowerBand  float64
}

func (bb *BollingerBandDataItem) U() float64 {
	return bb.upperBand
}

func (bb *BollingerBandDataItem) L() float64 {
	return bb.lowerBand
}

func (bb *BollingerBandDataItem) M() float64 {
	return bb.middleBand
}

func NewBollingerBandDataItem(upperBand float64, middleBand float64, lowerBand float64) *BollingerBandDataItem {
	return &BollingerBandDataItem{upperBand: upperBand, middleBand: middleBand, lowerBand: lowerBand}
}

type AroonData interface {
	U() float64
	D() float64
}

type AroonDataItem struct {
	up  float64
	dwn float64
}

func NewAroonDataItem(up float64, dwn float64) *AroonDataItem {
	return &AroonDataItem{up: up, dwn: dwn}
}

func (adi *AroonDataItem) U() float64 {
	return adi.up
}

func (adi *AroonDataItem) D() float64 {
	return adi.dwn
}

type StochData interface {
	K() float64
	D() float64
}

type StochDataItem struct {
	k float64
	d float64
}

func NewStochDataItem(k float64, d float64) *StochDataItem {
	return &StochDataItem{k: k, d: d}
}

func (sdi *StochDataItem) K() float64 {
	return sdi.k
}

func (sdi *StochDataItem) D() float64 {
	return sdi.d
}
