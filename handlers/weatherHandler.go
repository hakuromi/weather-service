package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/hakuromi/weather-service/sideAPI/geocoding"
	"github.com/hakuromi/weather-service/sideAPI/openmeteo"
)

var (
	mu          = &sync.Mutex{}              // составная операция получения данных
	weatherData = make(map[string][]metrics) // хранилище для данных
)

type metrics struct { // данные о температуре
	Timestamp   string
	Temperature float64
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cityName := r.URL.Path[1:]
	if cityName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("City name is required!"))
		return
	}

	mu.Lock()
	defer mu.Unlock()
	fmt.Println("Requested city:", cityName)
	geoResp, err := geocoding.GetCoords(cityName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error getting coordinates of the city"))
		return
	}

	weatherResp, err := openmeteo.GetTemp(geoResp.Latitude, geoResp.Longitude)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error getting temperature"))
		return
	}

	weatherData[cityName] = append(weatherData[cityName], metrics{
		Timestamp:   weatherResp.Current.Time,
		Temperature: weatherResp.Current.Temperature2m,
	})

	response, err := json.Marshal(weatherData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error forming response"))
		return
	}
	w.Write(response)
}
