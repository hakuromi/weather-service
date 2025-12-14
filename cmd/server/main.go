package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
	"github.com/hakuromi/weather-service/http/geocoding"
	"github.com/hakuromi/weather-service/http/openmeteo"
)

const (
	httpPort = ":3000"
	city     = "moscow"
)

type Metrics struct {
	Timestamp   time.Time
	Temperature float64
}

type Storage struct {
	data map[string][]Metrics
	mu   sync.RWMutex
}

func main() {
	r := chi.NewRouter()     // роутер обрабатывет наши адреса
	r.Use(middleware.Logger) // мидлваре - для всех эндпоинтов

	storage := &Storage{
		data: make(map[string][]Metrics),
	}

	r.Get("/{city}", func(w http.ResponseWriter, r *http.Request) {
		cityName := chi.URLParam(r, "city")

		fmt.Println("Requested city:", cityName)

		storage.mu.RLock()
		defer storage.mu.RUnlock()

		metric, ok := storage.data[cityName]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
		}

		raw, err := json.Marshal(metric)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(raw)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	jobs, err := initJobs(s, storage)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("Starting server on port", httpPort)
		err := http.ListenAndServe(httpPort, r)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		fmt.Printf("Starting job: %v\n", jobs[0].ID())
		s.Start()
	}()

	wg.Wait()
}

func initJobs(scheduler gocron.Scheduler, storage *Storage) ([]gocron.Job, error) {
	j, err := scheduler.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				geoResponse, err := geocoding.GetCoords(city) // получаем координаты
				if err != nil {
					log.Println(err)
					return
				}
				// получаем температуру по ранее полученным координатам
				openMeteoResponse, err := openmeteo.GetTemp(geoResponse.Latitude, geoResponse.Longitude)
				if err != nil {
					log.Println(err)
					return
				}

				storage.mu.Lock()
				defer storage.mu.Unlock()

				timestamp, err := time.Parse("2006-01-02T15:04", openMeteoResponse.Current.Time)
				if err != nil {
					log.Println(err)
					return
				}
				storage.data[city] = append(storage.data[city], Metrics{
					Timestamp:   timestamp,
					Temperature: openMeteoResponse.Current.Temperature2m,
				})

				fmt.Printf("updated data for city: %s\n", city)
			},
		),
	)
	if err != nil {
		return nil, err
	}
	return []gocron.Job{j}, nil

}
