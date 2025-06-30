package services

import (
    "parking-lot-system/models"
    "errors"
)

type ParkingService struct {
    lots []*models.ParkingLot
}

func NewParkingService() *ParkingService {
    return &ParkingService{
        lots: make([]*models.ParkingLot, 0),
    }
}

func (ps *ParkingService) AddLot(lot *models.ParkingLot) {
    ps.lots = append(ps.lots, lot)
}

func (ps *ParkingService) ParkCar(car *models.Car) error {
    if car == nil {
        return errors.New("car cannot be nil")
    }
    
    for _, lot := range ps.lots {
        if err := lot.ParkCar(car); err == nil {
            return nil
        }
    }
    
    return errors.New("no available parking space")
}

func (ps *ParkingService) UnparkCar(licensePlate string) (*models.Car, error) {
    if licensePlate == "" {
        return nil, errors.New("license plate cannot be empty")
    }
    
    for _, lot := range ps.lots {
        if car, err := lot.UnparkCar(licensePlate); err == nil {
            return car, nil
        }
    }
    
    return nil, errors.New("car not found in any parking lot")
}

func (ps *ParkingService) FindCar(licensePlate string) (*models.ParkingSpace, error) {
    if licensePlate == "" {
        return nil, errors.New("license plate cannot be empty")
    }
    
    for _, lot := range ps.lots {
        if space := lot.FindCar(licensePlate); space != nil {
            return space, nil
        }
    }
    
    return nil, errors.New("car not found")
}

