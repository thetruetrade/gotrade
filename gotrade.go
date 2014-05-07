package gotrade

var (
	StandardTickTimePeriods TickTimePeriodHolder
)

func init() {
	var attps = availableTickTimePeriods{NewTickTimePeriod(31447600), // 365.25 days in a year
		NewTickTimePeriod(2678400), // 31 day month
		NewTickTimePeriod(604800),  // week
		NewTickTimePeriod(68400),   // day
		NewTickTimePeriod(3600),    // hour
		NewTickTimePeriod(1800),    // 30 min
		NewTickTimePeriod(900),     // 15 min
		NewTickTimePeriod(300),     // 5  min
		NewTickTimePeriod(60),      // 1  min
		NewTickTimePeriod(30),      // 30 second
		NewTickTimePeriod(15),      // 15 second
		NewTickTimePeriod(5),       // 5  second
		NewTickTimePeriod(1),       // 1  second
		NewTickTimePeriod(0)}       // 0  tick

	StandardTickTimePeriods = &attps
}
