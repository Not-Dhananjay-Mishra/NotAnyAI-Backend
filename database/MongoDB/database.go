package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"server/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var alldata = make(map[string]string)

type Data struct {
	ID       any    `bson:"_id,omitempty"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

var dbcoll *mongo.Collection

func DBinit() {
	uri := utils.MONGODB_CLUSTER
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// verify connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	coll := client.Database("nothing").Collection("user")
	dbcoll = coll
	log.Println(utils.Green("Database initiated sucessfully"))

}

func GetUserPassword(username string) (string, error) {
	var result Data
	filter := bson.M{"username": username}
	err := dbcoll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println("Error finding user:", err)
	} else {
		fmt.Printf("User found: %+v\n", result.ID)
	}
	//log.Println(result)
	if result.Password == "" {
		return "", errors.New("not exist")
	}

	return result.Password, nil
}

func AddUser(username string, password string) error {
	filter := bson.M{"username": username}
	var existing Data
	dbcoll.FindOne(context.TODO(), filter).Decode(&existing)
	//log.Println(existing)
	if existing.Username != "" {
		log.Println("user already exists")
		return errors.New("user already exists")
	}

	alldata[username] = password
	_, err := dbcoll.InsertOne(context.TODO(), bson.M{
		"username": username,
		"password": password,
	})
	if err != nil {
		log.Println("Error writing in DB:", err)
		return err
	}

	return nil
}

func TestPrintAllUser() {
	for i, j := range alldata {
		fmt.Println(i, j)
	}
}
