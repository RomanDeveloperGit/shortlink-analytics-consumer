package model

import "time"

const (
	LinkCreatedMetaEvent string = "link.created"
	LinkVisitedMetaEvent string = "link.visited"
)

type Meta struct {
	Event   string `json:"event"`
	EventId string `json:"eventId"`
}

type MessageEnvelope struct {
	Meta Meta `json:"meta"`
}

type Link struct {
	ID          int       `json:"id"`
	ShortCode   string    `json:"shortCode"`
	FullURL     string    `json:"fullUrl"`
	AuthorEmail *string   `json:"authorEmail"`
	ExpiresAt   time.Time `json:"expiresAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type LinkCreatedMessage struct {
	Meta Meta `json:"meta"`
	Link Link `json:"link"`
}

func (lcm *LinkCreatedMessage) ConvertToLinkCreatedEvent() *LinkCreatedEvent {
	return &LinkCreatedEvent{
		EventID:   lcm.Meta.EventId,
		ID:        lcm.Link.ID,
		ShortCode: lcm.Link.ShortCode,
		FullURL:   lcm.Link.FullURL,
		ExpiresAt: lcm.Link.ExpiresAt,
		CreatedAt: lcm.Link.CreatedAt,
	}
}

type LinkVisitedMessage struct {
	Meta      Meta  `json:"meta"`
	LinkID    int   `json:"linkId"`
	Timestamp int64 `json:"timestamp"`
}

func (lvm *LinkVisitedMessage) ConvertToLinkVisitedEvent() *LinkVisitedEvent {
	return &LinkVisitedEvent{
		EventID:   lvm.Meta.EventId,
		LinkID:    lvm.LinkID,
		Timestamp: lvm.Timestamp,
	}
}
