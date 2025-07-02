package models

import "time"

// CarLocation represents the location information of a parked car
type CarLocation struct {
	Car         *Car
	LotID       string
	SpaceID     string
	Row         string
	Position    int
	ParkedAt    time.Time
	AttendantID string
}

func NewCarLocation(car *Car, lotID, spaceID, row string, position int, attendantID string) *CarLocation {
	return &CarLocation{
		Car:         car,
		LotID:       lotID,
		SpaceID:     spaceID,
		Row:         row,
		Position:    position,
		ParkedAt:    time.Now(),
		AttendantID: attendantID,
	}
}

func (cl *CarLocation) GetLocationInfo() map[string]interface{} {
	return map[string]interface{}{
		"LicensePlate": cl.Car.LicensePlate,
		"LotID":        cl.LotID,
		"SpaceID":      cl.SpaceID,
		"Row":          cl.Row,
		"Position":     cl.Position,
		"ParkedAt":     cl.ParkedAt,
		"AttendantID":  cl.AttendantID,
	}
}
