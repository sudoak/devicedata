// main.go
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
)

// Device Data
type Device struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	DeviceID  string             `json:"device_id,omitempty" bson:"device_id,omitempty"`
	E1        float32            `json:"e1,omitempty" bson:"e1,omitempty"`
	E2        float32            `json:"e2,omitempty" bson:"e2,omitempty"`
	E3        float32            `json:"e3,omitempty" bson:"e3,omitempty"`
	E4        float32            `json:"e4,omitempty" bson:"e4,omitempty"`
	E5        float32            `json:"e5,omitempty" bson:"e5,omitempty"`
	Date      string             `json:"date" bson:"date"`
	Time      string             `json:"time" bson:"time"`
	TimeStamp string             `json:"timestamp" bson:"timestamp"`
}

// type TestDevice struct {
// 	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
// 	DeviceID  string             `json:"device_id,omitempty" bson:"device_id,omitempty"`
// 	E1        float32            `json:"e1,omitempty" bson:"e1,omitempty"`
// 	E2        float32            `json:"e2,omitempty" bson:"e2,omitempty"`
// 	E3        float32            `json:"e3,omitempty" bson:"e3,omitempty"`
// 	E4        float32            `json:"e4,omitempty" bson:"e4,omitempty"`
// 	E5        float32            `json:"e5,omitempty" bson:"e5,omitempty"`
// 	Date      string             `json:"date" bson:"date"`
// 	Time      string             `json:"time" bson:"time"`
// 	TimeStamp string             `json:"timestamp" bson:"timestamp"`
// }

var client *mongo.Client

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

// Mongo Routes Functions
func createPacketData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var device Device
	_ = json.NewDecoder(r.Body).Decode(&device)
	device.TimeStamp = time.Now().String()
	fmt.Println(device)
	collection := client.Database("askak").Collection("atom")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, device)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(result)
}

func GetDeviceData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var devices []Device
	collection := client.Database("askak").Collection("atom")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var device Device
		cursor.Decode(&device)
		devices = append(devices, device)
	}
	fmt.Println(devices)
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(devices)
}

// Exit
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/device/v1/data", GetDeviceData).Methods("GET")
	myRouter.HandleFunc("/device/v1/data", createPacketData).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://sudoak:sudoak1@ds237337.mlab.com:37337/askak?retryWrites=false")
	client, _ = mongo.Connect(ctx, clientOptions)

	fmt.Println("Connected to MongoDB!")
	handleRequests()
}
