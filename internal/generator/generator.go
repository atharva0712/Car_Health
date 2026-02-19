package generator

import (
	"car-monitoring-system/internal/config"
	"car-monitoring-system/internal/models"
	"math/rand"
	"time"
)

type TelemetryGenerator struct { //Responsible for generating telemetry data based on the provided configuration. Uses a random number generator to create realistic telemetry values for different metrics defined in the configuration.
	Config *config.Config
	rng    *rand.Rand
}

func NewTelemetryGenerator(cfg *config.Config) *TelemetryGenerator {
	return &TelemetryGenerator{
		Config: cfg,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateTelemetry generates telemetry data for a given VIN. It returns a TelemetryPayload struct containing the generated telemetry data.
func (g *TelemetryGenerator) GenerateTelemetry(vin string) models.TelemetryPayload {
	return models.TelemetryPayload{
		VIN:            vin,
		TirePressure:   g.generateMetricValue("tire_pressure"),
		EngineTemp:     g.generateMetricValue("engine_temp"),
		BatteryVoltage: g.generateMetricValue("battery_voltage"),
		OilLevel:       g.generateMetricValue("oil_level"),
		BrakeCondition: g.generateMetricValue("brake_condition"),
	}
}

func (g *TelemetryGenerator) generateMetricValue(metricName string) float64 {
	metric := g.Config.Metrics[metricName]
	dist := g.Config.Generation.Distribution

	roll := g.rng.Intn(100)

	if roll < dist.Normal {
		//Generate Normal value
		return g.randomInRange(metric.Normal)
	} else if roll < dist.Normal + dist.Warning {
		//Generate Warning value
		return g.randomInRange(metric.Warning)
	} else {
		//Generate Critical value
		return g.randomInRange(metric.Critical)
	}
}

// randomInRange generates a random float64 value within the specified range.
func (g *TelemetryGenerator) randomInRange(r config.Range) float64 {
	// Handle tire pressure with two warning ranges
	if r.Max2 > 0 {
		// Randomly pick first or second range
		if g.rng.Intn(2) == 0 {
			return r.Min + g.rng.Float64()*(r.Max-r.Min)
		} else {
			return r.Min2 + g.rng.Float64()*(r.Max2-r.Min2)
		}
	}
	// Single range
	return r.Min +g.rng.Float64()*(r.Max-r.Min)
}


