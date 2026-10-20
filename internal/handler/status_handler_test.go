package handler

import (
	"car-monitoring-system/internal/models"
	"testing"
)

func TestIsEmptyFullStatus(t *testing.T) {
	var status *models.AttributeStatus
	if !isEmptyFullStatus(status) {
		t.Fatal("typed nil full status should be treated as empty")
	}
	if isEmptyFullStatus(&models.AttributeStatus{}) {
		t.Fatal("initialized full status should not be treated as empty")
	}
	if isEmptyFullStatus(map[string]any{}) {
		t.Fatal("non-status response should not be treated as empty full status")
	}
}
