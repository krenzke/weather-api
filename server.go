package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)


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
  log.Println("getForecast", r.URL.Path)

	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","GET")

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
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Failed to load .env file",err)
    return
  }
  _, ok := os.LookupEnv("PIRATE_API_KEY")
  if (!ok) {
    log.Fatal("Missing PIRATE_API_KEY from .env")
    return
  }
  _, ok = os.LookupEnv("GOOGLE_API_KEY")
  if (!ok) {
    log.Fatal("Missing GOOGLE_API_KEY from .env")
    return
  }

  port, ok := os.LookupEnv("PORT")
  if (!ok) {
    port = "8000"
  }

  log.Printf("Starting Server, listening on port %s\n", port)

  http.HandleFunc("/weather/{zipcode}", getForecast)
  http.ListenAndServe(":" + port, nil)
}