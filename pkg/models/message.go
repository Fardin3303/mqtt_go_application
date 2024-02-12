package models

import "time"

type MQTTMessage struct {
	SessionID            int       `pg:",pk"`
	EnergyDeliveredInKWh float64   `json:"energy_delivered_in_kWh"`
	DurationInSeconds    int       `json:"duration_in_seconds"`
	SessionCostInCents   int       `json:"session_cost_in_cents"`
	Timestamp            time.Time `json:"timestamp"`
}
