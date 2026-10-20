package repository

import (
	"context"
	"errors"
	"testing"

	"car-monitoring-system/internal/models"
	"car-monitoring-system/internal/tenant"
)

func TestTenantRepositoriesFailClosedWithoutContext(t *testing.T) {
	ctx := context.Background()
	cars := NewCarRepository(nil)
	telemetry := NewTelemetryRepository(nil)

	tests := []struct {
		name string
		err  error
	}{
		{"create car", cars.Create(ctx, models.Car{})},
		{"list cars", func() error { _, err := cars.GetAll(ctx); return err }()},
		{"get car", func() error { _, err := cars.GetByVin(ctx, "VIN"); return err }()},
		{"delete car", cars.Delete(ctx, "VIN")},
		{"update car", cars.UpdateSimulationStatus(ctx, "VIN", true)},
		{"get metric status", func() error { _, err := telemetry.GetCurrentStatus(ctx, "VIN", "tire_pressure"); return err }()},
		{"get full status", func() error { _, err := telemetry.GetFullStatus(ctx, "VIN"); return err }()},
		{"insert event", telemetry.InsertEvent(ctx, "VIN", "tire_pressure", models.WARNING)},
		{"close event", telemetry.CloseEvent(ctx, "VIN", "tire_pressure")},
		{"list events", func() error { _, err := telemetry.GetEvents(ctx, "VIN", ""); return err }()},
		{"upsert status", telemetry.UpsertAttributeStatus(ctx, "VIN", nil, nil)},
		{"simulation status", func() error { _, err := telemetry.IsSimulationEnabled(ctx, "VIN"); return err }()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !errors.Is(test.err, tenant.ErrNoTenant) {
				t.Fatalf("expected ErrNoTenant, got %v", test.err)
			}
		})
	}
}
