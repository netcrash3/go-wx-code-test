package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-wx/internal/models"
)

const nwsBaseURL = "https://api.weather.gov"

type WeatherService struct {
	client *http.Client
}

func NewWeatherService() *WeatherService {
	return &WeatherService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *WeatherService) GetTodayForecast(lat, lon float64) (*models.ForecastResponse, error) {
	forecastURL, err := s.getForecastURL(lat, lon)
	if err != nil {
		return nil, fmt.Errorf("resolving grid point: %w", err)
	}

	period, err := s.getTodayPeriod(forecastURL)
	if err != nil {
		return nil, fmt.Errorf("fetching forecast: %w", err)
	}

	return &models.ForecastResponse{
		ShortForecast:               period.ShortForecast,
		Temperature:                 period.Temperature,
		TemperatureUnit:             period.TemperatureUnit,
		TemperatureCharacterization: characterizeTemp(period.Temperature),
	}, nil
}

func (s *WeatherService) getForecastURL(lat, lon float64) (string, error) {
	url := fmt.Sprintf("%s/points/%.4f,%.4f", nwsBaseURL, lat, lon)

	points, err := getJSON[models.NWSPointsResponse](s, url)
	if err != nil {
		return "", err
	}

	if points.Properties.Forecast == "" {
		return "", fmt.Errorf("no forecast URL returned for coordinates (%.4f, %.4f)", lat, lon)
	}

	return points.Properties.Forecast, nil
}

func (s *WeatherService) getTodayPeriod(forecastURL string) (*models.NWSForecastPeriod, error) {
	forecast, err := getJSON[models.NWSForecastResponse](s, forecastURL)
	if err != nil {
		return nil, err
	}

	for _, p := range forecast.Properties.Periods {
		if p.IsDaytime {
			return &p, nil
		}
	}

	if len(forecast.Properties.Periods) > 0 {
		return &forecast.Properties.Periods[0], nil
	}

	return nil, fmt.Errorf("no forecast periods available")
}

func getJSON[T any](s *WeatherService, url string) (*T, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "go-wx-app")
	req.Header.Set("Accept", "application/geo+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NWS API returned status %d for %s", resp.StatusCode, url)
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func characterizeTemp(tempF int) string {
	switch {
	case tempF > 75:
		return "hot"
	case tempF < 50:
		return "cold"
	default:
		return "moderate"
	}
}
