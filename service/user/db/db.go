package db

import (
	"context"
	"example/todo/service/user/model"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoListRepo struct {
	MongoCollection *mongo.Collection
}
type UserRepo struct {
	MongoCollection *mongo.Collection
}

func (r *TodoListRepo) InsertTodoinList(todo *model.TodoModel) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), todo)
	if err != nil {
		log.Fatal("unable to insert data", err)
		return nil, err
	}
	return result, nil

}

func (r *TodoListRepo) DeleteTodoinList(todoId string) (int64, error) {
	result, err := r.MongoCollection.DeleteOne(context.Background(), bson.D{{Key: "todoid", Value: todoId}})
	if err != nil {
		log.Fatal("unable to delete data", err)
		return 0, err
	}
	return result.DeletedCount, nil

}

func (r *TodoListRepo) FindAllList() ([]model.TodoModel, error) {

	result, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var emps []model.TodoModel

	err = result.All(context.Background(), &emps)

	if err != nil {
		return nil, fmt.Errorf("result not available")
	}
	return emps, nil

}

func (r *UserRepo) InsterRegisterUser(user *model.Users) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal("unable to insert user", err)
		return nil, err
	}
	return result, nil

}

func (r *UserRepo) FindUserInList(email string) (*model.Users, error) {
	var user *model.Users

	err := r.MongoCollection.FindOne(context.Background(), bson.M{"username": email}).Decode(&user)

	log.Printf(" user Name,", user)
	if err != nil {

		return nil, err
	}

	return user, nil

}

func (r *UserRepo) GetUserList() ([]model.Users, error) {
	result, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {

		return nil, err
	}

	var emps []model.Users

	err = result.All(context.Background(), &emps)

	if err != nil {
		return nil, fmt.Errorf("result not available")
	}

	return emps, nil

}
