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
func GetMessages(userid string) (messages []bson.D) {

	collection := Client.Database("linemessage").Collection("linemessage")

	fmt.Println("----------------------")

	projection := bson.D{
		{"message", 1},
		{"_id", 0},
	}

	cursor, err := collection.Find(
		context.TODO(),
		bson.D{
			{"userid", userid},
		},
		options.Find().SetProjection(projection),
	)

	if err = cursor.All(context.TODO(), &messages); err != nil {
		log.Fatal(err)
	}
	fmt.Println(messages)

	return messages

}
