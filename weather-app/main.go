package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type weatherResponse struct {
	Main    mainData `json:"main"`
	Weather []weather `json:"weather"`
}

type mainData struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
}

type weather struct {
	Description string `json:"description"`
}

type cityData struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

func getCoordinates(client *http.Client, city string, apiKey string) (float64, float64, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s",
		city,
		apiKey,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("geocoding API returned: %s", resp.Status)
	}

	var cities []cityData

	if err := json.NewDecoder(resp.Body).Decode(&cities); err != nil {
		return 0, 0, err
	}

	if len(cities) == 0 {
		return 0, 0, fmt.Errorf("city not found")
	}

	return cities[0].Lat, cities[0].Lon, nil
}

func getWeather(client *http.Client, lat float64, lon float64, apiKey string) (weatherResponse, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?lat=%v&lon=%v&appid=%s&units=metric",
		lat,
		lon,
		apiKey,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return weatherResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return weatherResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weatherResponse{}, fmt.Errorf("weather API returned: %s", resp.Status)
	}

	var weather weatherResponse

	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return weatherResponse{}, err
	}

	return weather, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env:", err)
		return
	}

	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		fmt.Println("API_KEY is not set")
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var city string

	fmt.Print("City: ")

	_, err = fmt.Scan(&city)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	lat, lon, err := getCoordinates(client, city, apiKey)
	if err != nil {
		fmt.Println("Error getting coordinates:", err)
		return
	}

	weather, err := getWeather(client, lat, lon, apiKey)
	if err != nil {
		fmt.Println("Error getting weather:", err)
		return
	}

	fmt.Printf("Temperature: %.1f°C\n", weather.Main.Temp)
	fmt.Printf("Feels like: %.1f°C\n", weather.Main.FeelsLike)

	if len(weather.Weather) > 0 {
		fmt.Printf("Condition: %s\n", weather.Weather[0].Description)
	}
}