package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var _ Collector = new(PrometheusCollector)

type PrometheusCollector struct {
	appName  string
	counters map[string]any
}

func (p *PrometheusCollector) NewCounter(name string) Counter {
	if _, exists := p.counters[name]; !exists {
		p.counters[name] = Counter{
			promauto.NewCounter(prometheus.CounterOpts{
				Namespace: p.appName,
				Name:      name,
			}),
		}
	}

	return p.counters[name].(Counter)
}

func (p *PrometheusCollector) NewDimensionedCounter(name string, labels ...string) DimensionedCounter {
	if _, exists := p.counters[name]; !exists {
		p.counters[name] = DimensionedCounter{
			promauto.NewCounterVec(prometheus.CounterOpts{
				Namespace: p.appName,
				Name:      name,
			}, labels),
		}
	}

	return p.counters[name].(DimensionedCounter)

}

func (p *PrometheusCollector) NewGauge(name string) Gauge {
	if _, exists := p.counters[name]; !exists {
		p.counters[name] = Gauge{
			promauto.NewGauge(prometheus.GaugeOpts{
				Namespace: p.appName,
				Name:      name,
			}),
		}
	}

	return p.counters[name].(Gauge)
}

func (p *PrometheusCollector) NewDimensionedGauge(name string, labels ...string) DimensionedGauge {
	if _, exists := p.counters[name]; !exists {
		p.counters[name] = DimensionedGauge{
			promauto.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: p.appName,
				Name:      name,
			}, labels),
		}
	}

	return p.counters[name].(DimensionedGauge)

}

type Counter struct {
	prometheus.Counter
}

type DimensionedCounter struct {
	*prometheus.CounterVec
}

type Gauge struct {
	prometheus.Gauge
}

type DimensionedGauge struct {
	*prometheus.GaugeVec
}
