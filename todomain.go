package main

import (
	"context"
	"example/todo/service/usecase"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// // type response struct {
// // 	message string
// // }

// type TASK struct {
// 	ID     int
// 	TASK   string
// 	ISDONE bool
// }

// // type TaskResponse struct {
// // 	task    TASK
// // 	message response
// // }

// var task = []TASK{
// 	{ID: 1, TASK: "DO GOLANG PROJECT", ISDONE: false},
// }

// func getTask(context *gin.Context) {
// 	context.IndentedJSON(http.StatusAccepted, task)
// }

// func addTask(context *gin.Context) {
// 	var newTask TASK

// 	if err := context.BindJSON(&newTask); err != nil {
// 		context.IndentedJSON(http.StatusBadRequest, TASK{
// 			ID:     0,
// 			TASK:   "PLEASE ADD CORRECT DATA",
// 			ISDONE: false,
// 		})
// 		return
// 	}

// 	if newTask.ID == 0 {
// 		context.IndentedJSON(http.StatusNotModified, TASK{
// 			ID:     0,
// 			TASK:   "PLEASE ADD CORRECT DATA",
// 			ISDONE: false,
// 		})
// 		return
// 	}
// 	task = append(task, newTask)
// 	context.IndentedJSON(http.StatusCreated, newTask)
// }

var mongoClient *mongo.Client

func mangoInitDb() *mongo.Client {

	mangoClieny, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODP_URI")))

	if err != nil {
		log.Fatal("unable to connnect the DB", err)
	}

	err = mangoClieny.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatal("PING faild", err)
	}

	return mangoClieny

}

func init() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("unable to load ENV", err)
	}

	log.Println("loaded ENV")

	log.Printf("mongo connected")

	mangoInitDb()

}

func main() {

	mongoClient = mangoInitDb()

	defer mongoClient.Disconnect(context.Background())

	log.Printf("mongoDB disconnected")

	collection := mongoClient.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_NAME"))
	userCollection := mongoClient.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_USERNAME"))

	log.Printf("mongoDB collection collected")
	todoService := usecase.TodoCollection{MongoCollection: collection}
	userService := usecase.UserModelCollection{MongoCollection: userCollection}
	r := mux.NewRouter()
	r.Use(recoverMiddleware)
	r.HandleFunc("/todos", todoListHandler).Methods(http.MethodGet)
	r.HandleFunc("/todo", todoService.CreateTodo).Methods(http.MethodPost)
	r.HandleFunc("/todoList", todoService.ListTodo).Methods(http.MethodGet)
	r.HandleFunc("/register", userService.RegisterFunc).Methods(http.MethodPost)
	r.HandleFunc("/login", userService.UserLoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/user-list", userService.ListUser).Methods(http.MethodGet)
	log.Println("Running Server 8000")

	http.ListenAndServe(":8000", r)

}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Recovered from panic: %v", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func todoListHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("running .."))
}
