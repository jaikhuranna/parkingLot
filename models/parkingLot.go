package models

import "errors"

type ParkingLot struct {
    ID       string
    Capacity int
    Spaces   []*ParkingSpace
}

func NewParkingLot(id string, capacity int) *ParkingLot {
    lot := &ParkingLot{
        ID:       id,
        Capacity: capacity,
        Spaces:   make([]*ParkingSpace, capacity),
    }
    
    for i := 0; i < capacity; i++ {
        lot.Spaces[i] = NewParkingSpace(i + 1)
    }
    
    return lot
}

func (pl *ParkingLot) ParkCar(car *Car) error {
    for _, space := range pl.Spaces {
        if space.Park(car) {
            return nil
        }
    }
    return errors.New("parking lot is full")
}

func (pl *ParkingLot) UnparkCar(licensePlate string) (*Car, error) {
    for _, space := range pl.Spaces {
        if space.IsOccupied && space.ParkedCar != nil && 
           space.ParkedCar.LicensePlate == licensePlate {
            return space.Unpark(), nil
        }
    }
    return nil, errors.New("car not found in parking lot")
}

func (pl *ParkingLot) FindCar(licensePlate string) *ParkingSpace {
    for _, space := range pl.Spaces {
        if space.IsOccupied && space.ParkedCar != nil && 
           space.ParkedCar.LicensePlate == licensePlate {
            return space
        }
    }
    return nil
}

