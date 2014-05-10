package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	. "github.com/thetruetrade/gotrade/indicators"
	"time"
)

var _ = Describe("when calculating a moving average convergence divergence (macd)", func() {
	var (
		shortPeriod  int = 3
		longPeriod   int = 6
		signalPeriod int = 2
		macd         *MACD
		indicator    Indicator
		sourceData   = []gotrade.DOHLCV{gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 5.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 6.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 7.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 8.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 9.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 9.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 10.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 11.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 12.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 13.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 14.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 15.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 16.0, 0.0)}
	)

	BeforeEach(func() {
		macd, _ = NewMACD(shortPeriod, longPeriod, signalPeriod, gotrade.UseClosePrice)
		indicator = macd
	})

	Context("and the macd has received less ticks than the slowEMA period", func() {

		BeforeEach(func() {
			for i := 0; i < longPeriod-1; i++ {
				macd.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the macd should have no result data", func() {
			Expect(len(macd.Data)).To(Equal(0))
		})

		It("the indicator stream length should be zero", func() {
			Expect(indicator.Length()).To(Equal(0))
		})

		It("the indicator stream should have no valid bars", func() {
			Expect(indicator.ValidFromBar()).To(Equal(-1))
		})
	})

	Context("and the macd has received ticks equal to the slowEMA period + the signalEMA period", func() {

		BeforeEach(func() {
			for i := 0; i <= longPeriod+signalPeriod-2; i++ {
				macd.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the macd should have result data with a single entry", func() {
			Expect(len(macd.Data)).To(Equal(1))
		})

		It("the indicator stream length should be one", func() {
			Expect(indicator.Length()).To(Equal(1))
		})

		It("the indicator stream should have valid bars from the slowEMA period + signalEMA period - 1", func() {
			Expect(indicator.ValidFromBar()).To(Equal(longPeriod + signalPeriod - 1))
		})
	})

	Context("and the macd has received more ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := range sourceData {
				macd.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the macd should have result data with entries equal to the number of ticks less the (slowEMA period + signalEMA -1) - 1 ", func() {
			Expect(len(macd.Data)).To(Equal(len(sourceData) - (longPeriod + signalPeriod - 1 - 1)))
		})

		It("the indicator stream min should equal the result data minimum", func() {
			Expect(indicator.MinValue()).To(Equal(GetDataMinMACD(macd.Data)))
		})

		It("the indicator stream max should equal the result data maximum", func() {
			Expect(indicator.MaxValue()).To(Equal(GetDataMaxMACD(macd.Data)))
		})
	})
})
