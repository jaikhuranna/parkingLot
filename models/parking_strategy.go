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

// HandicapPriorityStrategy prioritizes nearest available spaces for handicap drivers
type HandicapPriorityStrategy struct{}

func NewHandicapPriorityStrategy() *HandicapPriorityStrategy {
    return &HandicapPriorityStrategy{}
}

func (hps *HandicapPriorityStrategy) GetStrategyName() string {
    return "Handicap Priority Strategy"
}

func (hps *HandicapPriorityStrategy) FindParkingLot(lots []*ParkingLot, car *Car) (*ParkingLot, error) {
    if len(lots) == 0 {
        return nil, errors.New("no parking lots available")
    }
    
    if !car.IsHandicap {
        // For non-handicap cars, use even distribution as fallback
        strategy := NewEvenDistributionStrategy()
        return strategy.FindParkingLot(lots, car)
    }
    
    // For handicap drivers, prioritize lots with nearest available spaces
    // In this implementation, we'll prioritize the first available lot
    // In a real system, this would consider physical distance/proximity
    for _, lot := range lots {
        if !lot.IsFull() {
            return lot, nil
        }
    }
    
    return nil, errors.New("no available parking spaces for handicap vehicle")
}

// LargeVehicleStrategy directs large cars to lots with most free space
type LargeVehicleStrategy struct{}

func NewLargeVehicleStrategy() *LargeVehicleStrategy {
    return &LargeVehicleStrategy{}
}

func (lvs *LargeVehicleStrategy) GetStrategyName() string {
    return "Large Vehicle Strategy"
}

func (lvs *LargeVehicleStrategy) FindParkingLot(lots []*ParkingLot, car *Car) (*ParkingLot, error) {
    if len(lots) == 0 {
        return nil, errors.New("no parking lots available")
    }
    
    var bestLot *ParkingLot
    maxAvailable := -1
    
    if car.Size == LargeVehicle {
        // For large vehicles, find lot with most available spaces
        for _, lot := range lots {
            if !lot.IsFull() {
                available := lot.GetAvailableSpaces()
                if available > maxAvailable {
                    maxAvailable = available
                    bestLot = lot
                }
            }
        }
    } else {
        // For small/medium vehicles, use even distribution
        strategy := NewEvenDistributionStrategy()
        return strategy.FindParkingLot(lots, car)
    }
    
    if bestLot == nil {
        return nil, errors.New("no available parking spaces for large vehicle")
    }
    
    return bestLot, nil
}

// SmartParkingStrategy combines handicap priority and large vehicle strategies
type SmartParkingStrategy struct{}

func NewSmartParkingStrategy() *SmartParkingStrategy {
    return &SmartParkingStrategy{}
}

func (sps *SmartParkingStrategy) GetStrategyName() string {
    return "Smart Parking Strategy (Handicap + Large Vehicle)"
}

func (sps *SmartParkingStrategy) FindParkingLot(lots []*ParkingLot, car *Car) (*ParkingLot, error) {
    if len(lots) == 0 {
        return nil, errors.New("no parking lots available")
    }
    
    // Prioritize handicap drivers first
    if car.IsHandicap {
        strategy := NewHandicapPriorityStrategy()
        return strategy.FindParkingLot(lots, car)
    }
    
    // Then handle large vehicles
    if car.Size == LargeVehicle {
        strategy := NewLargeVehicleStrategy()
        return strategy.FindParkingLot(lots, car)
    }
    
    // Default to even distribution for regular cars
    strategy := NewEvenDistributionStrategy()
    return strategy.FindParkingLot(lots, car)
}
