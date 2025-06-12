
package sensors

import (
	"fmt"
	"github.com/warthog618/go-gpiocdev"
)

const FlameSensorGPIO = 19 // GPIO19 (Physical pin 35)

func ReadFlameDetected() bool {
	chip, err := gpiocdev.NewChip("gpiochip0")
	if err != nil {
		fmt.Println("[FlameSensor] Chip error:", err)
		return false
	}
	defer chip.Close()

	line, err := chip.RequestLine(FlameSensorGPIO, gpiocdev.AsInput)
	if err != nil {
		fmt.Println("[FlameSensor] Line request error:", err)
		return false
	}
	defer line.Close()

	val, err := line.Value()
	if err != nil {
		fmt.Println("[FlameSensor] Read error:", err)
		return false
	}

	return val == 0 // Flame detected when GPIO LOW
}
