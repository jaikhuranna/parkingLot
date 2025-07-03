package models

import "errors"

// ParkingStrategy interface for different parking allocation strategies
type ParkingStrategy interface {
	FindParkingLot(lots []*ParkingLot, car *Car) (*ParkingLot, error)
	GetStrategyName() string
}

// EvenDistributionStrategy implements even distribution across parking lots
type EvenDistributionStrategy struct{}

func NewEvenDistributionStrategy() *EvenDistributionStrategy {
	return &EvenDistributionStrategy{}
}

func (eds *EvenDistributionStrategy) GetStrategyName() string {
	return "Even Distribution Strategy"
}

func (eds *EvenDistributionStrategy) FindParkingLot(lots []*ParkingLot, car *Car) (*ParkingLot, error) {
	if len(lots) == 0 {
		return nil, errors.New("no parking lots available")
	}

	var bestLot *ParkingLot
	maxAvailable := -1

	// Find the lot with the most available spaces for even distribution
	for _, lot := range lots {
		if !lot.IsFull() {
			available := lot.GetAvailableSpaces()
			if available > maxAvailable {
				maxAvailable = available
				bestLot = lot
			}
		}
	}

	if bestLot == nil {
		return nil, errors.New("no available parking spaces in any lot")
	}

	return bestLot, nil
}

// LotUtilization represents the utilization statistics of a parking lot
type LotUtilization struct {
	LotID           string
	TotalSpaces     int
	OccupiedSpaces  int
	AvailableSpaces int
	UtilizationRate float64
}

func CalculateLotUtilization(lot *ParkingLot) *LotUtilization {
	occupied := lot.GetOccupiedSpaces()
	available := lot.GetAvailableSpaces()
	total := lot.Capacity

	var utilizationRate float64
	if total > 0 {
		utilizationRate = float64(occupied) / float64(total) * 100
	}

	return &LotUtilization{
		LotID:           lot.ID,
		TotalSpaces:     total,
		OccupiedSpaces:  occupied,
		AvailableSpaces: available,
		UtilizationRate: utilizationRate,
	}
}
