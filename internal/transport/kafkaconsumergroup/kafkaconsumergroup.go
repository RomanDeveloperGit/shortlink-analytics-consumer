package kafkaconsumergroup

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
)

const topic = "links.events"

type Kafkaconsumergroup struct {
	ConsumerGroup         sarama.ConsumerGroup
	maxHandleAttempts     int
	maxConsumeProcessTime time.Duration
	linkHandler           linkHandler
	logger                *slog.Logger
}

type Options struct {
	ServiceName           string
	Addrs                 []string
	User                  string
	Password              string
	MaxHandleAttempts     int
	LinkHandler           linkHandler
	Logger                *slog.Logger
	MaxConsumeProcessTime time.Duration
}

func NewKafkaConsumerGroup(opts *Options) *Kafkaconsumergroup {
	cfg := sarama.NewConfig()

	cfg.ClientID = opts.ServiceName

	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.AutoCommit.Enable = false
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	groupID := fmt.Sprintf("%s-group", opts.ServiceName)
	consumerGroup, err := sarama.NewConsumerGroup(
		opts.Addrs,
		groupID,
		cfg,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create kafka consumer: %w", err))
	}

	opts.Logger.Debug("kafka consumer created",
		slog.String("addrs", fmt.Sprintf("%v", opts.Addrs)),
		slog.String("group", groupID),
	)

	return &Kafkaconsumergroup{
		ConsumerGroup:         consumerGroup,
		maxHandleAttempts:     opts.MaxHandleAttempts,
		maxConsumeProcessTime: opts.MaxConsumeProcessTime,
		linkHandler:           opts.LinkHandler,
		logger:                opts.Logger,
	}
}

func (kc *Kafkaconsumergroup) Consume(ctx context.Context) {
	for {
		err := kc.ConsumerGroup.Consume(ctx, []string{topic}, &handler{
			maxAttempts:           kc.maxHandleAttempts,
			maxConsumeProcessTime: kc.maxConsumeProcessTime,
			linkHandler:           kc.linkHandler,
			logger:                kc.logger,
		})
		if err != nil {
			kc.logger.Error("failed to consume",
				slog.String("topic", topic),
				slog.String("error", err.Error()),
			)
		}

		if ctx.Err() != nil {
			kc.logger.Debug("consumer stopped gracefully")
			return
		}
	}
}
