package model

type Stats struct {
	Requests      int64 `json:"requests"`
	Keys          int   `json:"keys"`
	UptimeSeconds int64 `json:"uptime_seconds"`
}
