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
					SpaceID:     fmt.Sprintf("%d", space.ID),
					Reason:      "First available space strategy",
				}, nil
			}
		}
	}

	return nil, errors.New("no available parking spaces")
}

// Enhanced parking decision with strategy support
func (pa *ParkingAttendant) MakeParkingDecisionWithStrategy(lots []*ParkingLot, car *Car, strategy ParkingStrategy) (*ParkingDecision, error) {
	if !pa.IsActive {
		return nil, errors.New("attendant is not active")
	}

	if strategy == nil {
		// Fall back to original decision logic
		return pa.MakeParkingDecision(lots, car)
	}

	// Use the provided strategy to find the best lot
	selectedLot, err := strategy.FindParkingLot(lots, car)
	if err != nil {
		return nil, err
	}

	// Find an available space in the selected lot
	space := selectedLot.FindAvailableSpace()
	if space == nil {
		return nil, errors.New("no available space in selected lot")
	}

	return &ParkingDecision{
		AttendantID: pa.ID,
		LotID:       selectedLot.ID,
		SpaceID:     fmt.Sprintf("%d", space.ID),
		Reason:      fmt.Sprintf("Strategy: %s", strategy.GetStrategyName()),
	}, nil
}

func (pa *ParkingAttendant) SetStrategy(strategy ParkingStrategy) {
	// This could be expanded to store strategy in attendant
	// For now, strategy is passed per decision
}
