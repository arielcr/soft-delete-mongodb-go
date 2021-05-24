package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/arielcr/soft-delete-mongodb-go/entities"
	"github.com/arielcr/soft-delete-mongodb-go/repository"
	"github.com/caarlos0/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Read our configuration environment variables
	config := entities.MongoDBConfig{}
	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}

	// Create our context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Instantiate our MongoDB Client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI))
	if err != nil {
		log.Fatal(err)
	}

	// Create our repository
	repo := repository.NewMongoDbRepository(client, config)

	// Create a new User
	user := entities.User{
		Name:  "John Doe",
		Email: "john@doe.com",
	}

	userID, err := repo.CreateUser(ctx, user)
	if err != nil {
		fmt.Println(err)
	}

	// Get the User
	firstGet, err := repo.GetUser(ctx, userID)
	if err != nil {
		fmt.Println(err)
	}

	// Format nicely
	_ = formatOutput(firstGet)

	// Delete the User
	err = repo.DeleteUser(ctx, userID)
	if err != nil {
		fmt.Println(err)
	}

	// Get the User (again)
	_, err = repo.GetUser(ctx, userID)
	if err != nil {
		fmt.Println(err)
	}

}

func formatOutput(data interface{}) error {
	format, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(format))
	return nil
}
