package persist

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User is the structure for handling users in our application
type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserName string             `json:"username,omitempty" bson:"username,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
}

//GetUserByID finds a user in mongoDB
func GetUserByID(id string) (User, error) {
	mongoID, err := primitive.ObjectIDFromHex(id)
	var u User
	if err != nil {
		log.Println("Invalid hex")
		u.UserName = "Invalid hex id: " + id
		return u, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userCollection := collection("aqueduct", "users")
	if err := userCollection.FindOne(ctx, bson.M{"_id": mongoID}).Decode(&u); err != nil {
		u.UserName = "No user with id: " + id
		return u, err
	}
	//Clean up the return user to not send back the password
	u.Password = ""
	return u, nil
}

//GetUserByUsername finds a user in mongoDB
func GetUserByUsername(username string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userCollection := collection("aqueduct", "users")
	var u User
	if err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&u); err != nil {
		log.Println("Unable to find user")
		return u, err
	}
	//Clean up the return user to not send back the password
	u.Password = ""
	return u, nil
}

//InsertUser creates a new document for a user in mongoDB
func InsertUser(u User) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var returnUser User

	userCollection := collection("aqueduct", "users")
	userResult, err := userCollection.InsertOne(ctx, u)
	if err != nil {
		log.Println(err)
		return u, err
	}
	log.Println(userResult.InsertedID.(primitive.ObjectID).Hex())
	returnUser.ID = userResult.InsertedID.(primitive.ObjectID)
	return returnUser, nil

}
