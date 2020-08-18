package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)
var ClientOptions = options.Client().ApplyURI(Url)
var Client,_ = mongo.Connect(context.TODO(),ClientOptions)
func init() {
	err := Client.Ping(context.TODO(),nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}

func CloseConn() {
	err :=  Client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection shut down")
}
