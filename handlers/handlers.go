package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/renu-ramesh/robot-apocalypse-docker/helpers"
	"github.com/renu-ramesh/robot-apocalypse-docker/models"
	"github.com/renu-ramesh/robot-apocalypse-docker/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handlers struct {
	client *mongo.Client
	ctx    context.Context
	cancel context.CancelFunc
}

func New() Handlers {
	var env models.EnvVariables
	err := envconfig.Process("robot_apocalypse", &env)
	if err != nil {
		log.Fatal("unable to initialize environment variables", err.Error())
	}
	return Handlers{}
}

func (h Handlers) CreateNewsurvivors(w http.ResponseWriter, r *http.Request) {

	var env models.EnvVariables
	err := envconfig.Process("robot_apocalypse", &env)
	if err != nil {
		log.Fatal("unable to initialize environment variables", err.Error())
	}

	h.client, h.ctx, h.cancel, err = mongodb.MongoDBconnect(env.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer mongodb.MongoDBclose(h.client, h.ctx, h.cancel)

	var survivor_data models.Survivor
	reqBody, _ := ioutil.ReadAll(r.Body)
	err1 := json.Unmarshal(reqBody, &survivor_data)

	//check if survivor exist in DB
	filter := bson.M{
		"Id": bson.M{"$eq": survivor_data.Id},
	}
	//Fetch total document count
	count, err := mongodb.MongoDBCountDocuments(h.client, h.ctx, env.DatabaseName, env.Collection, filter)
	if err != nil {
		fmt.Println(err)
	}
	if count > 0 {
		response := models.Response{
			Message: "Survivor already exist",
			Success: false,
			Data:    nil,
			Error:   err,
		}
		fmt.Fprintf(w, "%+v", helpers.JSON_Marshell(response))
	} else {

		var document interface{}

		document = bson.D{
			{"Id", survivor_data.Id},
			{"Name", survivor_data.Name},
			{"Age", survivor_data.Age},
			{"Gender", survivor_data.Gender},
			{"Location", survivor_data.Location},
			{"Resource", survivor_data.Resource},
			{"Status", 0},
		}
		insertResult, err := mongodb.MongoDBinsertOne(h.client, h.ctx, env.DatabaseName, env.Collection, document)
		if err != nil {
			fmt.Fprintf(w, "%+v", err1)
		}
		response := models.Response{
			Message: "Survivor data inserted successfully",
			Success: true,
			Data:    insertResult.InsertedID,
			Error:   err,
		}
		fmt.Fprintf(w, "%+v", helpers.JSON_Marshell(response))
	}

}

// Function to modify survivors Location
func (h Handlers) UpdateSurvivorsLocation(w http.ResponseWriter, r *http.Request) {

	var env models.EnvVariables
	err := envconfig.Process("robot_apocalypse", &env)
	if err != nil {
		log.Fatal("unable to initialize environment variables", err.Error())
	}
	h.client, h.ctx, h.cancel, err = mongodb.MongoDBconnect(env.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer mongodb.MongoDBclose(h.client, h.ctx, h.cancel)

	var survivor models.Survivor
	reqBody, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(reqBody, &survivor)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
	}

	filter := bson.M{
		"Id": bson.M{"$eq": survivor.Id},
	}
	update := bson.M{
		"$set": bson.M{"Location": survivor.Location},
	}
	result, err := mongodb.MongoDBUpdateOne(h.client, h.ctx, env.DatabaseName, env.Collection, filter, update)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
	}
	response := models.Response{
		Message: "Survivor Location Modified successfully",
		Success: true,
		Data:    result.ModifiedCount,
		Error:   err,
	}
	fmt.Fprintf(w, "%+v", helpers.JSON_Marshell(response))
}

//Function to Update data of infected survivor
func (h Handlers) UpdateInfection(w http.ResponseWriter, r *http.Request) {

	var env models.EnvVariables
	err := envconfig.Process("robot_apocalypse", &env)
	if err != nil {
		log.Fatal("unable to initialize environment variables", err.Error())
	}
	h.client, h.ctx, h.cancel, err = mongodb.MongoDBconnect(env.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer mongodb.MongoDBclose(h.client, h.ctx, h.cancel)

	var survivor models.Survivor
	reqBody, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(reqBody, &survivor)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
	}

	//Update status
	filter := bson.M{
		"Id": bson.M{"$eq": survivor.Id},
	}
	update := bson.M{
		"$inc": bson.M{"Status": 1},
	}
	result, err := mongodb.MongoDBUpdateOne(h.client, h.ctx, env.DatabaseName, env.Collection, filter, update)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
	}

	response := models.Response{
		Message: "Survivor data Modified successfully",
		Success: true,
		Data:    result.ModifiedCount,
		Error:   err,
	}
	fmt.Fprintf(w, "%+v", helpers.JSON_Marshell(response))
}

// Function to find Percentage of infected/non-infected survivors
func (h Handlers) PercentageSpecification(w http.ResponseWriter, r *http.Request) {

	var env models.EnvVariables
	err := envconfig.Process("robot_apocalypse", &env)
	if err != nil {
		log.Fatal("unable to initialize environment variables", err.Error())
	}

	h.client, h.ctx, h.cancel, err = mongodb.MongoDBconnect(env.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer mongodb.MongoDBclose(h.client, h.ctx, h.cancel)

	var total_count, count int64
	filter := bson.M{
		"Id": bson.M{"$ne": ""},
	}
	//Fetch total document count
	total_count, err = mongodb.MongoDBCountDocuments(h.client, h.ctx, env.DatabaseName, env.Collection, filter)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
		fmt.Println(err)
	}

	vars := mux.Vars(r)
	key := vars["spec"]

	//check the request specification
	switch key {
	case "infected":
		filter = bson.M{
			"Status": bson.M{"$gt": 2},
		}
	case "non-infected":
		filter = bson.M{
			"Status": bson.M{"$lt": 3},
		}
	default:
		fmt.Fprintf(w, "%+v", "404 page not found")
	}
	count, err = mongodb.MongoDBCountDocuments(h.client, h.ctx, env.DatabaseName, env.Collection, filter)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
		fmt.Println(err)
	}
	// Calculate Percentage
	percentage := int((float64(count) / float64(total_count)) * 100)
	percentageStr := fmt.Sprintf("%d", percentage)

	response := models.Response{
		Message: key + " Percentage",
		Success: true,
		Data:    percentageStr + "%",
		Error:   err,
	}
	fmt.Fprintf(w, "%+v", helpers.JSON_Marshell(response))

}

//Function to List infected/non-infected survivors
func (h Handlers) ListSurvivors(w http.ResponseWriter, r *http.Request) {

	var env models.EnvVariables
	err := envconfig.Process("robot_apocalypse", &env)
	if err != nil {
		log.Fatal("unable to initialize environment variables", err.Error())
	}
	h.client, h.ctx, h.cancel, err = mongodb.MongoDBconnect(env.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// Free the resource when main function is  returned
	defer mongodb.MongoDBclose(h.client, h.ctx, h.cancel)

	// create a filter an option of type interface,
	// that stores bjson objects.
	var filter, option interface{}
	vars := mux.Vars(r)
	key := vars["spec"]
	switch key {
	case "infected":
		filter = bson.M{
			"Status": bson.M{"$gt": 2},
		}
	case "non-infected":
		filter = bson.M{
			"Status": bson.M{"$lt": 3},
		}
	default:
		fmt.Fprintf(w, "%+v", "404 page not found")
	}
	option = bson.D{{"_id", 0}}
	cursor, err := mongodb.MongoDBquery(h.client, h.ctx, env.DatabaseName, env.Collection, filter, option)
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err := cursor.All(h.ctx, &results); err != nil {
		panic(err)
	}

	for _, doc := range results {
		resJson, _ := json.Marshal(doc)
		fmt.Fprintf(w, "%+v\n", string(resJson))
	}

}

// Function to Connect to the Robot CPU system
func (h Handlers) ListAllRobots(w http.ResponseWriter, r *http.Request) {

	var env models.EnvVariables
	err := envconfig.Process("robot_apocalypse", &env)
	if err != nil {
		log.Fatal("unable to initialize environment variables", err.Error())
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", env.RoboCpuUrl, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	response := models.Response{
		Message: "Connect to the Robot CPU system",
		Success: true,
		Data:    string(bodyBytes),
		Error:   err,
	}

	// fmt.Fprintf(w, "%+v", helpers.JSON_Marshell(response))
	fmt.Fprintf(w, "%+v", response)
}
