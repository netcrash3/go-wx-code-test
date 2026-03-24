package models

type ForecastResponse struct {
	ShortForecast            string `json:"shortForecast"`
	Temperature              int    `json:"temperature"`
	TemperatureUnit          string `json:"temperatureUnit"`
	TemperatureCharacterization string `json:"temperatureCharacterization"`
}

type HealthResponse struct {
	Status string `json:"status"`
}
