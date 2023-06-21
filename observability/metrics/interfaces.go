package metrics

type Collector interface {
	NewCounter(string) Counter
	NewDimensionedCounter(string, ...string) DimensionedCounter
	NewGauge(string) Gauge
	NewDimensionedGauge(string, ...string) DimensionedGauge
}
