package generator

import (
	"log"
	"math/rand"
	"time"

	"github.com/fsipic/go-rivo/pkg/api"
)

func StartPriceFluctuation() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fluctuatePrices()
	}
}

func fluctuatePrices() {
	println("Fluctuating prices...")
	stations := api.GetStations()

	for _, station := range stations {
		for i, fuel := range station.Fuels {
			factor := 0.3 + rand.Float64()*(2.0-0.3)
			newPrice := fuel.Price * factor
			log.Printf("Updated price for %s at station %s from %.2f to %.2f", fuel.Type, station.Name, fuel.Price, newPrice)
			station.Fuels[i].Price = newPrice
			api.UpdateFuelPrice(fuel.Type, newPrice)
		}
	}
}
