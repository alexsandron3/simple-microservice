package controller

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

// TO-DO = Refact this code to use ENV vars and REMOTE database
func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected!")

	collection = client.Database("upvote").Collection("users")
	fmt.Println("Collection is ready")
}

func GetAllUsers() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}

	var users []primitive.M

	for cur.Next(context.Background()) {
		var user bson.M
		err := cur.Decode(&user)

		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	defer cur.Close(context.Background())

	return users

}

// TO-DO = Add feat to downvote user
func UpvoteUser(id string) primitive.M {

	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objId}
	update := bson.M{"$inc": bson.M{"votes": 1}}

	var updatedUser bson.M
	err := collection.FindOneAndUpdate(context.Background(), filter, update).Decode(&updatedUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return updatedUser
		}
		log.Fatal(err)
	}

	return updatedUser

}
