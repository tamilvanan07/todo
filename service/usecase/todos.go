package usecase

import (
	"context"
	"encoding/json"
	"example/todo/service/user/db"
	"time"

	"example/todo/service/user/model"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"

	"golang.org/x/crypto/bcrypt"
)

type TodoCollection struct {
	MongoCollection *mongo.Collection
}

type UserModelCollection struct {
	MongoCollection *mongo.Collection
}

type UserModel struct {
	Username *string
	Password *string
	Token    *string
}

type Response struct {
	Data    interface{}
	Success bool
	Error   string
}

type LoginResponse struct {
	Token   interface{} `json:"token"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}

var client *mongo.Client

func (svc *TodoCollection) CreateTodo(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	var todo model.TodoModel

	err := json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Panicln("bad request", err)
		res.Error = err.Error()
		return
	}
	repo := db.TodoListRepo{MongoCollection: svc.MongoCollection}

	inserId, err := repo.InsertTodoinList(&todo)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res.Success = false
		log.Fatal("Something went Wrong", err)
		return
	}

	res.Data = todo

	res.Error = "Data is Store successfully"

	res.Success = true

	w.WriteHeader(http.StatusOK)

	log.Println("Todo Added Suucessfully", inserId, todo)

}

func (svc *TodoCollection) DeletedTodo(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	var todo model.TodoModel

	err := json.NewDecoder(r.Body).Decode(&todo)

	var empId = mux.Vars(r)["id"]

	repo := db.TodoListRepo{MongoCollection: svc.MongoCollection}

	log.Printf("TaskId", empId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Panicln("bad request", err)
		res.Error = err.Error()
		return
	}

	inserId, err := repo.InsertTodoinList(&todo)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res.Success = false
		log.Fatal("Something went Wrong", err)
		return
	}

	res.Data = todo

	res.Error = "Data is Store successfully"

	res.Success = true

	w.WriteHeader(http.StatusOK)

	log.Println("Todo Added Suucessfully", inserId, todo)

}

func (svc *TodoCollection) ListTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Origin,Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,locale")
	// enableCors(&w)
	res := &Response{}

	defer json.NewEncoder(w).Encode(res)

	repo := db.TodoListRepo{MongoCollection: svc.MongoCollection}

	todoList, err := repo.FindAllList()

	if err != nil {
		res.Success = false
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("Something went Wrong", err)
		return
	}

	res.Data = todoList

	res.Success = true

	res.Error = "Success"

	w.WriteHeader(http.StatusOK)

}

func (svc *UserModelCollection) ListUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Origin,Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,locale")

	res := &Response{}

	defer json.NewEncoder(w).Encode(res)

	repo := db.UserRepo{MongoCollection: svc.MongoCollection}

	todoList, err := repo.GetUserList()

	if err != nil {
		res.Success = false
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("Something went Wrong", err)
		return
	}

	log.Printf("mongoDB collection collected,", todoList)
	res.Data = todoList

	res.Success = true

	res.Error = "Success"

	w.WriteHeader(http.StatusOK)

}

var jwtSecret = []byte("your_secret_key")

func (user *UserModelCollection) UserLoginHandler(response http.ResponseWriter, request *http.Request) {

	response.Header().Add("Content-Type", "application/json")

	var dbUser *model.Users
	json.NewDecoder(request.Body).Decode(&dbUser)

	loginRes := &LoginResponse{}
	defer json.NewEncoder(response).Encode(loginRes)

	repo := db.UserRepo{MongoCollection: user.MongoCollection}

	dbUser, err := repo.FindUserInList(*dbUser.Username)

	if err != nil {
		response.WriteHeader(http.StatusOK)

		loginRes.Success = false
		loginRes.Message = "User is not register please signup"

		// json.NewEncoder(response).Encode(loginRes)
		return
	}

	jwtToken, err := GenerateJWT(*dbUser.Username)

	if err != nil {
		loginRes.Success = false

		loginRes.Message = err.Error()

		loginRes.Token = ""

		response.WriteHeader(http.StatusInternalServerError)
		// response.Write([]byte(`{"message":"` + err.Error() + `","isSuccess": "`isSucces`",}`))
		return
	}

	loginRes.Success = true

	loginRes.Message = "Login Successfully"

	loginRes.Token = jwtToken

	response.WriteHeader(http.StatusOK)

}

func (user *UserModelCollection) RegisterFunc(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	// client = mangoInitDb()

	var dbUser *model.Users

	err := json.NewDecoder(request.Body).Decode(&dbUser)

	if err != nil {
		http.Error(response, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// Hash password before saving
	hashedPassword, err := HashPassword(dbUser.Password)
	if err != nil {
		http.Error(response, "Error hashing password", http.StatusInternalServerError)
		return
	}

	*dbUser.Password = hashedPassword

	repo := db.UserRepo{MongoCollection: user.MongoCollection}

	addUser, err := repo.InsterRegisterUser(dbUser)

	if err != nil {
		http.Error(response, "Error inserting user", http.StatusInternalServerError)
		log.Fatalln(addUser)
		return
	}

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	json.NewEncoder(response).Encode(map[string]string{"message": "User registered successfully"})

}

// HashPassword hashes a password using bcrypt
func HashPassword(password *string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with a plain password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a JWT token
func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
	})

	return token.SignedString(jwtSecret)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
