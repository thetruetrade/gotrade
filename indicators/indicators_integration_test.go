package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	"github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when executing the gotrade simple moving average with a years data and known output", func() {
	var (
		sma             *indicators.SMA
		period          int
		expectedResults []float64
		err             error
		priceStream     *gotrade.DOHLCVStream
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVPriceDataFromFile("sma_10_expectedresult.data")
		priceStream = gotrade.NewDOHLCVStream()
	})

	Describe("using a lookback period of 10", func() {

		BeforeEach(func() {
			period = 10
			sma, err = indicators.NewSMA(period, gotrade.UseClosePrice)
			priceStream.AddSubscription(sma)
			csvFeed.FillDOHLCVStream(priceStream)
		})

		It("the result set should have a length equal to the source data length less the period + 1", func() {
			Expect(len(sma.Data)).To(Equal(len(priceStream.Data) - sma.LookbackPeriod + 1))
		})

		It("it should have correctly calculated the simple moving average for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k]).To(BeNumerically("~", sma.Data[k], 0.01))
			}
		})
	})
})

var _ = Describe("when executing the gotrade exponential moving average with a years data and known output", func() {
	var (
		ema             *indicators.EMA
		period          int
		expectedResults []float64
		err             error
		priceStream     *gotrade.DOHLCVStream
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVPriceDataFromFile("ema_10_expectedresult.data")
		priceStream = gotrade.NewDOHLCVStream()
	})

	Describe("using a lookback period of 10", func() {

		BeforeEach(func() {
			period = 10
			ema, err = indicators.NewEMA(period, gotrade.UseClosePrice)
			priceStream.AddSubscription(ema)
			csvFeed.FillDOHLCVStream(priceStream)
		})

		It("the result set should have a length equal to the source data length less the period + 1", func() {
			Expect(len(ema.Data)).To(Equal(len(priceStream.Data) - ema.LookbackPeriod + 1))
		})

		It("it should have correctly calculated the exponential moving average for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k]).To(BeNumerically("~", ema.Data[k], 0.01))
			}
		})
	})
})

var _ = Describe("when executing the gotrade bollinger bands with a years data and known output", func() {
	var (
		bb              *indicators.BollingerBands
		period          int
		expectedResults []indicators.BollingerBandEntry
		err             error
		priceStream     *gotrade.DOHLCVStream
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVBollingerPriceDataFromFile("bb_10_expectedresult.data")
		priceStream = gotrade.NewDOHLCVStream()
	})

	Describe("using a lookback period of 10", func() {

		BeforeEach(func() {
			period = 10
			bb, err = indicators.NewBollingerBands(period, gotrade.UseClosePrice)
			priceStream.AddSubscription(bb)
			csvFeed.FillDOHLCVStream(priceStream)
		})

		It("the result set should have a length equal to the source data length less the period + 1", func() {
			Expect(len(bb.Data)).To(Equal(len(priceStream.Data) - bb.LookbackPeriod + 1))
		})

		It("it should have correctly calculated the bollinger upper, middle and lower bands for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k].UpperBand).To(BeNumerically("~", bb.Data[k].UpperBand, 0.01))
				Expect(expectedResults[k].MiddleBand).To(BeNumerically("~", bb.Data[k].MiddleBand, 0.01))
				Expect(expectedResults[k].LowerBand).To(BeNumerically("~", bb.Data[k].LowerBand, 0.01))
			}
		})
	})
})
