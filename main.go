package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"sensor_dashboard_influx/db"
	"sensor_dashboard_influx/internal/api"
	"sensor_dashboard_influx/internal/ws"
	"sensor_dashboard_influx/sensors"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	// 1. Create and run the hub
	hub := ws.NewHub()
	go hub.Run()

	// 2. Setup HTTP router
	r := mux.NewRouter()
	r.HandleFunc("/api/history", api.HistoryHandler).Methods("GET")
	r.HandleFunc("/api/aggregate", api.AggregateHandler).Methods("GET")
	r.HandleFunc("/stream", ws.ServeWS(hub))

	// 3. Wrap with CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	}).Handler(r)

	// 4. Start HTTP server in background
	go func() {
		log.Printf("listening on :%s", port)
		if err := http.ListenAndServe(
			":"+port,
			handler,
		); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	// 5. Initialize sensors and motor
	if err := sensors.InitMotor(); err != nil {
		log.Fatalf("Motor initialization failed: %v", err)
	}
	defer sensors.CloseMotor()

	// Setup shutdown signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	defer close(sigs)

	// Initialize ADC
	adc, err := sensors.InitPCF8591()
	if err != nil {
		log.Println("Warning: ADC init failed, continuing without analog temp sensor.")
		adc = nil
	}
	defer func() {
		if adc != nil {
			adc.Close()
		}
	}()

	// Initialize flow meters
	fm, err := sensors.NewFlowMeter(22, 450.0)
	if err != nil {
		log.Fatalf("Failed to init inflow meter: %v", err)
	}
	defer fm.Close()

	fom, err := sensors.NewFlowMeter(16, 450.0)
	if err != nil {
		log.Fatalf("Failed to init outflow meter: %v", err)
	}
	defer fom.Close()

	// 6. Main sensor-poll loop
	for {
		select {
		case <-sigs:
			sensors.SetMotor(false)
			defer sensors.CloseMotor()
			log.Println("Shutdown signal received, exiting.")
			return
		default:
		}

		// Read sensors
		tempDigital := sensors.ReadTempCelcius()
		level := sensors.ReadCoolantLevel()
		vibration := sensors.ReadVibration()
		leak := sensors.ReadCoolantLeak()
		flame := sensors.ReadFlameDetected()
		smoke := sensors.ReadGasAlert()

		// Analog temp with retry
		tempAnalog, err := sensors.ReadTemperatureFromADC(adc)
		if adc != nil && err != nil {
			log.Println("Analog Temp read error, reinitializing ADC:", err)
			adc.Close()
			time.Sleep(time.Second)
			adc, _ = sensors.InitPCF8591()
		}

		// Motor control
		motorOn := sensors.EvaluateMotorState(tempDigital, level)
		sensors.SetMotor(motorOn)

		// Flow rates
		inFlow, _ := fm.Measure(time.Second)
		outFlow, _ := fom.Measure(time.Second)

		// Log readings
		fmt.Printf("[TEMP] %.2f°C \n, [COOLANT] %.2fml \n, [VIBRATION] %.2f \n, [LEAK] %v \n, [FLAME] %v \n, [SMOKE] %v \n, [Motor Status] %v \n, [IN FLOW RATE] %.2f \n, [OUT FLOW RATE] %.2f \n",
			tempDigital, level, vibration, leak, flame, smoke, motorOn ,inFlow, outFlow)



		// Build and broadcast aggregate
		now := time.Now()
		agg := ws.AggregateUpdate{
			Time:          now,
			CPU_Temp:      tempDigital,
			CoolantTemp:   tempAnalog,
			Vibration:     vibration,
			LeakDetected:  leak,
			SmokeDetected: smoke,
			FlameDetected: flame,
			MotorOn:       motorOn,
			CoolantLevel:  level,
			InFlowrate:    inFlow,
			OutFlowrate:   outFlow,
		}
		if data, err := json.Marshal(agg); err == nil {
			hub.Broadcast <- data
		}

		// Push to InfluxDB
		timestamp := now.Truncate(time.Second)
		db.PushTemperature(tempDigital, timestamp)
		db.PushCoolant(level, timestamp)
		db.PushVibration(vibration, timestamp)
		db.PushCoolantLeak(leak, timestamp)
		db.PushFlameStatus(flame, timestamp)
		db.PushSmokeStatus(smoke, timestamp)
		db.PushMotorState(motorOn, timestamp)
		db.PushInFlowrate(inFlow, timestamp)
		db.PushOutFlowrate(outFlow, timestamp)
		db.PushCoolantTemp(tempAnalog, timestamp)

		time.Sleep(4 * time.Second)
	}
}
