package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

//Client instance
var DB *mongo.Client = ConnectDB()

//getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("admin_panel").Collection(collectionName)
	return collection
}

func GetUser(ctx context.Context, email string) (storedEmail string, password string, err error) {
	coll := GetCollection(DB, "users")
	var user bson.M
	err = coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		return
	}

	return user["email"].(string), user["password"].(string), nil
}

func GetUserByID(ctx context.Context, id primitive.ObjectID) (userDetails bson.M, err error) {
	coll := GetCollection(DB, "user_details")
	var user bson.M
	err = coll.FindOne(ctx, bson.M{"user_id": id}).Decode(&user)

	if err != nil {
		return
	}

	return user, nil
}
