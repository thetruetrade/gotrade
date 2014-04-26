// Exponential Moving Average (EMA)
package indicators

// A Simple Moving Average Indicator
type EMA struct {
	LookbackPeriod int
}

// NewEMA returns a new Exponential Moving Average (EMA) configured with the
// specified lookbackPeriod
func NewEMA(lookbackPeriod int) (indicator *EMA, err error) {
	return &EMA{LookbackPeriod: lookbackPeriod}, nil
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
	multiplier := float64(2.0 / float64(ema.LookbackPeriod+1.0))
	sourceLength := len(source)
	outputLength := sourceLength - ema.LookbackPeriod + 1

	// initialise the output data array
	results = make([]float64, outputLength)

	// initialise the previousEMA to a SMA for the same period
	previousEMA := initialiseEMAWithSMA(source, ema.LookbackPeriod)
	results[0] = previousEMA
	y := 1
	for i := ema.LookbackPeriod; i < sourceLength; i++ {
		results[y] = (source[i]-previousEMA)*multiplier + previousEMA
		previousEMA = results[y]
		y++
	}

	return results, nil
}

func initialiseEMAWithSMA(source []float64, lookbackPeriod int) float64 {
	periodTotal := float64(0.0)
	y := (lookbackPeriod * -1)
	i := 0
	for y < 0 {
		y += 1
		periodTotal += source[i]
		i += 1
	}
	return periodTotal / float64(lookbackPeriod)
}
