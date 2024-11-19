package fake_publisher

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/copito/runner/internal/entities"
)

func GenerateFakeData(logger *slog.Logger, conf *entities.Config, p *kafka.Producer) {
	logger.Info("Generating fake data...")

	messageCount := 0

	sourceData := entities.ChangeDataCaptureSource{
		Version: 1,
		App:     "service_1",
	}

	logger.Info(
		"using kafka connection",
		slog.String("kafka.server", conf.Kafka.Server),
		slog.String("kafka.topic", conf.Kafka.ChangeDataCaptureTopicExample),
	)

	ctx := context.Background()
	err := p.InitTransactions(ctx)
	if err != nil {
		logger.Error("Failed to initialize transactions", slog.Any("err", err))
		return
	}

	for {
		// Random delay between 1 to 5 seconds
		delay := time.Duration(rand.Intn(5)+1) * time.Second
		time.Sleep(delay)

		// Random choice between lemon and strawberry string
		colB := "lemon"
		if rand.Intn(2) == 1 {
			colB = "strawberry"
		}

		// Generate the message payload
		payload := entities.ChangeDataCaptureEventPayload{
			Op:     "i",
			Before: nil,
			After: &entities.ChangeDataCaptureMessage{
				ID:        messageCount + 1,
				ColA:      rand.Intn(100),
				ColB:      colB,
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			Source: sourceData,
		}

		messageBytes, err := json.Marshal(payload)
		if err != nil {
			logger.Error("Failed to marshal message")
			continue
		}

		// Send the message to Kafka
		topic := conf.Kafka.ChangeDataCaptureTopicExample

		err = p.BeginTransaction()
		if err != nil {
			logger.Error("Failed to begin transaction", slog.Any("err", err))
			return
		}

		message := &kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: messageBytes,
		}

		err = p.Produce(message, nil)
		if err != nil {
			logger.Error("Failed to send message", slog.Any("err", err))
			err = p.AbortTransaction(ctx)
			if err != nil {
				logger.Error("Failed to abort message", slog.Any("err", err))
			}

			return
		}

		err = p.CommitTransaction(ctx)
		if err != nil {
			logger.Error("Failed to commit transaction", slog.Any("err", err))
			return
		}

		messageCount += 1
		logger.Info("Message sent", slog.Int("messageCount", messageCount))

		// Emit thousands of messages
		if messageCount >= 1000 {
			break
		}
	}

	logger.Info("Stopped sending example messages...")
}
