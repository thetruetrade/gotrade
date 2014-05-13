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
			priceStream.AddTickSubscription(sma)
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
			priceStream.AddTickSubscription(ema)
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

var _ = Describe("when executing the gotrade weighted moving average with a years data and known output", func() {
	var (
		wma             *indicators.WMA
		period          int
		expectedResults []float64
		err             error
		priceStream     *gotrade.DOHLCVStream
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVPriceDataFromFile("wma_10_expectedresult.data")
		priceStream = gotrade.NewDOHLCVStream()
	})

	Describe("using a lookback period of 10", func() {

		BeforeEach(func() {
			period = 10
			wma, err = indicators.NewWMA(period, gotrade.UseClosePrice)
			priceStream.AddTickSubscription(wma)
			csvFeed.FillDOHLCVStream(priceStream)
		})

		It("the result set should have a length equal to the source data length less the period + 1", func() {
			Expect(len(wma.Data)).To(Equal(len(priceStream.Data) - wma.LookbackPeriod + 1))
		})

		It("it should have correctly calculated the weighted moving average for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k]).To(BeNumerically("~", wma.Data[k], 0.01))
			}
		})
	})
})

var _ = Describe("when executing the gotrade double exponential moving average with a years data and known output", func() {
	var (
		dema            *indicators.DEMA
		period          int
		expectedResults []float64
		err             error
		priceStream     *gotrade.DOHLCVStream
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVPriceDataFromFile("dema_10_expectedresult.data")
		priceStream = gotrade.NewDOHLCVStream()
	})

	Describe("using a lookback period of 10", func() {

		BeforeEach(func() {
			period = 10
			dema, err = indicators.NewDEMA(period, gotrade.UseClosePrice)
			priceStream.AddTickSubscription(dema)
			csvFeed.FillDOHLCVStream(priceStream)
		})

		It("the result set should have a length equal to the source data length less twice the lookback period -1", func() {
			Expect(len(dema.Data)).To(Equal(len(priceStream.Data) - (dema.LookbackPeriod - 1)))
		})

		It("it should have correctly calculated the double exponential moving average for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k]).To(BeNumerically("~", dema.Data[k], 0.01))
			}
		})
	})
})

var _ = Describe("when executing the gotrade triple exponential moving average with a years data and known output", func() {
	var (
		tema            *indicators.TEMA
		period          int
		expectedResults []float64
		err             error
		priceStream     *gotrade.DOHLCVStream
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVPriceDataFromFile("tema_10_expectedresult.data")
		priceStream = gotrade.NewDOHLCVStream()
	})

	Describe("using a lookback period of 10", func() {

		BeforeEach(func() {
			period = 10
			tema, err = indicators.NewTEMA(period, gotrade.UseClosePrice)
			priceStream.AddTickSubscription(tema)
			csvFeed.FillDOHLCVStream(priceStream)
		})

		It("the result set should have a length equal to the source data length less triple the looback period -1", func() {
			Expect(len(tema.Data)).To(Equal(len(priceStream.Data) - (tema.LookbackPeriod - 1)))
		})

		It("it should have correctly calculated the triple exponential moving average for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k]).To(BeNumerically("~", tema.Data[k], 0.01))
			}
		})
	})
})

var _ = Describe("when executing the gotrade variance with a years data and known output", func() {
	var (
		variance        *indicators.Variance
		period          int
		expectedResults []float64
		err             error
		priceStream     *gotrade.DOHLCVStream
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVPriceDataFromFile("variance_10_expectedresult.data")
		priceStream = gotrade.NewDOHLCVStream()
	})

	Describe("using a lookback period of 10", func() {

		BeforeEach(func() {
			period = 10
			variance, err = indicators.NewVariance(period, gotrade.UseClosePrice)
			priceStream.AddTickSubscription(variance)
			csvFeed.FillDOHLCVStream(priceStream)
		})

		It("the result set should have a length equal to the source data length less the period + 1", func() {
			Expect(len(variance.Data)).To(Equal(len(priceStream.Data) - variance.LookbackPeriod + 1))
		})

		It("it should have correctly calculated the variance for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k]).To(BeNumerically("~", variance.Data[k], 0.1))
			}
		})
	})
})

var _ = Describe("when executing the gotrade standard deviation with a years data and known output", func() {
	var (
		stdDev          *indicators.StdDeviation
		period          int
		expectedResults []float64
		err             error
		priceStream     *gotrade.DOHLCVStream
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVPriceDataFromFile("stddev_10_expectedresult.data")
		priceStream = gotrade.NewDOHLCVStream()
	})

	Describe("using a lookback period of 10", func() {

		BeforeEach(func() {
			period = 10
			stdDev, err = indicators.NewStdDeviation(period, gotrade.UseClosePrice)
			priceStream.AddTickSubscription(stdDev)
			csvFeed.FillDOHLCVStream(priceStream)
		})

		It("the result set should have a length equal to the source data length less the period + 1", func() {
			Expect(len(stdDev.Data)).To(Equal(len(priceStream.Data) - stdDev.LookbackPeriod + 1))
		})

		It("it should have correctly calculated the standard deviation for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k]).To(BeNumerically("~", stdDev.Data[k], 0.1))
			}
		})
	})
})

var _ = Describe("when executing the gotrade bollinger bands with a years data and known output", func() {
	var (
		bb              *indicators.BollingerBands
		period          int
		expectedResults []indicators.BollingerBand
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
			priceStream.AddTickSubscription(bb)
			csvFeed.FillDOHLCVStream(priceStream)
		})

		It("the result set should have a length equal to the source data length less the period + 1", func() {
			Expect(len(bb.Data)).To(Equal(len(priceStream.Data) - bb.LookbackPeriod + 1))
		})

		It("it should have correctly calculated the bollinger upper, middle and lower bands for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k].U()).To(BeNumerically("~", bb.Data[k].U(), 0.01))
				Expect(expectedResults[k].M()).To(BeNumerically("~", bb.Data[k].M(), 0.01))
				Expect(expectedResults[k].L()).To(BeNumerically("~", bb.Data[k].L(), 0.01))
			}
		})
	})
})

var _ = Describe("when executing the gotrade macd with a years data and known output", func() {
	var (
		macd            *indicators.MACD
		expectedResults []indicators.MACDData
		err             error
		priceStream     *gotrade.DOHLCVStream
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVMACDPriceDataFromFile("macd_12_26_9_expectedresult.data")
		priceStream = gotrade.NewDOHLCVStream()
	})

	Describe("using a lookback periods of 12, 26, 9", func() {

		BeforeEach(func() {
			macd, err = indicators.NewMACD(12, 26, 9, gotrade.UseClosePrice)
			priceStream.AddTickSubscription(macd)
			csvFeed.FillDOHLCVStream(priceStream)
		})

		It("the result set should have a length equal to the source data length less the period + 1", func() {
			Expect(len(macd.Data)).To(Equal(len(priceStream.Data) - (26 + 8) + 1))
		})

		It("it should have correctly calculated the macd, signal and histogram for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k].M()).To(BeNumerically("~", macd.Data[k].M(), 0.01))
				Expect(expectedResults[k].S()).To(BeNumerically("~", macd.Data[k].S(), 0.01))
				Expect(expectedResults[k].H()).To(BeNumerically("~", macd.Data[k].H(), 0.01))
			}
		})
	})
})
