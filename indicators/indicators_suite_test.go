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

var csvFeed *feeds.CSVFileFeed

func TestIndicators(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Indicators Suite")
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

func LoadCSVBollingerPriceDataFromFile(fileName string) (results []indicators.BollingerBandEntry, err error) {
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
		results = append(results, indicators.BollingerBandEntry{UpperBand: upperBandValue, MiddleBand: middleBandValue, LowerBand: lowerBandValue})
	}
	return results, nil
}

var _ = BeforeSuite(func() {
	csvFeed = feeds.NewCSVFileFeedWithDOHLCVFormat("../testdata/JSETOPI.2013.data",
		feeds.DashedYearDayMonthDateParserForLocation(time.Local))
})

var _ = AfterSuite(func() {
	csvFeed = nil
})

func GetDataMax(dohlcvArray []float64) float64 {
	max := math.SmallestNonzeroFloat64

	for i := range dohlcvArray {
		if max < dohlcvArray[i] {
			max = dohlcvArray[i]
		}
	}

	return max
}

func GetDataMin(dohlcvArray []float64) float64 {
	min := math.MaxFloat64

	for i := range dohlcvArray {
		if min > dohlcvArray[i] {
			min = dohlcvArray[i]
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

func GetDataMaxBollinger(dohlcvArray []indicators.BollingerBandEntry, selectData indicators.BollingerDataSelectionFunc) float64 {
	max := math.SmallestNonzeroFloat64

	for i := range dohlcvArray {
		var selectedData = selectData(dohlcvArray[i])
		if max < selectedData {
			max = selectedData
		}
	}

	return max
}

func GetDataMinBollinger(dohlcvArray []indicators.BollingerBandEntry, selectData indicators.BollingerDataSelectionFunc) float64 {
	min := math.MaxFloat64

	for i := range dohlcvArray {
		var selectedData = selectData(dohlcvArray[i])
		if min > selectedData {
			min = selectedData
		}
	}

	return min
}
