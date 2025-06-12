
package sensors

import (
	"fmt"
	"github.com/warthog618/go-gpiocdev"
)



const GasSensorGPIO = 26 // GPIO26 (Physical pin 37)

func ReadGasAlert() bool {
	chip, err := gpiocdev.NewChip("gpiochip0")
	if err != nil {
		fmt.Println("[GasSensor] Chip error:", err)
		return false
	}
	defer chip.Close()

	line, err := chip.RequestLine(GasSensorGPIO, gpiocdev.AsInput)
	if err != nil {
		fmt.Println("[GasSensor] Line request error:", err)
		return false
	}
	defer line.Close()

	val, err := line.Value()
	if err != nil {
		fmt.Println("[GasSensor] Read error:", err)
		return false
	}

	return val == 0 // Alert when gas is detected
}
