package main

import (
	"fmt"
	"net/http"

	"github.com/hakuromi/weather-service/handlers"
)

const (
	httpPort = ":3000"
)

func main() {
	http.HandleFunc("/", handlers.WeatherHandler)
	fmt.Println("Starting server on port", httpPort)
	err := http.ListenAndServe(httpPort, nil)
	if err != nil {
		panic(err)
	}
}
