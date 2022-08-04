package sqlbublic

import (
	"context"
	configpackage "golang_chatbot/config"
	"log"
	"time"

	model "golang_chatbot/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDB() (client *mongo.Client) {
	config, err := configpackage.InitConfig()

	log.Print(config.LineChannelSecret)
	URI := "mongodb://" + config.MongoAccount + ":" + config.MongoPassword + "@127.0.0.1:27017/?authSource=admin"
	clientOptions := options.Client().ApplyURI(URI)
	var ctx = context.TODO()

	// Connect to MongoDB
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client

}

func SaveMessage(client *mongo.Client, userid string, message string, time time.Time) {

	collection := client.Database("linemessage").Collection("linemessage")

	insertmodel := model.LineMessage{
		UserID:  userid,
		Message: message,
		Time:    time,
	}

	_, err := collection.InsertOne(context.TODO(), insertmodel)
	if err != nil {
		log.Fatal(err)
	}

}
func GetMessages(client *mongo.Client, userid string) (messages []bson.D) {

	collection := client.Database("linemessage").Collection("linemessage")

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

	return messages

}
