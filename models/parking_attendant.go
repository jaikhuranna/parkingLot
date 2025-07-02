package models

import (
	"errors"
	"fmt"
)

// ParkingAttendant represents a parking lot attendant
type ParkingAttendant struct {
    ID       string
    Name     string
    LotID    string
    IsActive bool
}

func NewParkingAttendant(id, name, lotID string) *ParkingAttendant {
    return &ParkingAttendant{
        ID:       id,
        Name:     name,
        LotID:    lotID,
        IsActive: true,
    }
}

func (pa *ParkingAttendant) SetActive(status bool) {
    pa.IsActive = status
}

func (pa *ParkingAttendant) GetInfo() map[string]interface{} {
    return map[string]interface{}{
        "ID":       pa.ID,
        "Name":     pa.Name,
        "LotID":    pa.LotID,
        "IsActive": pa.IsActive,
    }
}

// ParkingDecision represents an attendant's parking decision
type ParkingDecision struct {
    AttendantID string
    LotID       string
    SpaceID     string
    Reason      string
}

func (pa *ParkingAttendant) MakeParkingDecision(lots []*ParkingLot, car *Car) (*ParkingDecision, error) {
    if !pa.IsActive {
        return nil, errors.New("attendant is not active")
    }
    
    // Simple decision logic: find first available space
    for _, lot := range lots {
        if !lot.IsFull() {
            space := lot.FindAvailableSpace()
            if space != nil {
                return &ParkingDecision{
                    AttendantID: pa.ID,
                    LotID:       lot.ID,
		    SpaceID:     fmt.Sprintf( "%d", space.ID),
                    Reason:      "First available space strategy",
                }, nil
            }
        }
    }
    
    return nil, errors.New("no available parking spaces")
}
