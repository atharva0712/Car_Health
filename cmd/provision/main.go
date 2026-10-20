package main

import (
	"car-monitoring-system/internal/db"
	"car-monitoring-system/internal/logger"
	"car-monitoring-system/internal/repository"
	"car-monitoring-system/internal/tenant"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	name := flag.String("name", "", "tenant name")
	flag.Parse()
	if strings.TrimSpace(*name) == "" {
		fmt.Fprintln(os.Stderr, "tenant name is required: go run ./cmd/provision --name \"Development\"")
		os.Exit(2)
	}

	log := logger.NewLogger()
	key, err := tenant.GenerateAPIKey()
	if err != nil {
		log.Errorf("Failed to generate API key: %v", err)
		os.Exit(1)
	}

	database := db.Connect()
	defer database.Close()

	_, err = repository.NewTenantRepository(database).Create(context.Background(), *name, tenant.HashAPIKey(key))
	if err != nil {
		log.Errorf("Failed to provision tenant: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Tenant provisioned. Save this API key now; it will not be shown again:\n%s\n", key)
}
