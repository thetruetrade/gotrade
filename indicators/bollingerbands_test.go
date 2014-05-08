package indicators_test

import (
	//"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	. "github.com/thetruetrade/gotrade/indicators"
	//"math"
	"time"
)

var _ = Describe("when calculating bollinger bands", func() {
	var (
		period     int = 3
		bb         *BollingerBands
		indicator  Indicator
		sourceData = []gotrade.DOHLCV{gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 5.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 6.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 7.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 8.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 9.0, 0.0)}
	)

	BeforeEach(func() {
		bb, _ = NewBollingerBands(period, gotrade.UseClosePrice)
		indicator = bb
	})

	Context("when the bollinger band has received less ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i < period-1; i++ {
				bb.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the bollinger band should have no result data", func() {
			Expect(len(bb.Data)).To(Equal(0))
		})

		It("the indicator stream length should be zero", func() {
			Expect(indicator.Length()).To(Equal(0))
		})

		It("the indicator stream should have no valid bars", func() {
			Expect(indicator.ValidFromBar()).To(Equal(-1))
		})
	})

	Context("when the bollinger band has received ticks equal to the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i <= period-1; i++ {
				bb.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the bollinger band should have result data with a single entry", func() {
			Expect(len(bb.Data)).To(Equal(1))
		})

		It("the indicator stream length should be one", func() {
			Expect(indicator.Length()).To(Equal(1))
		})

		It("the indicator stream should have valid bars from the lookback period", func() {
			Expect(indicator.ValidFromBar()).To(Equal(period))
		})
	})

	Context("when the bollinger band has received more ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := range sourceData {
				bb.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the bollinger band should have result data with entries equal to the number of ticks less the (lookback period - 1)", func() {
			Expect(len(bb.Data)).To(Equal(len(sourceData) - (period - 1)))
		})

		It("the indicator stream min should equal the result data minimum", func() {
			Expect(indicator.MinValue()).To(Equal(GetDataMinBollinger(bb.Data, UseLowerBand)))
		})

		It("the indicator stream max should equal the result data maximum", func() {
			Expect(indicator.MaxValue()).To(Equal(GetDataMaxBollinger(bb.Data, UseUpperBand)))
		})
	})
})
