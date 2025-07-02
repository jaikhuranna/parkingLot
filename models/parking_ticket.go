package models

import (
	"fmt"
	"time"
)

type ParkingTicket struct {
	ID           string
	LicensePlate string
	LotID        string
	SpaceID      string
	ParkedAt     time.Time
	UnparkedAt   time.Time
	IsActive     bool
	AttendantID  string
}

func NewParkingTicket(licensePlate, lotID, spaceID string) *ParkingTicket {
	return &ParkingTicket{
		ID:           generateTicketID(licensePlate, lotID, spaceID),
		LicensePlate: licensePlate,
		LotID:        lotID,
		SpaceID:      spaceID,
		ParkedAt:     time.Now(),
		IsActive:     true,
		AttendantID:  "",
	}
}

func NewParkingTicketWithAttendant(licensePlate, lotID, spaceID, attendantID string) *ParkingTicket {
	ticket := NewParkingTicket(licensePlate, lotID, spaceID)
	ticket.AttendantID = attendantID
	return ticket
}

func generateTicketID(licensePlate, lotID, spaceID string) string {
	return fmt.Sprintf("%s_%s_%s_%d", licensePlate, lotID, spaceID, time.Now().Unix())
}

func (pt *ParkingTicket) GetParkingDuration() time.Duration {
	if pt.IsActive {
		return time.Since(pt.ParkedAt)
	}
	return pt.UnparkedAt.Sub(pt.ParkedAt)
}

func (pt *ParkingTicket) CompleteParking() {
	pt.IsActive = false
	pt.UnparkedAt = time.Now()
}

func (pt *ParkingTicket) GetTicketInfo() map[string]interface{} {
	return map[string]interface{}{
		"ID":           pt.ID,
		"LicensePlate": pt.LicensePlate,
		"LotID":        pt.LotID,
		"SpaceID":      pt.SpaceID,
		"ParkedAt":     pt.ParkedAt,
		"UnparkedAt":   pt.UnparkedAt,
		"IsActive":     pt.IsActive,
		"AttendantID":  pt.AttendantID,
		"Duration":     pt.GetParkingDuration(),
	}
}
