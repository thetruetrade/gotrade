package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	. "github.com/thetruetrade/gotrade/indicators"
	"time"
)

var _ = Describe("when calculating a triple exponential moving average (tema)", func() {
	var (
		period     int = 3
		tema       *TEMA
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
		tema, _ = NewTEMA(period, gotrade.UseClosePrice)
		indicator = tema
	})

	Context("and the tema has received less ticks than triple the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i < period*3-3; i++ {
				tema.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the tema should have no result data", func() {
			Expect(len(tema.Data)).To(Equal(0))
		})

		It("the indicator stream length should be zero", func() {
			Expect(tema.Length()).To(Equal(0))
		})
	})

	Context("and the tema has received ticks equal to triple the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i <= 3*period-3; i++ {
				tema.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the tema should have result data with a single entry", func() {
			Expect(len(tema.Data)).To(Equal(1))
		})

		It("the indicator stream length should be one", func() {
			Expect(tema.Length()).To(Equal(1))
		})

		It("the indicator stream min and max should be equal", func() {
			Expect(tema.MaxValue()).To(Equal(tema.MinValue()))
		})

		It("the indicator stream should have valid bars from triple the lookback period", func() {
			Expect(indicator.ValidFromBar()).To(Equal(3*period - 2))
		})
	})

	Context("and the tema has received more ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := range sourceData {
				tema.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the tema should have result data with entries equal to the number of ticks less triple the (lookback period - 3)", func() {
			Expect(len(tema.Data)).To(Equal(len(sourceData) - (3*period - 3)))
		})

		It("the indicator stream min should equal the result data minimum", func() {
			Expect(indicator.MinValue()).To(Equal(GetDataMin(tema.Data)))
		})

		It("the indicator stream max should equal the result data maximum", func() {
			Expect(indicator.MaxValue()).To(Equal(GetDataMax(tema.Data)))
		})
	})
})
