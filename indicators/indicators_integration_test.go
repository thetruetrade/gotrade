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
			Expect(len(sma.Data)).To(Equal(len(priceStream.Data) - sma.GetLookbackPeriod() + 1))
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
			Expect(len(ema.Data)).To(Equal(len(priceStream.Data) - ema.GetLookbackPeriod() + 1))
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
			Expect(len(wma.Data)).To(Equal(len(priceStream.Data) - wma.GetLookbackPeriod() + 1))
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
			Expect(len(dema.Data)).To(Equal(len(priceStream.Data) - (dema.GetLookbackPeriod() - 1)))
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
			Expect(len(tema.Data)).To(Equal(len(priceStream.Data) - (tema.GetLookbackPeriod() - 1)))
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
			Expect(len(variance.Data)).To(Equal(len(priceStream.Data) - variance.GetLookbackPeriod() + 1))
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
			Expect(len(stdDev.Data)).To(Equal(len(priceStream.Data) - stdDev.GetLookbackPeriod() + 1))
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
		expectedResults []BollingerBand
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
			Expect(bb.Length()).To(Equal(len(priceStream.Data) - bb.GetLookbackPeriod() + 1))
		})

		It("it should have correctly calculated the bollinger upper, middle and lower bands for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k].U()).To(BeNumerically("~", bb.UpperBand[k], 0.01))
				Expect(expectedResults[k].M()).To(BeNumerically("~", bb.MiddleBand[k], 0.01))
				Expect(expectedResults[k].L()).To(BeNumerically("~", bb.LowerBand[k], 0.01))
			}
		})
	})
})

var _ = Describe("when executing the gotrade macd with a years data and known output", func() {
	var (
		macd            *indicators.MACD
		expectedResults []MACDData
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
			Expect(macd.Length()).To(Equal(len(priceStream.Data) - (26 + 8) + 1))
		})

		It("it should have correctly calculated the macd, signal and histogram for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k].M()).To(BeNumerically("~", macd.MACD[k], 0.01))
				Expect(expectedResults[k].S()).To(BeNumerically("~", macd.Signal[k], 0.01))
				Expect(expectedResults[k].H()).To(BeNumerically("~", macd.Histogram[k], 0.01))
			}
		})
	})
})

var _ = Describe("when executing the gotrade aroon with a years data and known output", func() {
	var (
		aroon           *indicators.Aroon
		expectedResults []AroonData
		err             error
		priceStream     *gotrade.DOHLCVStream
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVAroonPriceDataFromFile("aroon_25_expectedresult.data")
		priceStream = gotrade.NewDOHLCVStream()
	})

	Describe("using a lookback periods of 25", func() {

		BeforeEach(func() {
			aroon, err = indicators.NewAroon(25)
			priceStream.AddTickSubscription(aroon)
			csvFeed.FillDOHLCVStream(priceStream)
		})

		It("the result set should have a length equal to the source data length less the period + 1", func() {
			Expect(aroon.Length()).To(Equal(len(priceStream.Data) - 25))
		})

		It("it should have correctly calculated the aroon up for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k].U()).To(BeNumerically("~", aroon.Up[k], 0.01))
			}
		})

		It("it should have correctly calculated the aroon down for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k].D()).To(BeNumerically("~", aroon.Down[k], 0.01))
			}
		})
	})
})

var _ = Describe("when executing the gotrade aroonosc with a years data and known output", func() {
	var (
		aroon           *indicators.AroonOsc
		expectedResults []float64
		err             error
		priceStream     *gotrade.DOHLCVStream
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVPriceDataFromFile("aroonosc_25_expectedresult.data")
		priceStream = gotrade.NewDOHLCVStream()
	})

	Describe("using a lookback periods of 25", func() {

		BeforeEach(func() {
			aroon, err = indicators.NewAroonOsc(25)
			priceStream.AddTickSubscription(aroon)
			csvFeed.FillDOHLCVStream(priceStream)
		})

		It("the result set should have a length equal to the source data length less the period + 1", func() {
			Expect(aroon.Length()).To(Equal(len(priceStream.Data) - 25))
		})

		It("it should have correctly calculated the aroon oscillator for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k]).To(BeNumerically("~", aroon.Data[k], 0.01))
			}
		})
	})
})
