package main

import (
	"bufio"
	"fmt"
	"net/http"
)

func main() {

	resp, err := http.Get("http://dataservice.accuweather.com/currentconditions/v1/329828?apikey=aqKf60jG6EIL9LqZKSf5KGnRAH1prGPe&language=en-us&details=false")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
