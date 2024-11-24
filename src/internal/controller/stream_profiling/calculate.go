package stream_profiling

// func CalculateStreamProfiling(logger *slog.Logger, config *entities.Config, p *kafka.Producer, c *kafka.Consumer) {
// }

// func processMessages(consumer *kafka.Consumer, producer *kafka.Producer, db DatabaseClient) {
// 	for {
// 		msg, err := consumer.ReadMessage(-1)
// 		if err != nil {
// 			// Handle consumer error
// 			continue
// 		}

// 		// Begin Kafka transaction
// 		err = producer.BeginTransaction()
// 		if err != nil {
// 			// Handle error
// 			continue
// 		}

// 		// Process the message
// 		err = processMessage(msg, db)
// 		if err != nil {
// 			// Abort transaction on error
// 			producer.AbortTransaction(nil)
// 			continue
// 		}

// 		// Send offsets to transaction
// 		offsets := []kafka.TopicPartition{{
// 			Topic:     &msg.TopicPartition.Topic,
// 			Partition: msg.TopicPartition.Partition,
// 			Offset:    msg.TopicPartition.Offset + 1,
// 		}}

// 		err = producer.SendOffsetsToTransaction(offsets, consumer.GroupMetadata(), nil)
// 		if err != nil {
// 			// Abort transaction on error
// 			producer.AbortTransaction(nil)
// 			continue
// 		}

// 		// Commit transaction
// 		err = producer.CommitTransaction(nil)
// 		if err != nil {
// 			// Handle commit error
// 			continue
// 		}
// 	}
// }
