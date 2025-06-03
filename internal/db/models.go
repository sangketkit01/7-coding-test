package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/sangketkit01/7-coding-test/internal/util"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
	CreatedAt time.Time              `bson:"created_at" json:"created_at"`
}

func New(mongo *mongo.Client) MongoClient {
	client = mongo

	collection := client.Database("users").Collection("users")
	if err := createEmailUniqueIndex(collection); err != nil {
		log.Println("failed to create unique index on email:", err)
	}

	return User{}
}

func createEmailUniqueIndex(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	return err
}


func (u User) LoginUser() (*User, error) {
	collection := client.Database("users").Collection("users")

	// ctx := context.TODO()

	var foundUser User
	fmt.Println("Decode target type:", reflect.TypeOf(foundUser))

	err := collection.FindOne(context.TODO(), bson.M{"email": u.Email}).Decode(&foundUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("user not found")
			return nil, errors.New("invalid email or password")
		}

		log.Println("error finding user:", err)
		return nil, err
	}

	if err := util.CheckPassword(foundUser.Password, u.Password); err != nil {
		log.Println("password mismatch")
		return nil, errors.New("invalid email or password")
	}

	log.Println("user logged in successfully:", u.Email)

	return &foundUser, nil
}

func (u User) Insert(user User) error {
	collection := client.Database("users").Collection("users")

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		log.Println("failed to hashed password:", err)
		return err
	}

	_, err = collection.InsertOne(context.TODO(), User{
		Name:      user.Name,
		Email:     user.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Println("email already exists")
			return errors.New("email already exists")
		}

		log.Println("failed to insert user:", err)
		return err
	}

	log.Println("insert user successfully")

	return nil
}

func (u User) FetchUserByID(id string) (*User, error) {
	collection := client.Database("users").Collection("users")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("invalid object id:", err)
		return nil, err
	}

	var user User
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("user not found")
			return nil, errors.New("user not found")
		}

		log.Println("error finding user by id:", err)
		return nil, err
	}

	return &user, nil
}

func (u User) ListAllUsers() ([]*User, error) {
	collection := client.Database("users").Collection("users")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("failed to fetch users:", err)
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var users []*User
	for cursor.Next(context.TODO()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			log.Println("failed to decode user:", err)
			return nil, err
		}

		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		log.Println("cursor error:", err)
		return nil, err
	}

	return users, nil
}

func (u User) UpdateUser() error {
	collection := client.Database("users").Collection("users")

	objectID, err := primitive.ObjectIDFromHex(u.ID.Hex())
	if err != nil {
		log.Println("invalid object id:", err)
		return errors.New("invalid user ID")
	}

	var current User
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&current)
	if err != nil {
		log.Println("failed to fetch current user:", err)
		return err
	}

	name := current.Name
	if u.Name != "" {
		name = u.Name
	}
	email := current.Email
	if u.Email != "" {
		email = u.Email
	}

	update := bson.M{
		"$set": bson.M{
			"name":  name,
			"email": email,
		},
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		update,
	)

	if err != nil {
		log.Println("failed to update user:", err)
		return err
	}

	log.Println("user updated successfully")
	return nil
}

func (u User) DeleteUser() error {
	collection := client.Database("users").Collection("users")

	objectID, err := primitive.ObjectIDFromHex(u.ID.Hex())
	if err != nil {
		log.Println("invalid object id:", err)
		return errors.New("invalid user ID")
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		log.Println("failed to delete user:", err)
		return err
	}

	log.Println("user deleted successfully")
	return nil
}

func (u User) GetUserByEmail(email string) (*User, error) {
	collection := client.Database("users").Collection("users")

	var user User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("user not found by email")
			return nil, errors.New("user not found")
		}

		log.Println("error finding user by email:", err)
		return nil, err
	}

	return &user, nil
}
