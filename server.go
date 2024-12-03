package main

import (
	"encoding/json"
	"net/http"
)

type successResponse struct {
	Success bool `json:"success"`
	Data []string `json:"data"`
}

type errorResponse struct {
	Success bool `json:"success"`
	ErrorMessage string `json:"error_message"`
}

func writeErrorResponse(w http.ResponseWriter, err error) {
  w.WriteHeader(401)
  response := errorResponse{
    Success: false,
    ErrorMessage: err.Error(),
  }
  json.NewEncoder(w).Encode(response)
}

func getForecast(w http.ResponseWriter, r *http.Request) {
	println("Go request for", r.PathValue("zipcode"))

	w.Header().Set("Content-Type","application/json")
  // TODO: setup cache
  location, err := geocodeZipcode(r.PathValue("zipcode"))
  if (err != nil) {
    writeErrorResponse(w, err)
    return
  }

  weatherData, err := getWeatherData(location)
  if (err != nil) {
    writeErrorResponse(w, err)
    return
  }

  w.Write(weatherData)
}

func main() {
	println("Starting Server")

  http.HandleFunc("/weather/{zipcode}", getForecast)
  http.ListenAndServe(":8000", nil)
}