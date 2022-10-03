package main

import (
	"fmt"
	loggly "github.com/jamespearly/loggly"
	"github.com/joho/godotenv"
	"log"
)

func main() {

	err1 := godotenv.Load(".env")

	if err1 != nil {
		log.Fatal(err1)
	}

	var tag string
	tag = "My-Go-Demo"

	// Instantiate the client
	client := loggly.New(tag)

	// Valid EchoSend (message echoed to console and no error returned)
	err := client.EchoSend("info", "Good morning!")
	fmt.Println("err:", err)

	// Valid Send (no error returned)
	err = client.Send("error", "Good morning! No echo.")
	fmt.Println("err:", err)

	// Invalid EchoSend -- message level error
	err = client.EchoSend("blah", "blah")
	fmt.Println("err:", err)

}
