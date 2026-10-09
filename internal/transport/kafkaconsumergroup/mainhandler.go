package kafkaconsumergroup

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"

	"github.com/RomanDeveloperGit/shortlink-analytics-consumer/internal/model"
)

type linkHandler interface {
	ProcessLinkCreatedMessage(
		ctx context.Context,
		linkCreatedMessage *model.LinkCreatedMessage,
	) error
	ProcessLinkVisitedMessage(
		ctx context.Context,
		linkVisitedMessage *model.LinkVisitedMessage,
	) error
}

type handler struct {
	maxAttempts           int
	maxConsumeProcessTime time.Duration
	linkHandler           linkHandler
	logger                *slog.Logger
}

func (h *handler) Setup(session sarama.ConsumerGroupSession) error {
	h.logger.Debug("consumer started")

	return nil
}

func (h *handler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *handler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		ctx, cancel := context.WithTimeout(context.Background(), h.maxConsumeProcessTime)
		defer cancel()

		h.logger.Info("message received",
			slog.Int("partition", int(msg.Partition)),
			slog.Int64("offset", msg.Offset),
			slog.String("value", string(msg.Value)),
		)

		err := h.processWithRetry(ctx, msg)
		if err != nil {
			h.logger.Error("failed to process message after retries",
				slog.Int("partition", int(msg.Partition)),
				slog.Int64("offset", msg.Offset),
				slog.String("value", string(msg.Value)),
				slog.Any("error", err),
			)

			continue
		}

		session.MarkMessage(msg, "")
		session.Commit()

		h.logger.Info("message processed",
			slog.Int("partition", int(msg.Partition)),
			slog.Int64("offset", msg.Offset),
			slog.String("value", string(msg.Value)),
		)
	}

	return nil
}

func (h *handler) processWithRetry(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var lastErr error

	for attempt := 1; attempt <= h.maxAttempts; attempt++ {
		err := h.processMessage(ctx, msg)
		if err == nil {
			return nil
		}

		h.logger.Error("failed to process message",
			slog.Int("attempt", attempt),
			slog.Any("error", err),
		)

		// экспоненциальный backoff: 100ms, 200ms, 400ms...
		time.Sleep(time.Duration(100*(1<<(attempt-1))) * time.Millisecond)
	}

	return lastErr
}

func (h *handler) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var messageEnvelope model.MessageEnvelope

	err := json.Unmarshal(msg.Value, &messageEnvelope)
	if err != nil {
		return fmt.Errorf("failed to unmarshal messageEnvelope: %w", err)
	}

	switch messageEnvelope.Meta.Event {
	case model.LinkCreatedMetaEvent:
		var linkCreatedMessage model.LinkCreatedMessage

		err = json.Unmarshal(msg.Value, &linkCreatedMessage)
		if err != nil {
			return fmt.Errorf("failed to unmarshal link created message: %w", err)
		}

		err = h.linkHandler.ProcessLinkCreatedMessage(ctx, &linkCreatedMessage)
		if err != nil {
			return fmt.Errorf("failed to process link created message: %w", err)
		}
	case model.LinkVisitedMetaEvent:
		var linkVisitedMessage model.LinkVisitedMessage

		err = json.Unmarshal(msg.Value, &linkVisitedMessage)
		if err != nil {
			return fmt.Errorf("failed to unmarshal link visited message: %w", err)
		}

		err = h.linkHandler.ProcessLinkVisitedMessage(ctx, &linkVisitedMessage)
		if err != nil {
			return fmt.Errorf("failed to process link visited message: %w", err)
		}
	default:
		return fmt.Errorf(
			"unknown event: %s with eventId: %s",
			messageEnvelope.Meta.Event,
			messageEnvelope.Meta.EventId,
		)
	}

	return nil
}
