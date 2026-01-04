package geocoding // отвечает за получение координат

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

/* стороннее api возвращает ответ в формате слайса структур:
{
  "results": [
    {
      "id": 524901,
      "name": "Москва",
      "latitude": 55.75222,
      "longitude": 37.61556,
      "elevation": 144,
      "feature_code": "PPLC",
      "country_code": "RU",
      "admin1_id": 524894,
      "timezone": "Europe/Moscow",
      "population": 10381222,
      "country_id": 2017370,
      "country": "Россия",
      "admin1": "Москва"
    }
  ],
  "generationtime_ms": 0.4042387
}
*/

type Response struct {
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type geoResponse struct {
	Results []Response `json:"results"`
}

func GetCoords(city string) (Response, error) {
	newHttpClient := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	response, err := newHttpClient.Get( // get request for take coords of a city
		fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=ru&format=json",
			city),
	)
	if err != nil { // check success get request
		return Response{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK { // check response status code
		return Response{}, fmt.Errorf("status code %d", response.StatusCode)
	}

	var georesp geoResponse // создаем экземпляр структуры georesponse

	err = json.NewDecoder(response.Body).Decode(&georesp)
	if err != nil {
		return Response{}, err
	}

	return georesp.Results[0], nil // берем первый элемент слайса структур
}
