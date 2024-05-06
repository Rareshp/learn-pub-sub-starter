package main

import (
	"fmt"
	"log"
  "os"
  "os/signal"
  "syscall"
	amqp "github.com/rabbitmq/amqp091-go"
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

  // handle Ctrl+C to quit 
  signalChan := make(chan os.Signal, 1) // the int is for buffering
  signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
  <-signalChan
  ctrlC()
}
