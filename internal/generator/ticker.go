// package generator

// import (
// 	"bytes"
// 	"car-monitoring-system/internal/config"
// 	"car-monitoring-system/internal/logger"
// 	"car-monitoring-system/internal/models"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"sync"
// 	"time"
// )

// type Ticker struct {
// 	generator  *TelemetryGenerator
// 	backendURL string
// 	interval   time.Duration
// 	Log        *logger.Logger
// 	workers    int
// }

// // NewTicker creates a new Ticker instance with the provided configuration.
// func NewTicker(cfg *config.Config, log *logger.Logger) *Ticker {
// 	return &Ticker{
// 		generator:  NewTelemetryGenerator(cfg),
// 		backendURL: cfg.Backend.URL,
// 		interval:   time.Duration(cfg.Generation.IntervalSeconds) * time.Second,
// 		Log:        log,
// 		workers:    5,
// 	}
// }

// // Start begins the ticker that generates and sends telemetry data at regular intervals.

// func (t *Ticker) Start() {
// 	ticker := time.NewTicker(t.interval)
// 	defer ticker.Stop()

// 	t.Log.Infof("Starting telemetry ticker with interval: %v", t.interval)

// 	t.tick() // Generate immediately on start

// 	// Then run on schedule
// 	for range ticker.C {
// 		t.tick()
// 	}
// }

// func (t *Ticker) tick() {
// 	//Fetch all cars from backend

// 	vins, err := t.fetchVINs()
// 	if err != nil {
// 		t.Log.Errorf("Error fetching VINs: %v", err)
// 		return
// 	}
// 	if len(vins) == 0 {
// 		t.Log.Info("No simulation-enabled cars found")
// 		return
// 	}

// 	t.Log.Infof("Generating telemetry for %d cars", len(vins))

// 	// Worker Pool
// 	jobs := make(chan string)
// 	var wg sync.WaitGroup

// 	// Start workers
// 	for i := 0; i < t.workers; i++ {
// 		wg.Add(1)
// 		go t.worker(i, jobs, &wg)
// 	}

// 	// send jobs
// 	for _, vin := range vins {
// 		jobs <- vin
// 	}

// 	close(jobs)

// 	wg.Wait()

// 	t.Log.Info("Telemetry cycle completed")

// }
// func (t *Ticker) worker(id int, jobs <-chan string, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	for vin := range jobs {
// 		payload := t.generator.GenerateTelemetry(vin)

// 		err := t.sendTelemetry(payload)
// 		if err != nil {
// 			t.Log.Errorf("[Worker %d] Error sending telemetry for VIN %s: %v", id, vin, err)
// 		} else {
// 			t.Log.Debugf("[Worker %d] Telemetry send for VIN %s", id, vin)
// 		}
// 	}
// }

// func (t *Ticker) fetchVINs() ([]string, error) {
// 	resp, err := http.Get(t.backendURL + "/cars")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var cars []models.Car
// 	err = json.NewDecoder(resp.Body).Decode(&cars)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Filter only simulation_enabled = true

// 	var vins []string
// 	for _, car := range cars {
// 		if car.SimulationEnabled {
// 			vins = append(vins, car.VIN)
// 		}
// 	}
// 	return vins, nil
// }

// func (t *Ticker) sendTelemetry(payload models.TelemetryPayload) error {
// 	data, err := json.Marshal(payload)
// 	if err != nil {
// 		return err
// 	}

// 	resp, err := http.Post(
// 		t.backendURL+"/telemetry",
// 		"application/json",
// 		bytes.NewBuffer(data),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

//		if resp.StatusCode != http.StatusOK {
//			return fmt.Errorf("backend returned status: %d", resp.StatusCode)
//		}
//		return nil
//	}
package generator

import (
	"bytes"
	"car-monitoring-system/internal/config"
	"car-monitoring-system/internal/logger"
	"car-monitoring-system/internal/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Ticker struct {
	config     *config.Config
	backendURL string
	interval   time.Duration
	workers    int
	Log        *logger.Logger
	httpClient *http.Client
}

func NewTicker(cfg *config.Config, log *logger.Logger) *Ticker {
	return &Ticker{
		config:     cfg,
		backendURL: cfg.Backend.URL,
		interval:   time.Duration(cfg.Generation.IntervalSeconds) * time.Second,
		workers:    5,
		Log:        log,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (t *Ticker) Start() {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	t.Log.Infof("Starting telemetry ticker with intervals: %v", t.interval)

	t.tick()

	for range ticker.C {
		t.tick()
	}
}

func (t *Ticker) tick() {
	vins, err := t.fetchVINs()
	if err != nil {
		t.Log.Errorf("Error fetching VINs: %v", err)
		return
	}
	if len(vins) == 0 {
		t.Log.Info("No simulation-enabled cars found")
		return
	}
	t.Log.Infof("Generating telemetry for %d cars using %d workers", len(vins), t.workers)

	jobs := make(chan string)
	var wg sync.WaitGroup

	// Starting workers
	for i := 0; i < t.workers; i++ {
		wg.Add(1)
		go t.worker(i, jobs, &wg)
	}

	// Sending VINs to workers
	for _, vin := range vins {
		jobs <- vin
	}

	close(jobs)

	wg.Wait()

	t.Log.Info("Telemetry cycle completed")
}

func (t *Ticker) worker(id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	seed := time.Now().UnixNano() + int64(id)*1000
	localGenerator := &TelemetryGenerator{
		Config: t.config,
		rng:    rand.New(rand.NewSource(seed)),
	}

	for vin := range jobs {
		payload := localGenerator.GenerateTelemetry(vin)

		err := t.sendTelemetry(payload)
		if err != nil {
			t.Log.Errorf("[Worker %d] Error sending telemetry for VIN %s: %v", id, vin, err)
		} else {
			t.Log.Debugf("[Worker %d] Telemetry send for VIN %s", id, vin)
		}
	}
}

func (t *Ticker) fetchVINs() ([]string, error) {
	resp, err := t.httpClient.Get(t.backendURL + "/cars")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cars []models.Car
	err = json.NewDecoder(resp.Body).Decode(&cars)
	if err != nil {
		return nil, err
	}

	//Filter only simulation_enabled = true

	var vins []string
	for _, car := range cars {
		if car.SimulationEnabled {
			vins = append(vins, car.VIN)
		}
	}
	return vins, nil
}

func (t *Ticker) sendTelemetry(payload models.TelemetryPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := t.httpClient.Post(
		t.backendURL+"/telemetry",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("backend returned status: %d", resp.StatusCode)
	}
	return nil
}
