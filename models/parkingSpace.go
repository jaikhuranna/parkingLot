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
        "SpaceID":   ps.ID,
        "ParkedAt":  ps.ParkedAt,
    }
}
