package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Raw model for incoming data from WeatherAPI.com
type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Humidity  int     `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

// Clean and formatted JSON data that API will return
type ClientResponse struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature_celsius"`
	Humidity    int     `json:"humidity_percentage"`
	Condition   string  `json:"condition"`
}

// Error struct for API error logs clients.
type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found.")
	}

	// Enforce that API key exists before starting the server
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is not set.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Fallback
	}

	// Define route and start up the server
	http.HandleFunc("/api/weather", handler(apiKey))

	log.Printf("Server starting at http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Wrap the API logic and inject API key dependency securely
func handler(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only accept HTTP GET requests
		if r.Method != http.MethodGet {
			respondWithError(w, http.StatusMethodNotAllowed, "Only GET request is allowed")
			return
		}

		// Extract and clean query parameters
		city := strings.TrimSpace(r.URL.Query().Get("city"))
		if city == "" {
			respondWithError(w, http.StatusBadRequest, "Missing required query parameter: 'city'")
			return
		}

		// Safe URL endpoint
		escapedCity := url.QueryEscape(city)
		apiURL := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, escapedCity)

		// Call 3rd-party API
		client := &http.Client{Timeout: 10 * time.Second}
		response, err := client.Get(apiURL)
		if err != nil {
			log.Printf("Network error fetching weather for %s: %v", city, err)
			respondWithError(w, http.StatusInternalServerError, "Failed to connect to upstream weather provider")
			return
		}
		defer response.Body.Close()

		// Handle non-200 responses safely
		if response.StatusCode != http.StatusOK {
			if response.StatusCode == http.StatusBadRequest {
				respondWithError(w, http.StatusNotFound, fmt.Sprintf("City '%s' could not be found", city))
				return
			}
			log.Printf("Upstream provider returned an error status: %d", response.StatusCode)
			respondWithError(w, http.StatusBadGateway, "Upstream weather provider returned an error status")
			return
		}

		// Decode the JSON stream
		var rawWeather Weather
		if err := json.NewDecoder(response.Body).Decode(&rawWeather); err != nil {
			log.Printf("JSON parsing error: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to process weather data payload")
			return
		}

		// Map ClientResponse
		clientData := ClientResponse{
			Location:    fmt.Sprintf("%s, %s", rawWeather.Location.Name, rawWeather.Location.Country),
			Temperature: rawWeather.Current.TempC,
			Humidity:    rawWeather.Current.Humidity,
			Condition:   rawWeather.Current.Condition.Text,
		}

		// Send final JSON package to user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(clientData); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	}
}

// Ensure that all error responses share the same clean JSON
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
