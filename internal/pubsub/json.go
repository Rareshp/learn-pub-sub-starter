package pubsub

import (
  "fmt"
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {

  fmt.Println("Printing message...", val)

  m, err := json.Marshal(val)
  if err != nil { 
    return err 
  }

  // needed struct for Publishing 
  msg := amqp.Publishing{
    // DeliveryMode: amqp.Persistent,
    // Timestamp:    time.Now(),
    ContentType:  "application/json",
    Body:         m,
  }

  return ch.PublishWithContext(context.Background(), exchange, key, false, false, msg)
}
