package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RomanDeveloperGit/shortlink-analytics-consumer/internal/config"
	"github.com/RomanDeveloperGit/shortlink-analytics-consumer/internal/infra/db"
	"github.com/RomanDeveloperGit/shortlink-analytics-consumer/internal/logger"
	linkRepository "github.com/RomanDeveloperGit/shortlink-analytics-consumer/internal/repository/link"
	"github.com/RomanDeveloperGit/shortlink-analytics-consumer/internal/transport/kafkaconsumergroup"
	linkHandler "github.com/RomanDeveloperGit/shortlink-analytics-consumer/internal/transport/kafkaconsumergroup/handler/link"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	log := logger.MustSetup(&logger.Options{
		File:           os.Stdout,
		Env:            cfg.Env,
		ServiceName:    cfg.ServiceName,
		ServiceVersion: cfg.ServiceVersion,
	})

	dbClient := db.MustConnect(&db.Options{
		ClientID:        cfg.ServiceName,
		ClientVersion:   cfg.ServiceVersion,
		Host:            cfg.ClickHouse.Host,
		Port:            cfg.ClickHouse.Port,
		User:            cfg.ClickHouse.User,
		Password:        cfg.ClickHouse.Password,
		DBName:          cfg.ClickHouse.DB,
		MaxIdleConns:    cfg.ClickHouse.MaxIdleConns,
		MaxOpenConns:    cfg.ClickHouse.MaxOpenConns,
		ConnMaxLifetime: cfg.ClickHouse.ConnMaxLifetime,
		DialTimeout:     cfg.ClickHouse.DialTimeout,
		PingTimeout:     cfg.ClickHouse.PingTimeout,
		Logger:          log,
	})

	linkRepo := linkRepository.NewRepository(dbClient.DB)
	linkH := linkHandler.NewHandler(log, linkRepo)

	consumerGroupCtx, cancelConsumerGroupCtx := context.WithCancel(ctx)
	defer cancelConsumerGroupCtx()

	for range cfg.Kafka.ConsumersCount {
		consumerGroupClient := kafkaconsumergroup.NewKafkaConsumerGroup(&kafkaconsumergroup.Options{
			ServiceName:           cfg.ServiceName,
			Addrs:                 cfg.Kafka.Brokers,
			User:                  cfg.Kafka.User,
			Password:              cfg.Kafka.Password,
			MaxHandleAttempts:     cfg.Kafka.MaxConsumeHandleAttempts,
			MaxConsumeProcessTime: cfg.Kafka.MaxConsumeProcessTime,
			LinkHandler:           linkH,
			Logger:                log,
		})

		go consumerGroupClient.Consume(consumerGroupCtx)
	}

	time.Sleep(1 * time.Millisecond)

	log.Info("server started")

	<-ctx.Done()

	log.Info("shutting down server")

	// Останавливаем слушатель ОС сигналов, чтобы последующий второй CTRL+C мгновенно завершал работу.
	// Это работает, потому что снова произойдет INTERRUPT, на который никто не подписан(т.к. stop убирает эту подписку как раз)
	// А по умолчанию без подписки программа сразу же завершает своё исполнение
	// Сделано на случай, если нужно срочно остановить сервер, не дожидаясь полного graceful shutdown (часто в local-режиме)
	stop()

	cancelConsumerGroupCtx()
	dbClient.Close()

	log.Info("server stopped")
}
