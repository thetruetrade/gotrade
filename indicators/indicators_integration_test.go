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

		It("the result set should have a length equal to the source data length less the lookbackperiod", func() {
			Expect(len(sma.Data)).To(Equal(len(priceStream.Data) - sma.GetLookbackPeriod()))
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

		It("the result set should have a length equal to the source data length less the lookbackperiod", func() {
			Expect(len(ema.Data)).To(Equal(len(priceStream.Data) - ema.GetLookbackPeriod()))
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

		It("the result set should have a length equal to the source data length less the lookbackperiod", func() {
			Expect(len(wma.Data)).To(Equal(len(priceStream.Data) - wma.GetLookbackPeriod()))
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

		It("the result set should have a length equal to the source data length less the lookbackperiod", func() {
			Expect(len(dema.Data)).To(Equal(len(priceStream.Data) - (dema.GetLookbackPeriod())))
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

		It("the result set should have a length equal to the source data length less the lookbackperiod", func() {
			Expect(len(tema.Data)).To(Equal(len(priceStream.Data) - (tema.GetLookbackPeriod())))
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

		It("the result set should have a length equal to the source data length less the lookbackperiod", func() {
			Expect(len(variance.Data)).To(Equal(len(priceStream.Data) - variance.GetLookbackPeriod()))
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

		It("the result set should have a length equal to the source data length less the lookbackperiod", func() {
			Expect(len(stdDev.Data)).To(Equal(len(priceStream.Data) - stdDev.GetLookbackPeriod()))
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

		It("the result set should have a length equal to the source data length less the lookbackperiod", func() {
			Expect(bb.Length()).To(Equal(len(priceStream.Data) - bb.GetLookbackPeriod()))
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

		It("the result set should have a length equal to the source data length less the lookbackperiod", func() {
			Expect(macd.Length()).To(Equal(len(priceStream.Data) - macd.GetLookbackPeriod()))
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

		It("the result set should have a length equal to the source data length less the lookbackperiod", func() {
			Expect(aroon.Length()).To(Equal(len(priceStream.Data) - aroon.GetLookbackPeriod()))
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

		It("the result set should have a length equal to the source data length less the lookbackperiod", func() {
			Expect(aroon.Length()).To(Equal(len(priceStream.Data) - aroon.GetLookbackPeriod()))
		})

		It("it should have correctly calculated the aroon oscillator for each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k]).To(BeNumerically("~", aroon.Data[k], 0.01))
			}
		})
	})
})

var _ = Describe("when executing the gotrade truerange with a years data and known output", func() {
	var (
		trueRange       *indicators.TrueRange
		expectedResults []float64
		err             error
		priceStream     *gotrade.DOHLCVStream
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVPriceDataFromFile("truerange_expectedresult.data")
		priceStream = gotrade.NewDOHLCVStream()
	})

	Describe("using an implicit lookback period of 1", func() {

		BeforeEach(func() {
			trueRange, err = indicators.NewTrueRange()
			priceStream.AddTickSubscription(trueRange)
			csvFeed.FillDOHLCVStream(priceStream)
		})

		It("the result set should have a length equal to the source data length less the lookback period", func() {
			Expect(trueRange.Length()).To(Equal(len(priceStream.Data) - trueRange.GetLookbackPeriod()))
		})

		It("it should have correctly calculated the truerangefor each item in the result set accurate to two decimal places", func() {
			for k := range expectedResults {
				Expect(expectedResults[k]).To(BeNumerically("~", trueRange.Data[k], 0.01))
			}
		})
	})

	var _ = Describe("when executing the gotrade average truerange with a years data and known output", func() {
		var (
			avgTrueRange    *indicators.ATR
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("atr_14_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using an implicit lookback period of 14", func() {

			BeforeEach(func() {
				avgTrueRange, err = indicators.NewATR(14)
				priceStream.AddTickSubscription(avgTrueRange)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length less the lookback period", func() {
				Expect(avgTrueRange.Length()).To(Equal(len(priceStream.Data) - avgTrueRange.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the truerangefor each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", avgTrueRange.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade accumulation distribution line with a years data and known output", func() {
		var (
			adl             *indicators.ADL
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("adl_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using no lookback", func() {

			BeforeEach(func() {
				adl, err = indicators.NewADL()
				priceStream.AddTickSubscription(adl)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(adl.Length()).To(Equal(len(priceStream.Data)))
			})

			It("it should have correctly calculated the truerangefor each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", adl.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade chaikin oscilator with a years data and known output", func() {
		var (
			chaikinOsc      *indicators.ChainkinOsc
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
			fastPeriod      int
			slowPeriod      int
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("chaikinosc_3_10_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
			fastPeriod = 3
			slowPeriod = 10
		})

		Describe("using no a fast Time Period of 3 and a slow Time Period of 10", func() {

			BeforeEach(func() {
				chaikinOsc, err = indicators.NewChainkinOsc(3, 10)
				priceStream.AddTickSubscription(chaikinOsc)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length - the lookback Period", func() {
				Expect(chaikinOsc.Length()).To(Equal(len(priceStream.Data) - chaikinOsc.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the chaikin oscillator for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", chaikinOsc.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade on balance volume indicator with a years data and known output", func() {
		var (
			obv             *indicators.OBV
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("obv_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using no lookback period", func() {

			BeforeEach(func() {
				obv, err = indicators.NewOBV()
				priceStream.AddTickSubscription(obv)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(obv.Length()).To(Equal(len(priceStream.Data)))
			})

			It("it should have correctly calculated the on balance volume for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", obv.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade average price indicator with a years data and known output", func() {
		var (
			avgPrice        *indicators.AvgPrice
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("avgprice_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using no lookback period", func() {

			BeforeEach(func() {
				avgPrice, err = indicators.NewAvgPrice()
				priceStream.AddTickSubscription(avgPrice)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(avgPrice.Length()).To(Equal(len(priceStream.Data)))
			})

			It("it should have correctly calculated the avg price for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", avgPrice.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade median price indicator with a years data and known output", func() {
		var (
			medPrice        *indicators.MedPrice
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("medprice_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using no lookback period", func() {

			BeforeEach(func() {
				medPrice, err = indicators.NewMedPrice()
				priceStream.AddTickSubscription(medPrice)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(medPrice.Length()).To(Equal(len(priceStream.Data)))
			})

			It("it should have correctly calculated the median price for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", medPrice.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade typical price indicator with a years data and known output", func() {
		var (
			typPrice        *indicators.TypicalPrice
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("typprice_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using no lookback period", func() {

			BeforeEach(func() {
				typPrice, err = indicators.NewTypicalPrice()
				priceStream.AddTickSubscription(typPrice)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(typPrice.Length()).To(Equal(len(priceStream.Data)))
			})

			It("it should have correctly calculated the typical price for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", typPrice.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade plus directional movement indicator (1) with a years data and known output", func() {
		var (
			plusDM          *indicators.PlusDM
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("plusdm_1_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a time period of 1", func() {

			BeforeEach(func() {
				plusDM, err = indicators.NewPlusDM(1)
				priceStream.AddTickSubscription(plusDM)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(plusDM.Length()).To(Equal(len(priceStream.Data) - plusDM.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the plus directional movement for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", plusDM.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade plus directional movement indicator (14) with a years data and known output", func() {
		var (
			plusDM          *indicators.PlusDM
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("plusdm_14_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a time period of 14", func() {

			BeforeEach(func() {
				plusDM, err = indicators.NewPlusDM(14)
				priceStream.AddTickSubscription(plusDM)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(plusDM.Length()).To(Equal(len(priceStream.Data) - plusDM.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the plus directional movement for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", plusDM.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade minus directional movement indicator (1) with a years data and known output", func() {
		var (
			minusDM         *indicators.MinusDM
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("minusdm_1_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a time period of 1", func() {

			BeforeEach(func() {
				minusDM, err = indicators.NewMinusDM(1)
				priceStream.AddTickSubscription(minusDM)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(minusDM.Length()).To(Equal(len(priceStream.Data) - minusDM.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the minus directional movement for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", minusDM.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade minus directional movement indicator (14) with a years data and known output", func() {
		var (
			minusDM         *indicators.MinusDM
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("minusdm_14_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a time period of 14", func() {

			BeforeEach(func() {
				minusDM, err = indicators.NewMinusDM(14)
				priceStream.AddTickSubscription(minusDM)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(minusDM.Length()).To(Equal(len(priceStream.Data) - minusDM.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the minus directional movement for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", minusDM.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade plus directional indicator (1) with a years data and known output", func() {
		var (
			plusDI          *indicators.PlusDI
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("plusdi_1_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a time period of 1", func() {

			BeforeEach(func() {
				plusDI, err = indicators.NewPlusDI(1)
				priceStream.AddTickSubscription(plusDI)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(plusDI.Length()).To(Equal(len(priceStream.Data) - plusDI.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the plus directional movement for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", plusDI.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade plus directional indicator (14) with a years data and known output", func() {
		var (
			plusDI          *indicators.PlusDI
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("plusdi_14_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a time period of 14", func() {

			BeforeEach(func() {
				plusDI, err = indicators.NewPlusDI(14)
				priceStream.AddTickSubscription(plusDI)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(plusDI.Length()).To(Equal(len(priceStream.Data) - plusDI.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the plus directional movement for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", plusDI.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade minus directional indicator (1) with a years data and known output", func() {
		var (
			minusDI         *indicators.MinusDI
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("minusdi_1_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a time period of 1", func() {

			BeforeEach(func() {
				minusDI, err = indicators.NewMinusDI(1)
				priceStream.AddTickSubscription(minusDI)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(minusDI.Length()).To(Equal(len(priceStream.Data) - minusDI.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the minus directional movement for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", minusDI.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade minus directional indicator (14) with a years data and known output", func() {
		var (
			minusDI         *indicators.MinusDI
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("minusdi_14_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a time period of 14", func() {

			BeforeEach(func() {
				minusDI, err = indicators.NewMinusDI(14)
				priceStream.AddTickSubscription(minusDI)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(minusDI.Length()).To(Equal(len(priceStream.Data) - minusDI.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the minus directional movement for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", minusDI.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade directional movement indicator (14) with a years data and known output", func() {
		var (
			dx              *indicators.DX
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("dx_14_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a time period of 14", func() {

			BeforeEach(func() {
				dx, err = indicators.NewDX(14)
				priceStream.AddTickSubscription(dx)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(dx.Length()).To(Equal(len(priceStream.Data) - dx.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the minus directional movement for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", dx.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade average directional movement indicator (14) with a years data and known output", func() {
		var (
			adx             *indicators.ADX
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("adx_14_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a time period of 14", func() {

			BeforeEach(func() {
				adx, err = indicators.NewADX(14)
				priceStream.AddTickSubscription(adx)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(adx.Length()).To(Equal(len(priceStream.Data) - adx.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the average directional movement for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", adx.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade average directional movement rating (14) with a years data and known output", func() {
		var (
			adxr            *indicators.ADXR
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("adxr_14_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a time period of 14", func() {

			BeforeEach(func() {
				adxr, err = indicators.NewADXR(14)
				priceStream.AddTickSubscription(adxr)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(adxr.Length()).To(Equal(len(priceStream.Data) - adxr.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the average directional rating for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", adxr.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade average directional movement rating (1) with a years data and known output", func() {
		var (
			adxr            *indicators.ADXR
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("adxr_1_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a time period of 1", func() {

			BeforeEach(func() {
				adxr, err = indicators.NewADXR(1)
				priceStream.AddTickSubscription(adxr)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length", func() {
				Expect(adxr.Length()).To(Equal(len(priceStream.Data) - adxr.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the average directional rating for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", adxr.Data[k], 0.01))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade relative strength index with a years data and known output", func() {
		var (
			rsi             *indicators.RSI
			period          int
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("rsi_14_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a lookback period of 14", func() {

			BeforeEach(func() {
				period = 14
				rsi, err = indicators.NewRSI(period, gotrade.UseClosePrice)
				priceStream.AddTickSubscription(rsi)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length less the lookbackperiod", func() {
				Expect(len(rsi.Data)).To(Equal(len(priceStream.Data) - rsi.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the standard deviation for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", rsi.Data[k], 0.1))
				}
			})
		})
	})

	var _ = Describe("when executing the gotrade rate of change with a years data and known output", func() {
		var (
			ind             *indicators.ROC
			period          int
			expectedResults []float64
			err             error
			priceStream     *gotrade.DOHLCVStream
		)

		BeforeEach(func() {
			// load the expected results data
			expectedResults, _ = LoadCSVPriceDataFromFile("roc_10_expectedresult.data")
			priceStream = gotrade.NewDOHLCVStream()
		})

		Describe("using a lookback period of 10", func() {

			BeforeEach(func() {
				period = 10
				ind, err = indicators.NewROC(period, gotrade.UseClosePrice)
				priceStream.AddTickSubscription(ind)
				csvFeed.FillDOHLCVStream(priceStream)
			})

			It("the result set should have a length equal to the source data length less the lookbackperiod", func() {
				Expect(len(ind.Data)).To(Equal(len(priceStream.Data) - ind.GetLookbackPeriod()))
			})

			It("it should have correctly calculated the standard deviation for each item in the result set accurate to two decimal places", func() {
				for k := range expectedResults {
					Expect(expectedResults[k]).To(BeNumerically("~", ind.Data[k], 0.1))
				}
			})
		})
	})
})
