package models

import (
	"time"
)

type Click struct {
	ID         string    `json:"id" db:"id"`
	LinkID     string    `json:"link_id" db:"link_id"`
	ClickedAt  time.Time `json:"clicked_at" db:"clicked_at"`
	IPAddress  string    `json:"ip_address" db:"ip_address"`
	UserAgent  string    `json:"user_agent" db:"user_agent"`
	Referrer   string    `json:"referrer" db:"referrer"`
	Country    string    `json:"country" db:"country"`
	DeviceType string    `json:"device_type" db:"device_type"`
}

type Analytics struct {
	TotalClicks       int             `json:"total_clicks"`
	DailyClicks       []DailyClick    `json:"daily_clicks"`
	CountryBreakdown  []CountryStats  `json:"country_breakdown"`
	DeviceBreakdown   []DeviceStats   `json:"device_breakdown"`
	ReferrerBreakdown []ReferrerStats `json:"referrer_breakdown"`
}

type DailyClick struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type CountryStats struct {
	Country string `json:"country"`
	Count   int    `json:"count"`
}

type DeviceStats struct {
	DeviceType string `json:"device_type"`
	Count      int    `json:"count"`
}

type ReferrerStats struct {
	Referrer string `json:"referrer"`
	Count    int    `json:"count"`
}

type AnalyticsSummary struct {
	LinkID         string `json:"link_id"`
	TotalClicks    int    `json:"total_clicks"`
	UniqueVisitors int    `json:"unique_visitors"`
}
