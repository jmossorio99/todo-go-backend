package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Todo struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
}

var client, err = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://joseosso:littlebird0926*@cluster0.ofgit.mongodb.net"))

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func GetTodos(res http.ResponseWriter, req *http.Request) {
	fmt.Println("GetTodos: called")
	res.Header().Add("content-type", "application/json")
	var todos []Todo
	collection := client.Database("todoVueDB").Collection("todos")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var todo Todo
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}
	if err := cursor.Err(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(res).Encode(todos)
}

func GetTodo(res http.ResponseWriter, req *http.Request) {
	fmt.Println("GetTodo: called")
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var todo Todo
	collection := client.Database("todoVueDB").Collection("todos")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Todo{ID: id}).Decode(&todo)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(res).Encode(todo)
}

func CreateTodo(res http.ResponseWriter, req *http.Request) {
	fmt.Println("createNewTodo: called")
	res.Header().Add("content-type", "application/json")
	var todo Todo
	json.NewDecoder(req.Body).Decode(&todo)
	collection := client.Database("todoVueDB").Collection("todos")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, todo)
	json.NewEncoder(res).Encode(result)
}

func DeleteTodo(res http.ResponseWriter, req *http.Request) {
	fmt.Println("DeleteTodo: called")
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("todoVueDB").Collection("todos")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.DeleteOne(ctx, Todo{ID: id})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(res).Encode(result)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/todos", GetTodos).Methods("GET")
	myRouter.HandleFunc("/todos/{id}", GetTodo).Methods("GET")
	myRouter.HandleFunc("/todos", CreateTodo).Methods("POST")
	myRouter.HandleFunc("/todos/{id}", DeleteTodo).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	// //client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://joseosso:littlebird0926*@cluster0.ofgit.mongodb.net"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	handleRequests()
}
