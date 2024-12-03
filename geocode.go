package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GoogleGeocodeResponse struct {
	Results []struct {
		// AddressComponents []struct {
		// 	LongName  string   `json:"long_name"`
		// 	ShortName string   `json:"short_name"`
		// 	Types     []string `json:"types"`
		// } `json:"address_components"`
		// FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			// Bounds struct {
			// 	Northeast struct {
			// 		Lat float64 `json:"lat"`
			// 		Lng float64 `json:"lng"`
			// 	} `json:"northeast"`
			// 	Southwest struct {
			// 		Lat float64 `json:"lat"`
			// 		Lng float64 `json:"lng"`
			// 	} `json:"southwest"`
			// } `json:"bounds"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			// Viewport     struct {
			// 	Northeast struct {
			// 		Lat float64 `json:"lat"`
			// 		Lng float64 `json:"lng"`
			// 	} `json:"northeast"`
			// 	Southwest struct {
			// 		Lat float64 `json:"lat"`
			// 		Lng float64 `json:"lng"`
			// 	} `json:"southwest"`
			// } `json:"viewport"`
		} `json:"geometry"`
		// PlaceID            string   `json:"place_id"`
		// PostcodeLocalities []string `json:"postcode_localities"`
		// Types              []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

type Location struct{
	Lat float64
	Lng float64
}

var ZipcodeCache = make(map[string]Location)

func geocodeZipcode(zipcode string) (Location,error) {
	// first, check the cache
	location, ok := ZipcodeCache[zipcode]
	if (ok) {
		println("Cache hit for ", zipcode)
		return location, nil
	}

	// no cache hit, go ask google
	println("Cache miss for ", zipcode)

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", zipcode,  os.Getenv("GOOGLE_API_KEY"))
  res, err := http.Get(url)
  if (err != nil) {
    println(err)
    return Location{}, err
  }
  defer res.Body.Close()

  body, err := io.ReadAll(res.Body)
  if (err != nil) {
    println(err)
    return Location{}, err
  }

  var parsedResponse GoogleGeocodeResponse
  json.Unmarshal(body, &parsedResponse)
	location = Location{
		Lat: parsedResponse.Results[0].Geometry.Location.Lat,
		Lng: parsedResponse.Results[0].Geometry.Location.Lng,
	}

	ZipcodeCache[zipcode] = location

	return location, nil
}