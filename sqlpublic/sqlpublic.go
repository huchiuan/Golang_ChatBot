package sqlbublic

import (
	"context"
	"fmt"
	configpackage "golang_chatbot/config"
	"log"
	"time"

	model "golang_chatbot/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectToDB() {
	config, err := configpackage.InitConfig()

	log.Print(config.LineChannelSecret)
	URI := "mongodb://" + config.MongoAccount + ":" + config.MongoPassword + "@127.0.0.1:27017/?authSource=admin"
	clientOptions := options.Client().ApplyURI(URI)
	var ctx = context.TODO()
	// Connect to MongoDB
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// defer client.Disconnect(ctx)

}

func List() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	databases, err := Client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
}

func SaveMessage(userid string, message string, time time.Time) {

	collection := Client.Database("linemessage").Collection("linemessage")

	insertmodel := model.LineMessage{
		UserID:  userid,
		Message: message,
		Time:    time,
	}

	insertResult, err := collection.InsertOne(context.TODO(), insertmodel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

}
