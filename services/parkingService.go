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

