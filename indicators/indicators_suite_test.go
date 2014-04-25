package indicators_test

import (
	"encoding/csv"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"os"
	"strconv"
	"testing"
)

var TestInputData []float64

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

var _ = BeforeSuite(func() {
	file, err := os.Open("../testdata/JSETOPI.2013.data")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		closePrice, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		TestInputData = append(TestInputData, closePrice)
	}
})

var _ = AfterSuite(func() {
	TestInputData = nil
})
