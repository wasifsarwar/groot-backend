package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"os"
	"reflect"

	"github.com/nlopes/slack"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var connectionString = os.Getenv("MONGO_STRING")
var api = slack.New(os.Getenv("SLACK_TOKEN"))

const dbName = "slackdb"
const collName = "userlist"

// Declare Context type object for managing multiple API requests
var ctx context.Context

var collection *mongo.Collection

//Initmongo connection
func Initmongo() {

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to MongoDB!")

	collection = client.Database(dbName).Collection(collName)
	ctx, _ = context.WithTimeout(context.Background(), 15*time.Second)

	InsertNewOrUpsert()
}

//InsertNewOrUpsert ...
func InsertNewOrUpsert() {

	mongoPayLoad := GetAllMongoUsers()
	slackPayLoad := SlackAPIUsers(api)

	// creating an array of ids from the slackpayload
	var mongoIDlist []string
	for _, mongoUser := range mongoPayLoad {
		mongoIDlist = append(mongoIDlist, mongoUser.SlackID)
	}

	var slackIDlist []string
	for _, slackUser := range slackPayLoad {
		slackIDlist = append(slackIDlist, slackUser.SlackID)
	}

	for i := 0; i < len(slackIDlist); i++ {
		slackID := slackIDlist[i]
		if !itemExists(mongoIDlist, slackID) {
			for _, slackUser := range slackPayLoad {
				if slackUser.SlackID == slackID && slackID != "USLACKBOT" {
					InsertOne(ctx, slackUser, collection)
				}
			}
		} else {
			for _, slackUser := range slackPayLoad {
				if slackUser.SlackID == slackID && slackID != "USLACKBOT" {
					filter := bson.M{"slackid": slackID}

					updatedData := bson.M{"name": slackUser.Name, "email": slackUser.Email, "deletestatus": slackUser.DeleteStatus}
					UpdateOne(filter, updatedData)
				}
			}
		}
	}
}

//UpdateOne updates a user document in Mongo
func UpdateOne(updatedData interface{}, filter bson.M) {
	setUpdate := bson.D{{Key: "$set", Value: updatedData}}
	result, updaterr := collection.UpdateOne(ctx, filter, setUpdate)
	if updaterr != nil {
		log.Fatal("Error on updating: ", updaterr)
	} else {
		fmt.Println("matched count: ", result.MatchedCount, " modified count: ", result.ModifiedCount, " upserted count: ", result.UpsertedCount, " upserted id: ", result.UpsertedID)
	}
}

//helper to see if slack id exists in mongo array
func itemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

//InsertOne inserts one user document in Mongo
func InsertOne(ctx context.Context, mem MemoryStruct, collection *mongo.Collection) {

	result, inserterr := collection.InsertOne(ctx, mem)
	if inserterr != nil {
		log.Println(inserterr)
	} else {
		fmt.Println("InsertOne() result type: ", reflect.TypeOf(result))
		fmt.Println("InsertOne() API result:", result)

		// get the inserted ID string
		newID := result.InsertedID
		fmt.Println("InsertOne() newID:", newID)
		fmt.Println("InsertOne() newID type:", reflect.TypeOf(newID))
	}
}

/*
	SlackAPIUsers connects to slack and collects all users in the api
	Stores info in in-memory struct AllStruct

*/
func SlackAPIUsers(api *slack.Client) []MemoryStruct {

	user, err := api.GetUsers()
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	var result []MemoryStruct

	for i := 0; i < len(user); i++ {
		var mem MemoryStruct
		mem.SlackID = user[i].ID
		mem.Name = user[i].RealName

		mem.DeleteStatus = user[i].Deleted
		mem.Email = user[i].Profile.Email
		result = append(result, mem)
	}

	return result
}

/*
 * This section gets all users from mongoDB and creates an endpoint for router to ingest for frontend
 */

//GetAllUser returns a struct full of users existing in mongodb
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	payload := GetAllMongoUsers()
	fmt.Println("endpoint hit. These are users being returned: ", payload)
	json.NewEncoder(w).Encode(payload)
}

//GetAllMongoUsers connects to mongo and returns an array full of userlist data
func GetAllMongoUsers() []MemoryStruct {
	curr, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var results []MemoryStruct
	for curr.Next(context.Background()) {
		var user MemoryStruct
		err := curr.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, user)
	}

	if err := curr.Err(); err != nil {
		log.Fatal(err)
	}

	curr.Close(context.Background())
	return results
}
