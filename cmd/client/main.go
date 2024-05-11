package main

import (
  "os"
  "os/signal"
  "syscall"
	"fmt"
	"log"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const connectionString = "amqp://guest:guest@localhost:5672/"

func ctrlC() {
  fmt.Println("\nShutting down")
  os.Exit(1)
}

func main() {
	fmt.Println("Starting Peril client...")

  // create a new connection for rabbitmq
  connection, err := amqp.Dial(connectionString)
  // ensure the connection is closed when the program exits
  if err != nil {
    log.Print(err.Error())
  }
  defer connection.Close()

  fmt.Println("Connection to RabbitMQ was successful")

  userName, err := gamelogic.ClientWelcome()
  if err != nil {
    log.Printf("getting username error: %v", err)
  }

  publishChannel, q, err := pubsub.DeclareAndBind(
    connection,
    routing.ExchangePerilDirect,
    fmt.Sprintf("%v.%v", routing.PauseKey, userName),
    routing.PauseKey,
    2, // transient
  )
  if err != nil {
    log.Fatalf("could not create channel or q, %v:", err)
  }

  log.Println(publishChannel, q)

  // handle Ctrl+C to quit
  signalChan := make(chan os.Signal, 1) // the int is for buffering
  signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
  <-signalChan
  ctrlC()
}
