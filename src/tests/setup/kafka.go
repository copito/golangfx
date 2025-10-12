package setup

import (
	"context"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redpanda"
	"github.com/testcontainers/testcontainers-go/wait"
)

type KafkaContainer struct {
	Container *redpanda.Container
}

type KafkaConfig struct {
	User          string
	Password      string
	AdminUser     string
	AdminPassword string
}

func NewKafka(ctx context.Context, config KafkaConfig) *KafkaContainer {
	container, err := SetupKafka(ctx, config)
	if err != nil {
		log.Fatalf("failed to setup Kafka container: %v", err)
	}
	return &KafkaContainer{
		Container: container,
	}
}

func SetupKafka(ctx context.Context, config KafkaConfig) (*redpanda.Container, error) {
	kafkaContainer, err := redpanda.Run(ctx,
		"docker.redpanda.com/redpandadata/redpanda:v23.3.3",
		redpanda.WithEnableSASL(),
		redpanda.WithEnableKafkaAuthorization(),
		redpanda.WithEnableWasmTransform(),
		redpanda.WithBootstrapConfig("data_transforms_per_core_memory_reservation", 33554432),
		redpanda.WithBootstrapConfig("data_transforms_per_function_memory_limit", 16777216),
		redpanda.WithNewServiceAccount(config.User, config.Password),
		redpanda.WithSuperusers(config.AdminUser, config.AdminPassword),
		redpanda.WithEnableSchemaRegistryHTTPBasicAuth(),
		redpanda.WithAutoCreateTopics(),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForListeningPort("9092/tcp"),
				wait.ForLog("Ready"),
			),
		),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, err
	}

	return kafkaContainer, nil
}

func (s *KafkaContainer) Teardown(ctx context.Context) error {
	err := testcontainers.TerminateContainer(s.Container)
	return err
}
