package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"go-wx/internal/models"
)

type mockForecaster struct {
	result *models.ForecastResponse
	err    error
}

func (m *mockForecaster) GetTodayForecast(lat, lon float64) (*models.ForecastResponse, error) {
	return m.result, m.err
}

func setupRouter(svc Forecaster) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewForecastHandler(svc)
	r.GET("/api/forecast", h.GetForecast)
	return r
}

func TestGetForecast_MissingLat(t *testing.T) {
	r := setupRouter(&mockForecaster{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/forecast?lon=-97.0892", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
	expectErrorContains(t, w, "lat and lon query parameters are required")
}

func TestGetForecast_MissingLon(t *testing.T) {
	r := setupRouter(&mockForecaster{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/forecast?lat=39.7456", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
	expectErrorContains(t, w, "lat and lon query parameters are required")
}

func TestGetForecast_MissingBoth(t *testing.T) {
	r := setupRouter(&mockForecaster{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/forecast", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
	expectErrorContains(t, w, "lat and lon query parameters are required")
}

func TestGetForecast_InvalidLat(t *testing.T) {
	r := setupRouter(&mockForecaster{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/forecast?lat=abc&lon=-97.0892", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
	expectErrorContains(t, w, "lat must be a valid number")
}

func TestGetForecast_InvalidLon(t *testing.T) {
	r := setupRouter(&mockForecaster{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/forecast?lat=39.7456&lon=xyz", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
	expectErrorContains(t, w, "lon must be a valid number")
}

func TestGetForecast_ServiceError(t *testing.T) {
	mock := &mockForecaster{err: fmt.Errorf("NWS API returned status 500 for https://api.weather.gov/points/0.0000,0.0000")}
	r := setupRouter(mock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/forecast?lat=0&lon=0", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Fatalf("expected status 502, got %d", w.Code)
	}
	expectErrorContains(t, w, "NWS API returned status 500")
}

func TestGetForecast_Success(t *testing.T) {
	mock := &mockForecaster{
		result: &models.ForecastResponse{
			ShortForecast:               "Partly Cloudy",
			Temperature:                 68,
			TemperatureUnit:             "F",
			TemperatureCharacterization: "moderate",
		},
	}
	r := setupRouter(mock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/forecast?lat=39.7456&lon=-97.0892", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp models.ForecastResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.ShortForecast != "Partly Cloudy" {
		t.Errorf("expected ShortForecast 'Partly Cloudy', got '%s'", resp.ShortForecast)
	}
	if resp.Temperature != 68 {
		t.Errorf("expected Temperature 68, got %d", resp.Temperature)
	}
	if resp.TemperatureUnit != "F" {
		t.Errorf("expected TemperatureUnit 'F', got '%s'", resp.TemperatureUnit)
	}
	if resp.TemperatureCharacterization != "moderate" {
		t.Errorf("expected TemperatureCharacterization 'moderate', got '%s'", resp.TemperatureCharacterization)
	}
}

func expectErrorContains(t *testing.T, w *httptest.ResponseRecorder, substr string) {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}
	errMsg, ok := body["error"]
	if !ok {
		t.Fatal("expected 'error' key in response body")
	}
	if len(errMsg) == 0 {
		t.Fatal("expected non-empty error message")
	}
	for i := range errMsg {
		if i+len(substr) <= len(errMsg) && errMsg[i:i+len(substr)] == substr {
			return
		}
	}
	t.Errorf("expected error to contain '%s', got '%s'", substr, errMsg)
}
