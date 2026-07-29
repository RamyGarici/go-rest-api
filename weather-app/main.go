package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

type cityData struct{
	Name string `json:"name"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}


func main() {
	client := &http.Client{Timeout: time.Second * 10}
	var city string
	fmt.Println("City:")
	_,err := fmt.Scan(&city)
	if err != nil {
        fmt.Println("Error reading input:", err)
        return
    }
	var responseCity cityData
	req, err := http.NewRequest(http.MethodGet, "https://api.openweathermap.org/data/2.5/weather?lat=48.85&lon=2.35&appid=b73142ec6a46fe59f577ee9fcd7b6182&units=metric", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&responseCity); err != nil {
		fmt.Println("error decoding response body")
		return
	}
	lat := responseCity.Lat
	lon:=responseCity.Lon
	

url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%v&lon=%v&appid=b73142ec6a46fe59f577ee9fcd7b6182&units=metric",lat, lon)
	
	


	var response weatherResponse
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	decoder = json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		fmt.Println("error decoding response body")
		return
	}
fmt.Printf("Temp: %v\n",response.Main.Temp)
fmt.Printf("Feels like: %v\n",response.Main.FeelsLike)
fmt.Println(response.Weather[0].Description)

}

