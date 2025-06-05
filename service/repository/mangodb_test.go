package repository

import (
	"context"
	"example/todo/service/user/db"
	"example/todo/service/user/model"
	"log"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MangoDbConnection() *mongo.Client {

	mongoTestClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://tamilvananvanan120:7EEsCJevBt8OoWnz@go-todo.4ywwcgu.mongodb.net/?retryWrites=true&w=majority&appName=go-todo"))

	if err != nil {
		log.Fatal("Error while connecting the Db", err)
	}

	log.Printf("Connection Done Successful")

	err = mongoTestClient.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatal("Error")
	}
	return mongoTestClient

}

func TestMongoDBConnection(t *testing.T) {
	mongoTestClient := MangoDbConnection()

	defer mongoTestClient.Disconnect(context.Background())

	cell := mongoTestClient.Database("todo").Collection("todo-task")

	todoRepo := db.TodoListRepo{MongoCollection: cell}

	t.Run("insert Todo 1st", func(t *testing.T) {
		emp := model.TodoModel{
			TodoId:      "2",
			Title:       "Today Ram temple going",
			Description: "This is second Step",
			Time:        time.Now().String(),
		}

		result, err := todoRepo.InsertTodoinList(&emp)

		if err != nil {

			t.Fatal("insert 1 opertation fail", err)

		}
		t.Fatal("insert 1 opertation success", result)

	})
}
