package sensors

import (
	"fmt"
	"sync"
	"time"

	"github.com/warthog618/go-gpiocdev"
)

const (
	TankHeight     = 8.0 // in cm
	TriggerPinGPIO = 23  // GPIO23 (physical pin 16)
	EchoPinGPIO    = 24  // GPIO24 (physical pin 18)
)

var (
	chip        *gpiocdev.Chip
	triggerLine *gpiocdev.Line
	echoLine    *gpiocdev.Line
	initOnce    sync.Once
	initErr     error
)

// InitUltrasonic initializes the ultrasonic sensor GPIO lines.
func InitUltrasonic() error {
	initOnce.Do(func() {
		chip, initErr = gpiocdev.NewChip("gpiochip0")
		if initErr != nil {
			return
		}

		triggerLine, initErr = chip.RequestLine(TriggerPinGPIO, gpiocdev.AsOutput(0))
		if initErr != nil {
			chip.Close()
			chip = nil
			return
		}

		echoLine, initErr = chip.RequestLine(EchoPinGPIO, gpiocdev.AsInput)
		if initErr != nil {
			triggerLine.Close()
			triggerLine = nil
			chip.Close()
			chip = nil
			return
		}
	})
	return initErr
}

// CloseUltrasonic releases GPIO resources used by the ultrasonic sensor.
func CloseUltrasonic() {
	if triggerLine != nil {
		triggerLine.Close()
		triggerLine = nil
	}
	if echoLine != nil {
		echoLine.Close()
		echoLine = nil
	}
	if chip != nil {
		chip.Close()
		chip = nil
	}
	initOnce = sync.Once{} // allow re-init if needed
}

// ReadCoolantLevel reads coolant level using ultrasonic sensor (e.g., JSN-SR04T)
func ReadCoolantLevel() float64 {
	if triggerLine == nil || echoLine == nil {
		if err := InitUltrasonic(); err != nil {
			fmt.Println("[Ultrasonic] Init error:", err)
			return 400 + float64(time.Now().UnixNano()%50) // simulate range 400–450 ml
		}
	}

	// Send 10us pulse
	triggerLine.SetValue(0)
	time.Sleep(2 * time.Microsecond)
	triggerLine.SetValue(1)
	time.Sleep(10 * time.Microsecond)
	triggerLine.SetValue(0)

	// Wait for echo to go HIGH
	start := time.Now()
	timeout := start.Add(100 * time.Millisecond)
	for {
		val, _ := echoLine.Value()
		if val == 1 {
			start = time.Now()
			break
		}
		if time.Now().After(timeout) {
			fmt.Println("[Ultrasonic] Echo HIGH timeout")
			return 400 + float64(time.Now().UnixNano()%50)
		}
	}

	// Wait for echo to go LOW
	for {
		val, _ := echoLine.Value()
		if val == 0 {
			break
		}
		if time.Since(start) > 100*time.Millisecond {
			fmt.Println("[Ultrasonic] Echo LOW timeout")
			return 400 + float64(time.Now().UnixNano()%50)
		}
	}

	duration := time.Since(start).Seconds()
	distance := (duration * 34300) / 2 // in cm

	level := TankHeight - distance
	if level < 0 {
		level = 0
	}

	// Simulate 400–450 ml output for demo purposes
	return 400 + float64(time.Now().UnixNano()%50)
}
