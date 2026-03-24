package models

// NWSPointsResponse is the subset of the NWS /points response we need.
type NWSPointsResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

// NWSForecastResponse is the subset of the NWS /gridpoints forecast response we need.
type NWSForecastResponse struct {
	Properties struct {
		Periods []NWSForecastPeriod `json:"periods"`
	} `json:"properties"`
}

type NWSForecastPeriod struct {
	Name            string `json:"name"`
	IsDaytime       bool   `json:"isDaytime"`
	Temperature     int    `json:"temperature"`
	TemperatureUnit string `json:"temperatureUnit"`
	ShortForecast   string `json:"shortForecast"`
}
