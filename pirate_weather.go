package main

import (
	"fmt"
	"io"
	"net/http"
)

const WEATHER_API_KEY = "TODO"

func getWeatherData(location Location) ([]byte, error){
	url := fmt.Sprintf("https://api.pirateweather.net/forecast/%s/%f,%f", WEATHER_API_KEY, location.Lat, location.Lng)
	println("Getting weather from", url)

  res, err := http.Get(url)
  if (err != nil) {
    println(err)
    return nil, err
  }
  defer res.Body.Close()

  body, err := io.ReadAll(res.Body)
  if (err != nil) {
    println(err)
    return nil, err
  }
	return body, nil
}