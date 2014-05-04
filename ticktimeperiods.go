package gotrade

type TickTimePeriod struct {
	secondsInPeriod int
}

func (ttp *TickTimePeriod) SecondsInPeriod() int {
	return ttp.secondsInPeriod
}

func NewTickTimePeriod(secondsInPeriod int) TickTimePeriod {
	return TickTimePeriod{secondsInPeriod}
}

type TickTimePeriodHolder interface {
	Yearly() TickTimePeriod
	Monthly() TickTimePeriod
	Weekly() TickTimePeriod
	Daily() TickTimePeriod
	Hourly() TickTimePeriod
	ThirtyMinute() TickTimePeriod
	FifteenMinute() TickTimePeriod
	FiveMinute() TickTimePeriod
	OneMinute() TickTimePeriod
	ThirtySecond() TickTimePeriod
	FifteenSecond() TickTimePeriod
	FiveSecond() TickTimePeriod
	OneSecond() TickTimePeriod
	Tick() TickTimePeriod
}

type availableTickTimePeriods struct {
	yearly        TickTimePeriod
	monthly       TickTimePeriod
	weekly        TickTimePeriod
	daily         TickTimePeriod
	hourly        TickTimePeriod
	thirtyMinute  TickTimePeriod
	fifteenMinute TickTimePeriod
	fiveMinute    TickTimePeriod
	oneMinute     TickTimePeriod
	thirtySecond  TickTimePeriod
	fifteenSecond TickTimePeriod
	fiveSecond    TickTimePeriod
	oneSecond     TickTimePeriod
	tick          TickTimePeriod
}

func (attp *availableTickTimePeriods) Yearly() TickTimePeriod {
	return attp.yearly
}

func (attp *availableTickTimePeriods) Monthly() TickTimePeriod {
	return attp.monthly
}

func (attp *availableTickTimePeriods) Weekly() TickTimePeriod {
	return attp.weekly
}

func (attp *availableTickTimePeriods) Daily() TickTimePeriod {
	return attp.daily
}

func (attp *availableTickTimePeriods) Hourly() TickTimePeriod {
	return attp.hourly
}

func (attp *availableTickTimePeriods) ThirtyMinute() TickTimePeriod {
	return attp.thirtyMinute
}

func (attp *availableTickTimePeriods) FifteenMinute() TickTimePeriod {
	return attp.fifteenMinute
}

func (attp *availableTickTimePeriods) FiveMinute() TickTimePeriod {
	return attp.fiveMinute
}

func (attp *availableTickTimePeriods) OneMinute() TickTimePeriod {
	return attp.oneMinute
}

func (attp *availableTickTimePeriods) ThirtySecond() TickTimePeriod {
	return attp.thirtySecond
}

func (attp *availableTickTimePeriods) FifteenSecond() TickTimePeriod {
	return attp.fifteenSecond
}

func (attp *availableTickTimePeriods) FiveSecond() TickTimePeriod {
	return attp.fiveSecond
}

func (attp *availableTickTimePeriods) OneSecond() TickTimePeriod {
	return attp.oneSecond
}

func (attp *availableTickTimePeriods) Tick() TickTimePeriod {
	return attp.tick
}
