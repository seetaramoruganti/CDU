package sensors

import (
	"fmt"

	"github.com/warthog618/go-gpiocdev"
)

const VibrationGPIO = 27

// ReadVibrationFromSensor reads digital input from GPIO 21.
// Returns 1.0 for vibration detected, 0.0 otherwise.
func ReadVibration() float64 {
	chip, err := gpiocdev.NewChip("gpiochip0")
	if err != nil {
		fmt.Println("[Vibration] Chip error:", err)
		return 0.0
	}
	defer chip.Close()

	line, err := chip.RequestLine(VibrationGPIO, gpiocdev.AsInput)
	if err != nil {
		fmt.Println("[Vibration] Line request error:", err)
		return 0.0
	}
	defer line.Close()

	val, err := line.Value()
	if err != nil {
		fmt.Println("[Vibration] Read error:", err)
		return 0.0
	}

	// Optionally debounce or read multiple times here
	return float64(val)
}
