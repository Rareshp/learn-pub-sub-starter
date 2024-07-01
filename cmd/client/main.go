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

  gameState := gamelogic.NewGameState(userName)

  err = pubsub.SubscribeJSON(
    connection, 
    routing.ExchangePerilDirect,
    routing.PauseKey + "." + gameState.GetUsername(),
    routing.PauseKey,
    2, //transient
    handlerPause(gameState),
  )

  if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

  err = pubsub.SubscribeJSON(
    connection,
    routing.ExchangePerilTopic,
    routing.ArmyMovesPrefix + "." + gameState.GetUsername(),
		routing.ArmyMovesPrefix + ".*",
    2, // transient
    handlerMove(gameState),
  )

  if err != nil {
    log.Fatalf("could not subscribe to any moves: %v", err)
  }

  exit := 0 
  for exit == 0 {
    userInput := gamelogic.GetInput()
    switch userInput[0] {
      case "": 
        continue
      case "spawn":
        // spawn europe infantry
        gameState.CommandSpawn(userInput)
      case "move":
        // move europe 1
        mv, err := gameState.CommandMove(userInput)
        if err != nil {
          fmt.Println(err)
          continue
        }
        err = pubsub.PublishJSON(
          publishChannel,
          routing.ExchangePerilTopic,
          routing.ArmyMovesPrefix + "." + mv.Player.Username,
          mv,
        )
        if err != nil {
          fmt.Println(err)
        }
        fmt.Printf("Moved %v units to %s\n", len(mv.Units), mv.ToLocation)
      case "status":
        gameState.CommandStatus()
      case "spam":
        fmt.Println("Spamming not allowed yet!")
      case "help":
        gamelogic.PrintServerHelp()
      case "quit":
        gamelogic.PrintQuit()
        exit = 1
      default:
        fmt.Println("did not understand command")
        gamelogic.PrintServerHelp()
        continue
    }
  }

  log.Println(publishChannel, q)

  // handle Ctrl+C to quit
  signalChan := make(chan os.Signal, 1) // the int is for buffering
  signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
  <-signalChan
  ctrlC()
}
