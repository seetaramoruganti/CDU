package sensors

import (
	"fmt"
	"time"

	"github.com/warthog618/go-gpiocdev"
)

// FlowMeter encapsulates a single FL-308D sensor on a BCM pin.
type FlowMeter struct {
	chip        *gpiocdev.Chip
	line        *gpiocdev.Line
	calibration float64 // pulses per liter
}

// NewFlowMeter initializes the GPIO line (with pull-down) for the given BCM pin.
//
//	pin:         BCM pin number (after you’ve level-shifted the sensor output to 3.3 V)
//	calibration: pulses per liter (e.g. 450.0)
func NewFlowMeter(pin int, calibration float64) (*FlowMeter, error) {
	chip, err := gpiocdev.NewChip("gpiochip0")
	if err != nil {
		return nil, fmt.Errorf("open gpiochip0: %w", err)
	}

	// Note: AsInput and WithPullDown are constants, not functions.
	line, err := chip.RequestLine(
		pin,
		gpiocdev.AsInput,
		gpiocdev.WithPullDown,
	)
	if err != nil {
		chip.Close()
		return nil, fmt.Errorf("request GPIO %d: %w", pin, err)
	}

	return &FlowMeter{
		chip:        chip,
		line:        line,
		calibration: calibration,
	}, nil
}

// Measure blocks for exactly "duration", polls the GPIO every millisecond,
// counts LOW→HIGH edges, then returns flow in liters per minute.
func (fm *FlowMeter) Measure(duration time.Duration) (float64, error) {
	var prevVal int
	var pulseCount uint

	start := time.Now()
	for time.Since(start) < duration {
		// The current API uses Value() instead of GetValue()
		val, err := fm.line.Value()
		if err != nil {
			return 0, fmt.Errorf("read GPIO: %w", err)
		}

		// rising edge detection
		if prevVal == 0 && val == 1 {
			pulseCount++
		}
		prevVal = val

		time.Sleep(1 * time.Millisecond)
	}

	// Compute flow (L/min):
	//   liters during “duration” = pulseCount / calibration
	//   L/min = liters * (60 / durationSeconds)
	liters := float64(pulseCount) / fm.calibration
	flowLpm := liters * (60.0 / duration.Seconds())
	return flowLpm, nil
}

// Close releases the GPIO line and chip.
func (fm *FlowMeter) Close() {
	fm.line.Close()
	fm.chip.Close()
}
