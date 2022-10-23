package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/radugaf/RelantinV3/logger/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	port     = "80"
	rpcPort  = "5001"
	grpcPort = "50001"
	mongoURL = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct {
	Models models.Models
}

func main() {
	mongoClient, err := connectMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// Create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Close the connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Create a new config
	config := Config{
		Models: models.New(client),
	}

	// Start the server
	// go config.serve()

	log.Println("Starting server on port: ", port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: config.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// func (config *Config) serve() {
// 	server := &http.Server{
// 		Addr:    fmt.Sprintf(":%s", port),
// 		Handler: config.routes(),
// 	}

// 	err := server.ListenAndServe()
// 	if err != nil {
// 		log.Panic(err)
// 	}
// }

func connectMongo() (*mongo.Client, error) {
	// Create connectin options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// Connect
	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connectiong to mongo: ", err)
		return nil, err
	}

	return conn, nil
}
