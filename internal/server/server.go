package server

import (
	"car-monitoring-system/internal/config"
	"car-monitoring-system/internal/handler"
	"car-monitoring-system/internal/logger"
	"car-monitoring-system/internal/repository"
	"car-monitoring-system/internal/service"
	"car-monitoring-system/internal/tenant"
	"database/sql"
	"net/http"
)

func NewHandler(database *sql.DB, cfg *config.Config, log *logger.Logger) http.Handler {
	carRepo := repository.NewCarRepository(database)
	carService := service.NewCarService(carRepo, log)
	carHandler := handler.NewCarHandler(carService, log)

	telemetryRepo := repository.NewTelemetryRepository(database)
	telemetryService := service.NewTelemetryService(telemetryRepo, cfg, log)
	telemetryHandler := handler.NewTelemetryHandler(telemetryService, log)
	statusHandler := handler.NewStatusHandler(telemetryService, log)

	tenantRepo := repository.NewTenantRepository(database)
	auth := tenant.NewAuthMiddleware(tenantRepo, log)

	mux := http.NewServeMux()
	mux.Handle("POST /cars", auth.Wrap(http.HandlerFunc(carHandler.RegisterCar)))
	mux.Handle("GET /cars", auth.Wrap(http.HandlerFunc(carHandler.GetCars)))
	mux.Handle("GET /cars/{vin}", auth.Wrap(http.HandlerFunc(carHandler.GetCarByVIN)))
	mux.Handle("DELETE /cars/{vin}", auth.Wrap(http.HandlerFunc(carHandler.DeleteCar)))
	mux.Handle("PATCH /cars/{vin}/simulation", auth.Wrap(http.HandlerFunc(carHandler.UpdateSimulationStatus)))
	mux.Handle("POST /telemetry", auth.Wrap(http.HandlerFunc(telemetryHandler.ReceiveTelemetry)))
	mux.Handle("GET /cars/{vin}/status", auth.Wrap(http.HandlerFunc(statusHandler.GetCarStatus)))
	mux.Handle("GET /cars/{vin}/events", auth.Wrap(http.HandlerFunc(statusHandler.GetEvents)))

	return mux
}
