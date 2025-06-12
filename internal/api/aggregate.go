package api

import (
	"encoding/json"
	"net/http"
	"time"

	"sensor_dashboard_influx/internal/ws"
	"sensor_dashboard_influx/sensors"
)

// AggregateHandler returns a one‐off snapshot of all sensor readings.
// GET /api/aggregate
func AggregateHandler(w http.ResponseWriter, r *http.Request) {
	// 1) Read each sensor
	tempDigital := sensors.ReadTempCelcius()
	tempAnalog, _ := sensors.ReadTemperatureFromADC(nil) // or pass your ADC instance if available
	level := sensors.ReadCoolantLevel()
	vibration := sensors.ReadVibration()
	leak := sensors.ReadCoolantLeak()
	smoke := sensors.ReadGasAlert()
	flame := sensors.ReadFlameDetected()
	motorOn := sensors.EvaluateMotorState(tempDigital, level)

	inFlowMeter, err := sensors.NewFlowMeter(22, 450.0)
	var inFlow float64
	if err == nil {
		inFlow, _ = inFlowMeter.Measure(1 * time.Second)
	} else {
		inFlow = 0 // or handle error as needed
	}

	outFlowMeter, err := sensors.NewFlowMeter(16, 450.0)
	var outFlow float64
	if err == nil {
		outFlow, _ = outFlowMeter.Measure(1 * time.Second)
	} else {
		outFlow = 0 // or handle error as needed
	}

	// 2) Build the same AggregateUpdate struct
	agg := ws.AggregateUpdate{
		Time:          time.Now(),
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

	// 3) JSON‐encode and return
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(agg); err != nil {
		http.Error(w, "encode error: "+err.Error(), http.StatusInternalServerError)
	}
}
