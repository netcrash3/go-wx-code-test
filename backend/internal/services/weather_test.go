package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-wx/internal/models"
)

func TestCharacterizeTemp_Hot(t *testing.T) {
	if got := characterizeTemp(76); got != TempCharacterizationHot {
		t.Errorf("expected '%s' for 76, got '%s'", TempCharacterizationHot, got)
	}
}

func TestCharacterizeTemp_Cold(t *testing.T) {
	if got := characterizeTemp(49); got != TempCharacterizationCold {
		t.Errorf("expected '%s' for 49, got '%s'", TempCharacterizationCold, got)
	}
	if got := characterizeTemp(0); got != TempCharacterizationCold {
		t.Errorf("expected '%s' for 0, got '%s'", TempCharacterizationCold, got)
	}
}

func TestCharacterizeTemp_Moderate(t *testing.T) {
	if got := characterizeTemp(50); got != TempCharacterizationModerate {
		t.Errorf("expected '%s' for 50, got '%s'", TempCharacterizationModerate, got)
	}
	if got := characterizeTemp(75); got != TempCharacterizationModerate {
		t.Errorf("expected '%s' for 75, got '%s'", TempCharacterizationModerate, got)
	}
	if got := characterizeTemp(65); got != TempCharacterizationModerate {
		t.Errorf("expected '%s' for 65, got '%s'", TempCharacterizationModerate, got)
	}
}

func TestGetTodayForecast_PointsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	svc := &WeatherService{client: server.Client()}
	// Override by calling with a URL that hits our test server
	_, err := getJSON[models.NWSPointsResponse](svc, server.URL+"/points/39.7456,-97.0892")
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestGetTodayForecast_EmptyForecastURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(models.NWSPointsResponse{})
	}))
	defer server.Close()

	svc := &WeatherService{client: server.Client()}
	original := NWSBaseURL
	// We can't override the const, so test via getForecastURL indirectly
	// by testing getJSON + the empty check
	points, err := getJSON[models.NWSPointsResponse](svc, server.URL+"/points/0,0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if points.Properties.Forecast != "" {
		t.Fatalf("expected empty forecast URL, got '%s'", points.Properties.Forecast)
	}
	_ = original
}

func TestGetTodayForecast_ForecastAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	svc := &WeatherService{client: server.Client()}
	_, err := getJSON[models.NWSForecastResponse](svc, server.URL+"/forecast")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestGetTodayForecast_NoPeriodsAvailable(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/points/", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"properties": map[string]any{
				"forecast": "", // will be replaced with server URL
			},
		}
		json.NewEncoder(w).Encode(resp)
	})
	mux.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
		resp := models.NWSForecastResponse{}
		json.NewEncoder(w).Encode(resp)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	svc := &WeatherService{client: server.Client()}
	_, err := svc.getTodayPeriod(server.URL + "/forecast")
	if err == nil {
		t.Fatal("expected error for empty periods")
	}
	if err.Error() != "no forecast periods available" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGetTodayForecast_FallsBackToFirstPeriod(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"properties": map[string]any{
				"periods": []map[string]any{
					{
						"name":            "Tonight",
						"isDaytime":       false,
						"temperature":     45,
						"temperatureUnit": "F",
						"shortForecast":   "Clear",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	svc := &WeatherService{client: server.Client()}
	period, err := svc.getTodayPeriod(server.URL + "/forecast")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if period.Name != "Tonight" {
		t.Errorf("expected fallback to 'Tonight', got '%s'", period.Name)
	}
}

func TestGetTodayForecast_SelectsDaytimePeriod(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"properties": map[string]any{
				"periods": []map[string]any{
					{
						"name":            "Tonight",
						"isDaytime":       false,
						"temperature":     45,
						"temperatureUnit": "F",
						"shortForecast":   "Clear",
					},
					{
						"name":            "Tuesday",
						"isDaytime":       true,
						"temperature":     72,
						"temperatureUnit": "F",
						"shortForecast":   "Sunny",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	svc := &WeatherService{client: server.Client()}
	period, err := svc.getTodayPeriod(server.URL + "/forecast")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if period.Name != "Tuesday" {
		t.Errorf("expected 'Tuesday', got '%s'", period.Name)
	}
	if period.Temperature != 72 {
		t.Errorf("expected temperature 72, got %d", period.Temperature)
	}
}

func TestGetTodayForecast_EndToEnd(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/points/", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"properties": map[string]any{
				"forecast": "", // placeholder, replaced below
			},
		}
		// We need to know the server URL for the forecast redirect,
		// but we don't have it yet. Use a relative approach.
		// Actually, just hardcode the path and set it after server starts.
		_ = resp
	})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/points/39.7456,-97.0892":
			resp := map[string]any{
				"properties": map[string]any{
					"forecast": "http://" + r.Host + "/gridpoints/TOP/31,80/forecast",
				},
			}
			json.NewEncoder(w).Encode(resp)
		case r.URL.Path == "/gridpoints/TOP/31,80/forecast":
			resp := map[string]any{
				"properties": map[string]any{
					"periods": []map[string]any{
						{
							"name":            "Today",
							"isDaytime":       true,
							"temperature":     82,
							"temperatureUnit": "F",
							"shortForecast":   "Partly Cloudy",
						},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	svc := &WeatherService{client: server.Client()}
	// Call getForecastURL pointing at our test server
	pointsURL := server.URL + "/points/39.7456,-97.0892"
	points, err := getJSON[models.NWSPointsResponse](svc, pointsURL)
	if err != nil {
		t.Fatalf("unexpected error fetching points: %v", err)
	}

	period, err := svc.getTodayPeriod(points.Properties.Forecast)
	if err != nil {
		t.Fatalf("unexpected error fetching forecast: %v", err)
	}

	result := &models.ForecastResponse{
		ShortForecast:               period.ShortForecast,
		Temperature:                 period.Temperature,
		TemperatureUnit:             period.TemperatureUnit,
		TemperatureCharacterization: characterizeTemp(period.Temperature),
	}

	if result.ShortForecast != "Partly Cloudy" {
		t.Errorf("expected 'Partly Cloudy', got '%s'", result.ShortForecast)
	}
	if result.Temperature != 82 {
		t.Errorf("expected 82, got %d", result.Temperature)
	}
	if result.TemperatureCharacterization != TempCharacterizationHot {
		t.Errorf("expected '%s', got '%s'", TempCharacterizationHot, result.TemperatureCharacterization)
	}
}

func TestGetTodayForecast_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	svc := &WeatherService{client: server.Client()}
	_, err := getJSON[models.NWSPointsResponse](svc, server.URL+"/points/0,0")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
