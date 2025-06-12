package sensors

import (
	"fmt"
	"os/exec"
    "strconv"
    "strings"
	


)
 
func ReadHumidity() float64 {
    out, err := exec.Command("python3", "read_dht.py").Output()
    if err != nil {
        fmt.Println("[ERROR] Python script failed:", err)
        return -1.0
    }

    result := strings.TrimSpace(string(out))
    humidity, err := strconv.ParseFloat(result, 64)
    if err != nil {
        fmt.Println("[ERROR] Failed to parse humidity:", result)
        return -1.0
    }

    return humidity
}


