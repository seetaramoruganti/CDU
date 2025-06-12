package sensors

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ReadTempCelcius() float64 {
	basePath := "/sys/bus/w1/devices/"
	entries, err := os.ReadDir(basePath)
	if err != nil {
		fmt.Println("[TempSensor] Directory read error:", err)
		return 0.0
	}

	var sensorFile string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "28-") {
			sensorFile = filepath.Join(basePath, entry.Name(), "w1_slave")
			break
		}
	}

	if sensorFile == "" {
		fmt.Println("[TempSensor] No DS18B20 sensor found")
		return 0.0
	}

	data, err := os.ReadFile(sensorFile)
	if err != nil {
		fmt.Println("[TempSensor] Read error:", err)
		return 0.0
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 || !strings.HasSuffix(lines[0], "YES") {
		fmt.Println("[TempSensor] CRC check failed or bad data")
		return 0.0
	}

	parts := strings.Split(lines[1], "t=")
	if len(parts) < 2 {
		fmt.Println("[TempSensor] Temperature not found")
		return 0.0
	}

	raw, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		fmt.Println("[TempSensor] Parse error:", err)
		return 0.0
	}

	return float64(raw) / 1000.0
}
