package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

  defer rabbitConn.Close()
  log.Println("Connected to RabbitMQ!")

	// start listening for messages

	// create consumer

	// watch the queue and consume messages
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// Don't continue unitl rabbit is ready
	for {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			fmt.Println("Failed to connect to RabbitMQ", err)
			counts++
		} else {
			connection = conn
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off for", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
