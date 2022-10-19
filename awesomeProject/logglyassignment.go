package main

import (
	_ "context"
	"encoding/json"
	_ "errors"
	"fmt"
	_ "fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jamespearly/loggly"
	"github.com/joho/godotenv"
	"io"
	_ "io/ioutil"
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	_ "os"
	_ "strconv"
)
import _ "github.com/jamespearly/loggly"
import _ "github.com/joho/godotenv"
import _ "github.com/aws/aws-sdk-go/aws"
import _ "github.com/aws/aws-sdk-go/aws/session"
import _ "github.com/aws/aws-sdk-go/service/dynamodb"
import _ "github.com/aws/aws-sdk-go/aws/credentials"
import _ "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
import _ "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
import _ "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
import _ "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

// type of JSON response struct

type malUser struct {
	CurrentCondition []CurrentCondition `json:"CurrentCondition"`
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

func main() {

	response, err := http.Get(`http://dataservice.accuweather.com/currentconditions/v1/329828?apikey=aqKf60jG6EIL9LqZKSf5KGnRAH1prGPe&language=en-us&details=false`)
	if err != nil {
		log.Fatal(err)
	}

	//load env var
	err1 := godotenv.Load(".env")
	if err1 != nil {
		log.Fatal(err1)
	}

	//FIRST PRINT
	//log.Print(response.Body)
	bytes, errRead := io.ReadAll(response.Body)

	defer func() {
		e := response.Body.Close()
		if e != nil {
			log.Fatal(e)
		}
	}()

	if errRead != nil {
		log.Fatal(errRead)
	}
	//SECOND PRINT
	//log.Print(string(bytes))

	var currentCondition []CurrentCondition

	errUnmarshal := json.Unmarshal(bytes, &currentCondition)
	if errUnmarshal != nil {
		log.Fatal(errUnmarshal)
	}

	b, err := json.MarshalIndent(currentCondition, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(string(b))

	//THIRD PRINT
	//log.Printf("%+v", b)

	var tag string
	tag = "LogglyAssignment"

	client := loggly.New(tag)
	r := response.StatusCode
	s := len(bytes)
	err2 := client.Send("info", fmt.Sprintf("Response Code:%+v Response Size: %+v,", r, s))
	if err2 != nil {
		log.Fatal(err2)
	}

	//time.Sleep(1 * time.Hour)

	//DYNAMODB STUFF
	sess := session.Must(session.NewSessionWithOptions(session.Options{

		/*Config: aws.Config{
			Region: aws.String("us-east-1"),
		},*/
		SharedConfigState: session.SharedConfigEnable,
	}))

	//dynamodb client creation
	svc := dynamodb.New(sess)

	// Create table Movies
	tableName := "rkanin-accuweather"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("LocalObservationDateTime"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("WeatherIcon"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("LocalObservationDateTime"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("WeatherIcon"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err = svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	//log.Print("Created the table", tableName)

	av, erra := dynamodbattribute.MarshalMap(currentCondition[0])
	if erra != nil {
		log.Fatalf("Got error marshalling new movie item: %s", erra)
	}

	//log.Print(av)

	input2 := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	_, errb := svc.PutItem(input2)
	if errb != nil {
		log.Fatalf("Got error calling PutItem: %s", errb)
	}

}
