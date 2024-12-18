package main

import (
	"encoding/json"
	"log"
	"net"
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

func healthcheck(w http.ResponseWriter, r *http.Request) {
  log.Println("healthcheck", r.URL.Path)

	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","GET")

  w.Write([]byte("{\"status\":\"OK\"}"))
}

func setupSocketListener(path string) (net.Listener, error) {
	_, err := os.Stat(path)
	if err == nil{
		os.Remove(path)
	}

	listener, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}

	err = os.Chmod(path, 0777)
	if err != nil {
		return nil, err
	}
	return listener, nil
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

	// Set up routing and handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/api/forecast/{zipcode}", getForecast)
	mux.HandleFunc("/api/healthcheck", healthcheck)

	// Listen to either a unix socket (higher priority)
	// or a port
  socketPath, ok := os.LookupEnv("UNIX_SOCKET")
  if (ok) {
		// Using Socket
		log.Println("Listening on socket")
		log.Println(socketPath)

		listener, err := setupSocketListener(socketPath)
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		defer listener.Close()

		server := http.Server{
			Handler: mux,
		}
		server.Serve(listener)
  } else {
		// Using Port
		port, ok := os.LookupEnv("PORT")
		if (!ok) {
			port = "8000"
		}
		log.Printf("Starting Server, listening on port %s\n", port)
		http.ListenAndServe(":" + port, mux)
	}
}