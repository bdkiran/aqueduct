package persist

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

//InitMongoClient initializes a conection to MongoDb
func InitMongoClient(userName string, password string) {
	log.Println("Establishing Connection to mongoDB")
	atlasConnectionString := "mongoconnectionstring"
	connectionURI := fmt.Sprintf(atlasConnectionString, userName, password)
	clientOptions := options.Client().ApplyURI(connectionURI)

	var err error
	mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Printf("Unable to make a connection to mongoDb: %s", connectionURI)
	}
	//ping the server to verify connection is made.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err = mongoClient.Ping(ctx, nil); err != nil {
		log.Println("Error when pinging the mongoDB database")
	}
}

//DisconnectClient closes the client connection to mongodb,
//this ensures no lingering connections
func DisconnectClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient.Disconnect(ctx)
}

func collection(databaseName string, collectionName string) *mongo.Collection {
	collection := mongoClient.Database(databaseName).Collection(collectionName)
	return collection
}
