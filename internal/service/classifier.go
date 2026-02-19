package service

import "car-monitoring-system/internal/models"

func ClassifyTirePressure(value float64) models.Status {
	if value >= 30 && value <= 35 {
		return models.NORMAL
	}
	if (value >= 28 && value <= 29) || (value >= 36 && value <= 37) {
		return models.WARNING
	}

	return models.CRITICAL
}

func ClassifyEngineTemp(value float64) models.Status {

	if value <= 95 {
		return models.NORMAL
	}

	if value <= 105 {
		return models.WARNING
	}

	return models.CRITICAL
}

func ClassifyBattery(value float64) models.Status {

	if value >= 12 {
		return models.NORMAL
	}

	if value >= 11 {
		return models.WARNING
	}

	return models.CRITICAL
}

func ClassifyOil(value float64) models.Status {
	if value >= 50 {
		return models.NORMAL
	}

	if value >= 40 {
		return models.WARNING
	}

	return models.CRITICAL
}

func ClassifyBrake(value float64) models.Status {

	if value >= 60 {
		return models.NORMAL
	}

	if value >= 50 {
		return models.WARNING
	}

	return models.CRITICAL
}