package link

import (
	"context"
	"log/slog"

	"github.com/RomanDeveloperGit/shortlink-analytics-consumer/internal/model"
)

type linkRepository interface {
	CreateLinkCreatedEvent(
		ctx context.Context,
		link *model.LinkCreatedEvent,
	) error
	CreateLinkVisitedEvent(
		ctx context.Context,
		link *model.LinkVisitedEvent,
	) error
}

type handler struct {
	repository linkRepository
	logger     *slog.Logger
}

func NewHandler(logger *slog.Logger, repository linkRepository) *handler {
	return &handler{
		logger: logger,
		repository: repository,
	}
}

func (h *handler) ProcessLinkCreatedMessage(
	ctx context.Context,
	linkCreatedMessage *model.LinkCreatedMessage,
) error {
	return h.repository.CreateLinkCreatedEvent(ctx, linkCreatedMessage.ConvertToLinkCreatedEvent())
}

func (h *handler) ProcessLinkVisitedMessage(
	ctx context.Context,
	linkVisitedMessage *model.LinkVisitedMessage,
) error {
	return h.repository.CreateLinkVisitedEvent(ctx, linkVisitedMessage.ConvertToLinkVisitedEvent())
}
