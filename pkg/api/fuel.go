package api

import (
	"sync"
	"time"
)

type FuelType string

const (
	Diesel   FuelType = "Diesel"
	Gasoline FuelType = "Gasoline"
	Gas      FuelType = "Gas"
)

type Fuel struct {
	Type  FuelType `json:"type"`
	Price float64  `json:"price"`
}

var Stations []Station

type FuelPriceRecord struct {
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

var FuelPriceHistory = struct {
	sync.RWMutex
	m map[FuelType][]FuelPriceRecord
}{m: make(map[FuelType][]FuelPriceRecord)}

func UpdateFuelPrice(fuelType FuelType, price float64) {
	record := FuelPriceRecord{
		Price:     price,
		Timestamp: time.Now(),
	}
	FuelPriceHistory.Lock()
	FuelPriceHistory.m[fuelType] = append(FuelPriceHistory.m[fuelType], record)
	FuelPriceHistory.Unlock()
}

func GetFuelPriceHistory(fuelType FuelType) []FuelPriceRecord {
	FuelPriceHistory.RLock()
	history := FuelPriceHistory.m[fuelType]
	FuelPriceHistory.RUnlock()
	return history
}
