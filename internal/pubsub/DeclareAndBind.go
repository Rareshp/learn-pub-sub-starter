package pubsub

import (
  "log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
  ) (*amqp.Channel, amqp.Queue, error) {

  publishChannel, err := conn.Channel()
  if err != nil {
    log.Fatalf("could not creat channel: %v", err)
  }

  // transient by default
  durable := false 
  autoDelete := true 
  exclusive := true

  if simpleQueueType == 1 {
    durable = true 
    autoDelete = false 
    exclusive = false
  }

  q, err := publishChannel.QueueDeclare(
    queueName,
    durable, 
    exclusive, 
    autoDelete,
    false, // no wait
    nil,   // args
  )
  if err != nil {
    log.Fatalf("could not create Queue: %v", err)
  }

  publishChannel.QueueBind(
    queueName,
    key,
    exchange,
    false,
    nil,
  )

  return publishChannel, q, nil
}
