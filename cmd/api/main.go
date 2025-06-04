package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sangketkit01/7-coding-test/internal/config"
	"github.com/sangketkit01/7-coding-test/internal/db"
	"github.com/sangketkit01/7-coding-test/internal/token"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "8090"
	gRpcPort = "50001"
)

var client *mongo.Client

type App struct {
	router   *fiber.App
	model    db.MongoClient
	jwtMaker token.Maker
	config   *config.Config
}

func init() {
	if err := godotenv.Load("../../.env.local"); err != nil {
		log.Println("ไม่พบ .env.local ลองโหลด .env.production แทน")
		if err := godotenv.Load("../../.env.production"); err != nil {
			log.Println("ไม่พบ .env.production ด้วย")
		}
	}
}

func main() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "local"
	}

	config, err := config.NewConfig("../../", env)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("environment:", config.Environment)

	mongoClient, err := connectToMongo(config.MongoUrl, config.MongoUsername, config.MongoPassword)
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

	jwtMaker, err := token.NewMaker(config.SecretKey)
	if err != nil {
		log.Panic(err)
	}

	app := App{
		model:    db.New(client),
		jwtMaker: jwtMaker,
		config:   config,
	}

	app.router = app.routes()

	go app.LogsNumberOfUser()
	go app.gRPCListen()

	app.router.Listen(fmt.Sprintf(":%s", webPort))

}

func connectToMongo(url, username, password string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(url)
	clientOptions.SetAuth(options.Credential{
		Username: username,
		Password: password,
	})

	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting Mongo:", err)
		return nil, err
	}

	return conn, nil
}

func (app *App) LogsNumberOfUser() {
	for {
		users, err := app.model.ListAllUsers()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Number of users:%d\n", len(users))
		time.Sleep(10 * time.Second)
	}
}
