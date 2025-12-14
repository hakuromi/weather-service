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

const httpPort = ":3000"

func main() {
	r := chi.NewRouter()     // роутер обрабатывет наши адреса
	r.Use(middleware.Logger) // мидлваре - для всех эндпоинтов

	r.Get("/{city}", func(w http.ResponseWriter, r *http.Request) {
		city := chi.URLParam(r, "city")

		fmt.Println("Requested city:", city)

		geoResponse, err := geocoding.GetCoords(city) // получаем координаты
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// получаем температуру по ранее полученным координатам
		openMeteoResponse, err := openmeteo.GetTemp(geoResponse.Latitude, geoResponse.Longitude)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		raw, err := json.Marshal(openMeteoResponse)
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

	jobs, err := initJobs(s)
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
		fmt.Printf("Starting job: %v", jobs[0].ID())
		s.Start()
	}()

	wg.Wait()
}

func initJobs(scheduler gocron.Scheduler) ([]gocron.Job, error) {
	j, err := scheduler.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("Hello!")
			},
		),
	)
	if err != nil {
		return nil, err
	}
	return []gocron.Job{j}, nil

}
