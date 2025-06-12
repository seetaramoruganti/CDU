package sensors

import (
	"fmt"

	"github.com/warthog618/go-gpiocdev"
)

const MotorGPIO = 17 // BCM 17 (physical pin 11)

var line *gpiocdev.Line



func InitMotor() error { 
	chip, err := gpiocdev.NewChip("gpiochip0")
	if err != nil {
		return fmt.Errorf("InitMotor: failed to open gpiochip: %w", err)
	}
	l, err := chip.RequestLine(MotorGPIO,
		gpiocdev.AsOutput(0),
		gpiocdev.WithConsumer("motor"))
	if err != nil {
		return fmt.Errorf("InitMotor: failed to request line: %w", err)
	}
	line = l
	return nil
}

func SetMotor(on bool) {
	if line == nil {
		fmt.Println("[Motor] not initialized")
		return
	}
	val := 0
	if on {
		val = 1
	}
	if err := line.SetValue(val); err != nil {
		fmt.Println("[Motor] SetValue error:", err)
	}
}

func CloseMotor() {
	if line != nil {
		_ = line.SetValue(0)
		line.Close()
		line = nil
	}
}

func EvaluateMotorState(tempC, levelCm float64) bool {
	return tempC > 27.5 && levelCm > 400
}
