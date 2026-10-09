package model

import "time"

type LinkCreatedEvent struct {
	EventID   string    `db:"event_id"   json:"eventId"`
	ID        int       `db:"id"         json:"id"`
	ShortCode string    `db:"short_code" json:"shortCode"`
	FullURL   string    `db:"full_url"   json:"fullUrl"`
	ExpiresAt time.Time `db:"expires_at" json:"expiresAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type LinkVisitedEvent struct {
	EventID   string `db:"event_id"   json:"eventId"`
	LinkID    int    `db:"link_id"   json:"linkId"`
	Timestamp int64  `db:"timestamp" json:"timestamp"`
}
