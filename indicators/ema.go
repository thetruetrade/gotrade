// Exponential Moving Average (EMA)
package indicators

// A Simple Moving Average Indicator
type EMA struct {
	usage          indicatorUsageType
	LookbackPeriod int
}

// NewEMA returns a new Exponential Moving Average (EMA) configured with the
// specified lookbackPeriod
func NewEMA(lookbackPeriod int) (indicator *EMA, err error) {
	return &EMA{usage: SubChart, LookbackPeriod: lookbackPeriod}, nil
}

// Calculates the Exponential Moving Average (EMA) for the specified source.
// The source data must contain more items than the configured LookbackPeriod.
func (ema *EMA) Calculate(source []float64) (results []float64, err error) {

	// perform some sanity checks on the source data
	err = checkSourceDataIsNotEmpty(source)
	if err != nil {
		return nil, err
	}
	// and again on the source data with regards to the lookback period
	err = checkSourceValidForLookbackPeriod(source, ema.LookbackPeriod)
	if err != nil {
		return nil, err
	}

	// compute local variables for use in the loop
	periodTotal := float64(0.0)
	sourceLength := len(source)
	outputLength := sourceLength - ema.LookbackPeriod + 1

	// initialise the output data array
	results = make([]float64, outputLength)

	// iterate the source data

	y := (ema.LookbackPeriod * -1)
	for i := 0; i < sourceLength; i++ {
		y += 1

		if y > 0 {
			periodTotal -= source[i-ema.LookbackPeriod]
		}
		periodTotal += source[i]
		result := periodTotal / float64(ema.LookbackPeriod)
		if y >= 0 {
			results[y] = result
		}
	}

	return results, nil
}
