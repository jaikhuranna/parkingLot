package models

type ParkingSpace struct {
    ID         int
    IsOccupied bool
    ParkedCar  *Car
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
    return true
}

