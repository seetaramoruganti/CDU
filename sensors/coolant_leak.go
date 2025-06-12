package sensors

import (
	"fmt"
	"github.com/warthog618/go-gpiocdev"
)



const (
	CoolantLeakGPIO  = 13    // GPIO13 (physical pin 33)
)


// ReadCoolantLeak reads digital signal from FC-28 sensor
func ReadCoolantLeak() bool {
	chip, err := gpiocdev.NewChip("gpiochip0")
	if err != nil {
		fmt.Println("[CoolantLeak] Chip error:", err)
		return false
	}
	defer chip.Close()

	line, err := chip.RequestLine(CoolantLeakGPIO, gpiocdev.AsInput)
	if err != nil {
		fmt.Println("[CoolantLeak] Line request error:", err)
		return false
	}
	defer line.Close()

	val, err := line.Value()
	if err != nil {
		fmt.Println("[CoolantLeak] Read error:", err)
		return false
	}

	return val == 0 // Leak detected when GPIO LOW
}
