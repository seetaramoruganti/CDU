package sensors

import (
	"fmt"
	"sync"

	"github.com/warthog618/go-gpiocdev"
)

const MotorGPIO = 17 // BCM 17 (physical pin 11)

var line *gpiocdev.Line

var (
	mu          sync.Mutex
	manual      bool
	manualState bool
)

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

func setMotorLocked(on bool) {
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

// SetMotor sets the motor state without enabling manual override.
func SetMotor(on bool) {
	mu.Lock()
	defer mu.Unlock()
	setMotorLocked(on)
}

func CloseMotor() {
	if line != nil {
		_ = line.SetValue(0)
		line.Close()
		line = nil
		mu.Lock()
		manual = false
		manualState = false
		mu.Unlock()
	}
}

func EvaluateMotorState(tempC, levelCm float64) bool {
	return tempC > 27.5 && levelCm > 400
}

// SetManualMotor enables manual override and sets the motor accordingly.
func SetManualMotor(on bool) {
	mu.Lock()
	defer mu.Unlock()
	manual = true
	manualState = on
	setMotorLocked(on)
}

// ClearManualMotor disables the manual override, allowing automatic control.
func ClearManualMotor() {
	mu.Lock()
	manual = false
	mu.Unlock()
}

// ManualOverride returns whether manual control is active.
func ManualOverride() bool {
	mu.Lock()
	defer mu.Unlock()
	return manual
}

// ManualState returns the desired manual motor state.
func ManualState() bool {
	mu.Lock()
	defer mu.Unlock()
	return manualState
}
