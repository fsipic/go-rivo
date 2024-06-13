package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var user User
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdUser := CreateUser(user.Name, user.Email, user.Password)
	err = json.NewEncoder(w).Encode(createdUser)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func createStationHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request at /api/station")

	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var stationInput struct {
		Name    string   `json:"name"`
		Address string   `json:"address"`
		Loc     Location `json:"location"`
		Fuels   []Fuel   `json:"fuels"`
	}

	err := decoder.Decode(&stationInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, f := range stationInput.Fuels {
		if f.Type != Diesel && f.Type != Gasoline && f.Type != Gas {
			http.Error(w, "Invalid fuel type", http.StatusBadRequest)
			return
		}
	}

	createdStation := CreateStation(stationInput.Name, stationInput.Address, stationInput.Loc, stationInput.Fuels)
	err = json.NewEncoder(w).Encode(createdStation)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func getFuelHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	fuelType := r.URL.Query().Get("type")
	ft, err := validateFuelType(fuelType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	history := GetFuelPriceHistory(ft)
	if len(history) == 0 {
		http.Error(w, "No history found for specified fuel type", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(history)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func validateFuelType(fuelType string) (FuelType, error) {
	switch fuelType {
	case string(Diesel):
		return Diesel, nil
	case string(Gasoline):
		return Gasoline, nil
	case string(Gas):
		return Gas, nil
	default:
		return "", fmt.Errorf("invalid fuel type provided")
	}
}

func getNearestStationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	location := Location{Latitude: lat, Longitude: lon}
	nearestStations := FindNearestStations(location)
	err = json.NewEncoder(w).Encode(nearestStations)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/api/user", CorsMiddleware(http.HandlerFunc(createUserHandler)))
	mux.Handle("/api/station", CorsMiddleware(http.HandlerFunc(createStationHandler)))
	mux.Handle("/api/fuel/history", CorsMiddleware(http.HandlerFunc(getFuelHistoryHandler)))
	mux.Handle("/api/station/nearest", CorsMiddleware(http.HandlerFunc(getNearestStationsHandler)))

	return mux
}
