package models

import "time"

type ParkingSpace struct {
	ID         int
	IsOccupied bool
	ParkedCar  *Car
	ParkedAt   time.Time
}

func NewParkingSpace(id int) *ParkingSpace {
	return &ParkingSpace{
		ID:         id,
		IsOccupied: false,
		ParkedCar:  nil,
	}
}

func (ps *ParkingSpace) Park(car *Car) bool {
	if ps.IsOccupied {
		return false
	}
	ps.IsOccupied = true
	ps.ParkedCar = car
	ps.ParkedAt = time.Now()
	return true
}

func (ps *ParkingSpace) Unpark() *Car {
	if !ps.IsOccupied {
		return nil
	}

	car := ps.ParkedCar
	ps.IsOccupied = false
	ps.ParkedCar = nil
	ps.ParkedAt = time.Time{}

	return car
}

func (ps *ParkingSpace) GetParkedCar() *Car {
	return ps.ParkedCar
}

func (ps *ParkingSpace) GetLocationDetails() map[string]interface{} {
	return map[string]interface{}{
		"SpaceID":  ps.ID,
		"ParkedAt": ps.ParkedAt,
	}
}

// UC16: Get row assignment based on space ID
func (ps *ParkingSpace) GetRowAssignment() string {
	// Map space IDs to rows (A, B, C, D)
	// Simple mapping: spaces 1-25 = A, 26-50 = B, 51-75 = C, 76-100 = D
	if ps.ID <= 25 {
		return "A"
	} else if ps.ID <= 50 {
		return "B"
	} else if ps.ID <= 75 {
		return "C"
	} else {
		return "D"
	}
}

// UC16: Get detailed location info including row
func (ps *ParkingSpace) GetDetailedLocationInfo() map[string]interface{} {
	return map[string]interface{}{
		"SpaceID":  ps.ID,
		"Row":      ps.GetRowAssignment(),
		"ParkedAt": ps.ParkedAt,
		"IsOccupied": ps.IsOccupied,
	}
}
