package modules

import (
	"context"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/copito/runner/internal/entities"
	"go.uber.org/fx"
)

type KafkaConsumerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *slog.Logger
	Config    *entities.Config
}

type KafkaConsumerResults struct {
	fx.Out

	Consumer *kafka.Consumer
}

func NewKafkaConsumer(params KafkaConsumerParams) (KafkaConsumerResults, error) {
	params.Logger.Info("setting up Kafka Consumer...")

	kafkaConfig := params.Config.Kafka

	config := &kafka.ConfigMap{
		"bootstrap.servers":    kafkaConfig.Server,
		"group.id":             "groupID1234",
		"enable.auto.commit":   false,
		"isolation.level":      "read_committed",
		"auto.offset.reset":    "earliest",
		"enable.partition.eof": true,
	}

	params.Logger.Info(
		"using kafka consumer configs",
		slog.String("server", kafkaConfig.Server),
	)

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return KafkaConsumerResults{}, err
	}

	// Use fx lifecycle hooks to manage the database connection
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			params.Logger.Info("Kafka Connection Consumer established...")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Closing consumer kafka connection")
			consumer.Close()
			return nil
		},
	})

	return KafkaConsumerResults{Consumer: consumer}, nil
}

var KafkaConsumerModule = fx.Provide(NewKafkaConsumer)
