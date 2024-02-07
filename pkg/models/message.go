package models

import "time"

type MQTTMessage struct {
	ID                   int       `pg:",pk"`
	SessionID            int       `json:"session_id"`
	EnergyDeliveredInKWh float64   `json:"energy_delivered_in_kWh"`
	DurationInSeconds    int       `json:"duration_in_seconds"`
	SessionCostInCents   int       `json:"session_cost_in_cents"`
	Timestamp            time.Time `json:"timestamp"`
}
