package link

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"

	"github.com/RomanDeveloperGit/shortlink-analytics-consumer/internal/model"
)

type repository struct {
	db clickhouse.Conn
}

func NewRepository(db clickhouse.Conn) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateLinkCreatedEvent(
	ctx context.Context,
	link *model.LinkCreatedEvent,
) error {
	q := `
		INSERT INTO link_created_event (event_id, link_id, short_code, full_url, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	fmt.Println(link)
	return r.db.Exec(
		ctx,
		q,
		link.EventID,
		link.ID,
		link.ShortCode,
		link.FullURL,
		link.ExpiresAt,
		link.CreatedAt,
	)
}

func (r *repository) CreateLinkVisitedEvent(
	ctx context.Context,
	link *model.LinkVisitedEvent,
) error {
	q := `
		INSERT INTO link_visited_event (event_id, link_id, timestamp)
		VALUES (?, ?, ?)
	`

	return r.db.Exec(
		ctx,
		q,
		link.EventID,
		link.LinkID,
		link.Timestamp,
	)
}
