package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sangketkit01/7-coding-test/internal/db"
	"github.com/sangketkit01/7-coding-test/internal/token"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "8090"
	gRpcPort = "50001"
	mongoUrl = "mongodb://mongo:27017"
	secretKey = "sangketketsangketkitsangketkit01"
)

var client *mongo.Client

type App struct {
	router *fiber.App
	model db.MongoClient
	jwtMaker token.Maker
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancle()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	jwtMaker, err := token.NewMaker(secretKey)
	if err !=  nil{
		log.Panic(err)
	}

	app := App{
		model: db.New(client),
		jwtMaker: jwtMaker,
	}

	app.router = app.routes()

	go app.LogsNumberOfUser()
	app.router.Listen(fmt.Sprintf(":%s", webPort))

}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "1234",
	})

	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting Mongo:", err)
		return nil, err
	}

	return conn, nil
}

func (app *App) LogsNumberOfUser(){
	for{
		users, err := app.model.ListAllUsers()
		if err != nil{
			log.Println(err)
			continue
		}

		log.Printf("Number of users:%d\n", len(users))
		time.Sleep(10 * time.Second)
	}
}
