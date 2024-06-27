package main

import (
	"fmt"
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
  "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
  "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)

const connectionString = "amqp://guest:guest@localhost:5672/"

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
    log.Fatalf("could not create channel: %v", err)
  } 

  gamelogic.PrintServerHelp()

  gameLogsChannel, q, err := pubsub.DeclareAndBind(
    connection,
    routing.ExchangePerilTopic,
    routing.ExchangeGameLogQName,
    routing.ExchangeGameLogSlug,
    1, // durable
  )
  if err != nil {
    log.Fatalf("could not create channel or q, %v:", err)
  }
  log.Println(gameLogsChannel, q)

  exit := 0 
  for exit == 0 {
    userInput := gamelogic.GetInput()
    switch userInput[0] {
      case "": 
        continue
      case "pause":
        message(publishChannel, true)
      case "resume":
        message(publishChannel, false)
      case "quit":
        gamelogic.PrintQuit()
        exit = 1
      default:
        fmt.Printf("did not understand command")
        gamelogic.PrintServerHelp()
        continue
    }
  }
}

func message(publishChannel *amqp.Channel, b bool) {
  if b {
    log.Println("sending pause message")
  } else {
    log.Println("sending resume message")
  }
  err := pubsub.PublishJSON(
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
  fmt.Println("Pause message sent")
}
