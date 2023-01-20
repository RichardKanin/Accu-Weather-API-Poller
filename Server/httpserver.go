package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jamespearly/loggly"
	"log"
	"net/http"
	"time"
)

type User struct {
	SystemTime time.Time `json:"SystemTime"`
	HTTPStatus int       `json:"ResponseCode"`
}

type CurrentCondition struct {
	LocalObservationDateTime string      `json:"LocalObservationDateTime"`
	EpochTime                int         `json:"EpochTime"`
	WeatherText              string      `json:"WeatherText"`
	WeatherIcon              int         `json:"WeatherIcon"`
	HasPrecipitation         bool        `json:"HasPrecipitation"`
	PrecipitationType        *string     `json:"PrecipitationType"`
	IsDayTime                bool        `json:"IsDayTime"`
	Temperature              Temperature `json:"Temperature"`
	MobileLink               string      `json:"MobileLink"`
	Link                     string      `json:"Link"`
}

type Temperature struct {
	Metric   Metric   `json:"Metric"`
	Imperial Imperial `json:"Imperial"`
}

type Metric struct {
	Value    float64 `json:"Value"`
	Unit     string  `json:"Unit"`
	UnitType int     `json:"UnitType"`
}

type Imperial struct {
	Value    float64 `json:"Value"`
	Unit     string  `json:"Unit"`
	UnitType int     `json:"UnitType"`
}
type Status struct {
	TableName   string `json:"Table"`
	RecordCount *int64 `json:"RecordCount"`
}

func createDynamoDBClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	//dynamodb client creation
	svc := dynamodb.New(sess)

	return svc

}
func helloStatus(w http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {

		w.Header().Set("Content-Type", "application/json")

		a := http.StatusOK

		user := User{
			SystemTime: time.Now(),
			HTTPStatus: a,
		}

		err := json.NewEncoder(w).Encode(user)
		if err != nil {
			return
		}

		var tag string
		tag = "HTTPServer on EC2"

		client := loggly.New(tag)
		methodType := req.Method
		sourceIP := req.RemoteAddr
		requestPath := req.URL.Path
		httpStatus := a
		err2 := client.Send("info", fmt.Sprintf("Method Type: %v Source IP Address: %v Request Path: %v HTTP Status: %v", methodType, sourceIP, requestPath, httpStatus))
		if err2 != nil {
			log.Fatal(err2)
		}

	} else {
		http.Error(w, "Invalid request method.", 405)
	}

}

func allEndpoint(w http.ResponseWriter, r *http.Request) {
	svc := createDynamoDBClient()

	out, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("rkanin-accuweather"),
	})
	if err != nil {
		panic(err)
	}

	var currentCondition []CurrentCondition

	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &currentCondition)

	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal, %v", err))
	}

	fmt.Fprintf(w, "%v", currentCondition)
}

func statusEndpoint(w http.ResponseWriter, r *http.Request) {
	svc := createDynamoDBClient()

	tableName := "rkanin-accuweather"

	out, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("rkanin-accuweather"),
	})
	if err != nil {
		panic(err)
	}

	status := Status{
		TableName:   tableName,
		RecordCount: out.Count,
	}

	errKa := json.NewEncoder(w).Encode(status)
	if errKa != nil {
		return
	}

}

func searchEndpoint(w http.ResponseWriter, r *http.Request) {
	svc := createDynamoDBClient()

	LocalObservationDateTime := r.FormValue("LocalObservationDateTime")
	weatherIcon := r.FormValue("WeatherIcon")

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("rkanin-accuweather"),
		Key: map[string]*dynamodb.AttributeValue{
			"LocalObservationDateTime": {
				S: aws.String(LocalObservationDateTime),
			},
			"WeatherIcon": {
				N: aws.String(weatherIcon),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	currentCondition := CurrentCondition{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &currentCondition)

	if result.Item == nil {
		http.Error(w, "INVALID REQUEST", 400)
	}

	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal, %v", err))
	}

	fmt.Fprintf(w, "%v", currentCondition)
}

func main() {

	http.HandleFunc("/rkanin/all", allEndpoint)
	http.HandleFunc("/rkanin/status", statusEndpoint)
	http.HandleFunc("/rkanin/search", searchEndpoint)

	err5 := http.ListenAndServe(":8080", nil)
	if err5 != nil {
		return
	}

}
