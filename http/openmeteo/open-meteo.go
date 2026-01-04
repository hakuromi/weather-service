package openmeteo // отвечает за получение температуры по координатам

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

/* Стороннее api возвращает ответ в виде:
{
  "latitude": 55.75,
  "longitude": 37.625,
  "generationtime_ms": 0.0175237655639648,
  "utc_offset_seconds": 0,
  "timezone": "GMT",
  "timezone_abbreviation": "GMT",
  "elevation": 141,
  "current_units": {
    "time": "iso8601",
    "interval": "seconds",
    "temperature_3m": "undefined"
  },
  "current": {
    "time": "2025-12-14T12:30",
    "interval": 900,
    "temperature_3m": null
  }
}
*/

type OpenMeteoResp struct {
	Current struct {
		Time          string  `json:"time"`
		Temperature2m float64 `json:"temperature"`
	}
}

func GetTemp(lat, lon float64) (OpenMeteoResp, error) {
	openMeteoHttpClient := http.Client{ // новый клиент для запросов к api
		Timeout: time.Duration(10 * time.Second),
	}

	openMeteoResponse, err := openMeteoHttpClient.Get( // get запрос
		fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%v&longitude=%v&current=temperature",
			lat, // передаем широту и долготу
			lon,
		),
	)

	if err != nil { // ошибка при get запросе
		return OpenMeteoResp{}, err
	}

	defer openMeteoResponse.Body.Close()

	if openMeteoResponse.StatusCode != http.StatusOK { // проверяем статус код ответа open-meteo
		return OpenMeteoResp{}, fmt.Errorf("status code %d", openMeteoResponse.StatusCode)
	}

	var result OpenMeteoResp
	err = json.NewDecoder(openMeteoResponse.Body).Decode(&result)
	if err != nil {
		return OpenMeteoResp{}, err
	}
	return result, nil
}
