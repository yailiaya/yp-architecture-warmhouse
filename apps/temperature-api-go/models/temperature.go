package models

import "time"

// TemperatureResponse mirrors services.TemperatureResponse in the Go smart_home
// app. The JSON tags must stay identical so the temperature service that calls
// this API can decode the payload.
type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}
