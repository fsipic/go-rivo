package api

import (
	"math"
	"sort"
	"sync"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Station struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Address string   `json:"address"`
	Loc     Location `json:"location"`
	Fuels   []Fuel   `json:"fuels"`
}

var stationData = struct {
	sync.RWMutex
	stations []Station
}{stations: make([]Station, 0)}

func AddStation(station Station) {
	stationData.Lock()
	station.ID = len(stationData.stations) + 1
	stationData.stations = append(stationData.stations, station)
	stationData.Unlock()
}

func GetStations() []Station {
	stationData.RLock()
	stationsCopy := make([]Station, len(stationData.stations))
	copy(stationsCopy, stationData.stations)
	stationData.RUnlock()
	return stationsCopy
}

func CreateStation(name, address string, loc Location, fuels []Fuel) Station {
	station := Station{
		Name:    name,
		Address: address,
		Loc:     loc,
		Fuels:   fuels,
	}
	AddStation(station)

	for _, fuel := range fuels {
		UpdateFuelPrice(fuel.Type, fuel.Price)
	}

	return station
}

func FindNearestStations(current Location) []Station {
	type StationDistance struct {
		Station  Station
		Distance float64
	}

	stations := GetStations()
	var distances []StationDistance

	for _, station := range stations {
		dist := distance(current, station.Loc)
		distances = append(distances, StationDistance{Station: station, Distance: dist})
	}

	sort.Slice(distances, func(i, j int) bool {
		return distances[i].Distance < distances[j].Distance
	})

	var nearest []Station
	for i := 0; i < len(distances) && i < 3; i++ {
		nearest = append(nearest, distances[i].Station)
	}
	return nearest
}

func distance(loc1, loc2 Location) float64 {
	var rad = math.Pi / 180
	var haversine = func(theta float64) float64 {
		return 0.5 * (1 - math.Cos(theta))
	}
	dLat := (loc2.Latitude - loc1.Latitude) * rad
	dLon := (loc2.Longitude - loc1.Longitude) * rad
	a := haversine(dLat) + math.Cos(loc1.Latitude*rad)*math.Cos(loc2.Latitude*rad)*haversine(dLon)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return 6371 * c
}
