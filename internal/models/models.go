package models

import "time"

type ActivityLog struct {
	ID               int       `json:"id"`
	Tanggal          string    `json:"tanggal"`
	JobTitle         string    `json:"job_title"`
	PIC              string    `json:"pic"`
	Application      string    `json:"application"`
	Label            string    `json:"label"`
	OldValueText     string    `json:"old_value_text"`
	OldValueImageURL string    `json:"old_value_image_url"`
	NewValueText     string    `json:"new_value_text"`
	NewValueImageURL string    `json:"new_value_image_url"`
	Status           string    `json:"status"`
	Category         string    `json:"category"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

var ValidJobTitle = map[string]bool{"DBA": true, "DCO": true, "NETWORK": true}
var ValidStatus = map[string]bool{"Open": true, "Process": true, "Hold": true, "Done": true}
var ValidCategory = map[string]bool{"Change": true, "Daily": true, "Incident": true}
