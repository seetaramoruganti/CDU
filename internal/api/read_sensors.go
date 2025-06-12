package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// Load InfluxDB config from environment
var (
	influxURL    = os.Getenv("INFLUX_URL")
	influxToken  = os.Getenv("INFLUX_TOKEN")
	influxOrg    = os.Getenv("INFLUX_ORG")
	influxBucket = os.Getenv("INFLUX_BUCKET")
)

// SensorPoint represents one timestamped value
type SensorPoint struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

// HistoryHandler streams historical data for one metric
// GET /api/history?metric=temperature&from=2025-06-08T00:00:00Z&to=2025-06-09T00:00:00Z
func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse query params
	q := r.URL.Query()
	metric := q.Get("metric")
	if metric == "" {
		http.Error(w, "metric is required", http.StatusBadRequest)
		return
	}
	from := q.Get("from")
	to := q.Get("to")
	if from == "" || to == "" {
		http.Error(w, "from and to are required", http.StatusBadRequest)
		return
	}
	// 2. Parse timestamps
	start, err := time.Parse(time.RFC3339, from)
	if err != nil {
		http.Error(w, "invalid from format", http.StatusBadRequest)
		return
	}
	stop, err := time.Parse(time.RFC3339, to)
	if err != nil {
		http.Error(w, "invalid to format", http.StatusBadRequest)
		return
	}
	// 3. Build Flux query
	flux := fmt.Sprintf(
		`from(bucket:"%s")
		  |> range(start: %s, stop: %s)
		  |> filter(fn: (r) => r._measurement == "%s")
		  |> sort(columns: ["_time"])`,
		influxBucket,
		start.Format(time.RFC3339),
		stop.Format(time.RFC3339),
		metric,
	)
	// 4. Execute query
	client := influxdb2.NewClient(influxURL, influxToken)
	defer client.Close()
	queryAPI := client.QueryAPI(influxOrg)
	result, err := queryAPI.Query(context.Background(), flux)
	if err != nil {
		http.Error(w, fmt.Sprintf("query error: %v", err), http.StatusInternalServerError)
		return
	}
	// 5. Collect points
	points := make([]SensorPoint, 0)
	for result.Next() {
		rec := result.Record()
		val, ok := rec.Value().(float64)
		if !ok {
			continue
		}
		points = append(points, SensorPoint{Time: rec.Time(), Value: val})
	}
	if err := result.Err(); err != nil {
		http.Error(w, fmt.Sprintf("iteration error: %v", err), http.StatusInternalServerError)
		return
	}
	// 6. Write JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(points); err != nil {
		http.Error(w, fmt.Sprintf("encode error: %v", err), http.StatusInternalServerError)
	}
}
