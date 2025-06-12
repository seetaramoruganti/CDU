package ws

import "time"

// AggregateUpdate contains one sample from each sensor
type AggregateUpdate struct {
	Time          time.Time `json:"time"`
	CPU_Temp      float64   `json:"cpuTemp"`
	CoolantTemp   float64   `json:"coolantTemp"`
	Vibration     float64   `json:"vibration"`
	LeakDetected  bool      `json:"leakDetected"`
	SmokeDetected bool      `json:"smokeDetected"`
	CoolantLevel  float64   `json:"coolantLevel"`
	FlameDetected bool      `json:"flameDetected"`
	MotorOn       bool      `json:"motorOn"`
	InFlowrate    float64   `json:"inFlowrate"`
	OutFlowrate   float64   `json:"outFlowrate"`
}
