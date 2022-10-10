package main

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	"github.com/jamespearly/loggly"
	"github.com/joho/godotenv"
	"io"
	_ "io/ioutil"
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	_ "os"
)
import _ "github.com/jamespearly/loggly"
import _ "github.com/joho/godotenv"

// type of JSON response struct

type malUser struct {
	CurrentCondition []CurrentCondition `json:"CurrentCondition"`
}
type CurrentCondition struct {
	LocalObservationDayTime string      `json:"LocalObservationDateTime"`
	EpochTime               int         `json:"EpochTime"`
	WeatherText             string      `json:"WeatherText"`
	WeatherIcon             int         `json:"WeatherIcon"`
	HasPrecipitation        bool        `json:"HasPrecipitation"`
	PrecipitationType       *string     `json:"PrecipitationType"`
	IsDayTime               bool        `json:"IsDayTime"`
	Temperature             Temperature `json:"Temperature"`
	MobileLink              string      `json:"MobileLink"`
	Link                    string      `json:"Link"`
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

	log.Print(response.Body)
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
	log.Print(string(bytes))

	var currentCondition []CurrentCondition

	errUnmarshal := json.Unmarshal(bytes, &currentCondition)
	if errUnmarshal != nil {
		log.Fatal(errUnmarshal)
	}

	log.Printf("%+v", currentCondition)

	var tag string
	tag = "LogglyAssignment"

	client := loggly.New(tag)
	r := response.StatusCode
	s := len(bytes)
	err2 := client.Send("info", fmt.Sprintf("Response Code:%+v Response Size: %+v,", r, s))
	if err2 != nil {
		log.Fatal(err2)
	}
}
