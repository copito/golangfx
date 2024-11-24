package modules

import (
	"context"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/copito/runner/src/internal/entities"
	"go.uber.org/fx"
)

type KafkaProducerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *slog.Logger
	Config    *entities.Config
}

type KafkaProducerResults struct {
	fx.Out

	Producer *kafka.Producer
}

func NewKafkaProducer(params KafkaProducerParams) (KafkaProducerResults, error) {
	params.Logger.Info("setting up Kafka Producer...")

	kafkaConfig := params.Config.Kafka

	config := &kafka.ConfigMap{
		"bootstrap.servers":  kafkaConfig.Server,
		"transactional.id":   "unique_transactional_id",
		"enable.idempotence": true,
	}

	params.Logger.Info(
		"using kafka producer configs",
		slog.String("server", kafkaConfig.Server),
		slog.String("transactional_id", "unique_transactional_id"),
		slog.Bool("idempotence", true),
	)

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return KafkaProducerResults{}, err
	}

	// Use fx lifecycle hooks to manage the database connection
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			params.Logger.Info("Kafka Connection established...")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Closing producer kafka connection")
			producer.Flush(100)
			producer.Close()
			return nil
		},
	})

	return KafkaProducerResults{Producer: producer}, nil
}

var KafkaProducerModule = fx.Provide(NewKafkaProducer)
