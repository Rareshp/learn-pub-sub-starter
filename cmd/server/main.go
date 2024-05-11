package main

import (
	"fmt"
	"log"
  "os"
	amqp "github.com/rabbitmq/amqp091-go"
  "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
  "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

const connectionString = "amqp://guest:guest@localhost:5672/"

func cleanup() {
  fmt.Println("\nShutting down")
}

func ctrlC() {
  cleanup()
  os.Exit(1)
}

func main() {
	fmt.Println("Starting Peril server...")

  // create a new connection for rabbitmq
  connection, err := amqp.Dial(connectionString)
  // ensure the connection is closed when the program exits

  if err != nil {
    log.Print(err.Error())
  }
  defer connection.Close()

  fmt.Println("Connection to RabbitMQ was successful")

  publishChannel, err := connection.Channel()
  if err != nil {
    log.Fatalf("could not creat channel: %v", err)
  } 

  fmt.Println("About to publish JSON...")

  err = pubsub.PublishJSON(
    publishChannel, 
    routing.ExchangePerilDirect,
    routing.PauseKey,
    routing.PlayingState{
      IsPaused: true,
    },
  )
  if err != nil {
    log.Printf("could not publish message: %v", err)
  }
  fmt.Println("Puase message sent")
}
