package db

import (
	"fmt"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
)

const influxURL = "http://192.168.116.176:8086"
const dbName = "sensor_data"


func createClient() influx.Client {
	client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr: influxURL,
	})
	if err != nil {
		fmt.Println("[InfluxDB] Connection error:", err)
	}
	return client
}

func writePoint(measurement string, tags map[string]string, fields map[string]interface{}, timestamp time.Time) {
	client := createClient()
	if client == nil {
		return
	}
	defer client.Close()

	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  dbName,
		Precision: "s",
	})
	if err != nil {
		fmt.Println("[InfluxDB] BatchPoints error:", err)
		return
	}

	pt, err := influx.NewPoint(measurement, tags, fields, timestamp)
	if err != nil {
		fmt.Println("[InfluxDB] Point error:", err)
		return
	}
	bp.AddPoint(pt)

	err = client.Write(bp)
	if err != nil {
		fmt.Println("[InfluxDB] Write error:", err)
	}
}

func PushTemperature(temp float64, timestamp time.Time) {
	tags := map[string]string{
		"rack_id":   "rack01",
		"node_id":   "node02",
		"sensor_id": "cpu_temp",
	}
	fields := map[string]interface{}{
		"value": temp,
	}
	writePoint("temperature", tags, fields, timestamp)
}

func PushCoolant(level float64, timestamp time.Time) {
	tags := map[string]string{
		"rack_id":   "rack01",
		"node_id":   "node02",
		"sensor_id": "coolant_lvl",
	}
	fields := map[string]interface{}{
		"value": level,
	}
	writePoint("coolant", tags, fields, timestamp)
}

func PushVibration(vibration float64, timestamp time.Time) {
	tags := map[string]string{
		"rack_id":   "rack01",
		"node_id":   "node02",
		"sensor_id": "vibration",
	}
	fields := map[string]interface{}{
		"value": vibration,
	}
	writePoint("vibration", tags, fields, timestamp)
}

func PushMotorState(on bool, timestamp time.Time) {
	tags := map[string]string{
		"rack_id":   "rack01",
		"node_id":   "node02",
		"sensor_id": "motor_status",
	}
	fields := map[string]interface{}{
		"on": on,
	}
	writePoint("motor", tags, fields, timestamp)
}

func PushCoolantLeak(leak bool, timestamp time.Time) {
	tags := map[string]string{
		"rack_id":   "rack01",
		"node_id":   "node02",
		"sensor_id": "coolant_leak",
	}
	fields := map[string]interface{}{
		"leak": leak,
	}
	writePoint("coolant_leak", tags, fields, timestamp)
}

func PushFlameStatus(flame bool, timestamp time.Time) {
	tags := map[string]string{
		"rack_id":   "rack01",
		"node_id":   "node02",
		"sensor_id": "flame_detected",
	}
	fields := map[string]interface{}{
		"flame": flame,
	}
	writePoint("flame", tags, fields, timestamp)
}

func PushHumidity(humidity float64, timestamp time.Time) {
	tags := map[string]string{
		"rack_id":   "rack01",
		"node_id":   "node02",
		"sensor_id": "humidity",
	}
	fields := map[string]interface{}{
		"value": humidity,
	}
	writePoint("humidity", tags, fields, timestamp)
}

func PushSmokeStatus(smoke bool, timestamp time.Time) {
	tags := map[string]string{
		"rack_id":   "rack01",
		"node_id":   "node02",
		"sensor_id": "smoke_detected",
	}
	fields := map[string]interface{}{
		"smoke": smoke,
	}
	writePoint("smoke", tags, fields, timestamp)
}

func PushInFlowrate(inflow float64, timestamp time.Time) {
	tags := map[string]string{
		"rack_id":   "rack01",
		"node_id":   "node02",
		"sensor_id": "in_flow_rate",
	}
	fields := map[string]interface{}{
		"inflow": inflow,
	}
	writePoint("inflow", tags, fields, timestamp)
}

func PushOutFlowrate(outflow float64, timestamp time.Time) {
	tags := map[string]string{
		"rack_id":   "rack01",
		"node_id":   "node02",
		"sensor_id": "out_flow_rate",
	}
	fields := map[string]interface{}{
		"outflow": outflow,
	}
	writePoint("outflow", tags, fields, timestamp)
}

func PushCoolantTemp(temp float64, timestamp time.Time) {
	tags := map[string]string{
		"rack_id":   "rack01",
		"node_id":   "node02",
		"sensor_id": "coolant_temp",
	}
	fields := map[string]interface{}{
		"coolant_temp": temp,
	}
	writePoint("coolant_temp", tags, fields, timestamp)
}
