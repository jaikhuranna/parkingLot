package models

import (
    "errors"
    "parking-lot-system/interfaces"
)

type ParkingLot struct {
    ID        string
    Capacity  int
    Spaces    []*ParkingSpace
    observers []interfaces.ParkingLotObserver
    wasFull   bool // Track previous state to avoid duplicate notifications
}

func NewParkingLot(id string, capacity int) *ParkingLot {
    lot := &ParkingLot{
        ID:        id,
        Capacity:  capacity,
        Spaces:    make([]*ParkingSpace, capacity),
        observers: make([]interfaces.ParkingLotObserver, 0),
        wasFull:   false,
    }
    
    for i := 0; i < capacity; i++ {
        lot.Spaces[i] = NewParkingSpace(i + 1)
    }
    
    return lot
}

func (pl *ParkingLot) AddObserver(observer interfaces.ParkingLotObserver) {
    pl.observers = append(pl.observers, observer)
}

func (pl *ParkingLot) RemoveObserver(observer interfaces.ParkingLotObserver) {
    for i, obs := range pl.observers {
        if obs == observer {
            pl.observers = append(pl.observers[:i], pl.observers[i+1:]...)
            break
        }
    }
}

func (pl *ParkingLot) notifyObservers(isFull bool) {
    for _, observer := range pl.observers {
        if isFull {
            observer.OnLotFull(pl.ID)
        } else {
            observer.OnLotAvailable(pl.ID)
        }
    }
}

// Enhanced methods with notifications
func (pl *ParkingLot) ParkCar(car *Car) error {
    for _, space := range pl.Spaces {
        if space.Park(car) {
            // Check if lot became full after parking
            if pl.IsFull() && !pl.wasFull {
                pl.wasFull = true
                pl.notifyObservers(true)
            }
            return nil
        }
    }
    return errors.New("parking lot is full")
}

func (pl *ParkingLot) UnparkCar(licensePlate string) (*Car, error) {
    wasFullBeforeUnpark := pl.IsFull()
    
    for _, space := range pl.Spaces {
        if space.IsOccupied && space.ParkedCar != nil && 
           space.ParkedCar.LicensePlate == licensePlate {
            car := space.Unpark()
            
            // Check if lot became available after unparking
            if wasFullBeforeUnpark && !pl.IsFull() {
                pl.wasFull = false
                pl.notifyObservers(false)
            }
            return car, nil
        }
    }
    return nil, errors.New("car not found in parking lot")
}

// Lot status methods
func (pl *ParkingLot) IsFull() bool {
    return pl.GetAvailableSpaces() == 0
}

func (pl *ParkingLot) GetAvailableSpaces() int {
    count := 0
    for _, space := range pl.Spaces {
        if !space.IsOccupied {
            count++
        }
    }
    return count
}

func (pl *ParkingLot) GetOccupiedSpaces() int {
    return pl.Capacity - pl.GetAvailableSpaces()
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

