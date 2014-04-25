// Simple Moving Average (SMA)
package indicators

// A Simple Moving Average Indicator
type SMA struct {
	usage          indicatorUsageType
	LookbackPeriod int
}

// NewSMA returns a new Simple Moving Average (SMA) configured with the
// specified lookbackPeriod
func NewSMA(lookbackPeriod int) (indicator *SMA, err error) {
	return &SMA{usage: SubChart, LookbackPeriod: lookbackPeriod}, nil
}

// Calculates the Simple Moving Average (SMA) for the specified source.
// The source data must contain more items than the configured LookbackPeriod.
func (sma *SMA) Calculate(source []float64) (results []float64, err error) {

	// perform some sanity checks on the source data
	err = checkSourceDataIsNotEmpty(source)
	if err != nil {
		return nil, err
	}
	// and again on the source data with regards to the lookback period
	err = checkSourceValidForLookbackPeriod(source, sma.LookbackPeriod)
	if err != nil {
		return nil, err
	}

	// compute local variables for use in the loop
	periodTotal := float64(0.0)
	sourceLength := len(source)
	outputLength := sourceLength - sma.LookbackPeriod + 1

	// initialise the output data array
	results = make([]float64, outputLength)

	// iterate the source data

	y := (sma.LookbackPeriod * -1)
	for i := 0; i < sourceLength; i++ {
		y += 1

		if y > 0 {
			periodTotal -= source[i-sma.LookbackPeriod]
		}
		periodTotal += source[i]
		result := periodTotal / float64(sma.LookbackPeriod)
		if y >= 0 {
			results[y] = result
		}
	}

	return results, nil
}
