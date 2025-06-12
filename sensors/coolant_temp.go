// sensors/coolant_temp.go

package sensors

import (
	"fmt"
	"time"

	i2c "github.com/d2r2/go-i2c"
)

// InitPCF8591 opens the I2C bus and returns a handle to the PCF8591 at address 0x48.
// The returned *i2c.I2C must be closed by the caller when no longer needed.
func InitPCF8591() (*i2c.I2C, error) {
	// bus 1, address 0x48
	return i2c.NewI2C(0x48, 1)
}

// ReadTemperatureFromADC performs one 8‐bit read from AIN0 of the PCF8591,
// converts that raw byte (0–255) into a voltage (0–3.3V), then converts into °C
// assuming an LM35 (10 mV/°C).
//
// The caller must have previously called InitPCF8591() and pass that *i2c.I2C here.
// If any I²C transaction times out or fails, it returns an error.
func ReadTemperatureFromADC(dev *i2c.I2C) (float64, error) {
	// 1) Write the control byte 0x40 to select AIN0 (single‐ended, no auto‐increment)
	if _, err := dev.WriteBytes([]byte{0x40}); err != nil {
		return 0, fmt.Errorf("I2C write failed: %w", err)
	}
	// small delay to let PCF8591 settle
	time.Sleep(10 * time.Millisecond)

	// 2) Dummy read: the first read after writing the control byte returns stale data
	dummy := make([]byte, 1)
	if _, err := dev.ReadBytes(dummy); err != nil {
		return 0, fmt.Errorf("dummy read failed: %w", err)
	}

	// 3) Actual read of the 8‐bit ADC value
	data := make([]byte, 1)
	if _, err := dev.ReadBytes(data); err != nil {
		return 0, fmt.Errorf("real read failed: %w", err)
	}

	raw := data[0]                             // 0–255
	voltage := float64(raw) * 3.3 / 255.0      // map to 0–3.300 V
	tempC := voltage * 10.0 - 2                  // LM35: 0.10 V per °C → multiply by 10

	// small delay to avoid hammering the bus if called rapidly
	time.Sleep(4 * time.Millisecond)
	return tempC, nil
}
