package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	. "github.com/thetruetrade/gotrade/indicators"
	"time"
)

var _ = Describe("when calculating a double exponential moving average (dema)", func() {
	var (
		period     int = 3
		dema       *DEMA
		indicator  Indicator
		sourceData = []gotrade.DOHLCV{gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 5.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 6.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 7.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 8.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 9.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 10.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 11.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 12.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 13.0, 0.0)}
	)

	BeforeEach(func() {
		dema, _ = NewDEMA(period, gotrade.UseClosePrice)
		indicator = dema
	})

	Context("and the dema has received less ticks than twice the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i < period*2-2; i++ {
				dema.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the dema should have no result data", func() {
			Expect(len(dema.Data)).To(Equal(0))
		})

		It("the indicator stream length should be zero", func() {
			Expect(dema.Length()).To(Equal(0))
		})
	})

	Context("and the dema has received ticks equal to twice the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i <= 2*period-2; i++ {
				dema.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the dema should have result data with a single entry", func() {
			Expect(len(dema.Data)).To(Equal(1))
		})

		It("the indicator stream length should be one", func() {
			Expect(dema.Length()).To(Equal(1))
		})

		It("the indicator stream min and max should be equal", func() {
			Expect(dema.MaxValue()).To(Equal(dema.MinValue()))
		})

		It("the indicator stream should have valid bars from twice the lookback period", func() {
			Expect(indicator.ValidFromBar()).To(Equal(2*period - 1))
		})
	})

	Context("and the dema has received more ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := range sourceData {
				dema.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the dema should have result data with entries equal to the number of ticks less twice the (lookback period - 2)", func() {
			Expect(len(dema.Data)).To(Equal(len(sourceData) - (2*period - 2)))
		})

		It("the indicator stream min should equal the result data minimum", func() {
			Expect(indicator.MinValue()).To(Equal(GetDataMin(dema.Data)))
		})

		It("the indicator stream max should equal the result data maximum", func() {
			Expect(indicator.MaxValue()).To(Equal(GetDataMax(dema.Data)))
		})
	})
})
