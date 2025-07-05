package services

import (
	"fmt"
	"parking-lot-system/models"
	"strings"
	"time"
)

type PoliceService struct {
	parkingService *ParkingService
}

func NewPoliceService(parkingService *ParkingService) *PoliceService {
	return &PoliceService{
		parkingService: parkingService,
	}
}

type VehicleInvestigationInfo struct {
	Car           *models.Car
	LotID         string
	SpaceID       string
	ParkedAt      time.Time
	AttendantID   string
	AttendantName string
}

// UC12: Find all white cars for bomb threat investigation
func (ps *PoliceService) FindWhiteCars() ([]*VehicleInvestigationInfo, error) {
	var whiteCars []*VehicleInvestigationInfo

	for _, lot := range ps.parkingService.lots {
		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil {
				car := space.ParkedCar
				if strings.ToLower(car.Color) == "white" {
					info := &VehicleInvestigationInfo{
						Car:      car,
						LotID:    lot.ID,
						SpaceID:  fmt.Sprintf("%d", space.ID), // FIXED: Convert int to string
						ParkedAt: space.ParkedAt,
					}

					// Try to find attendant info if available
					if ticket, err := ps.parkingService.GetActiveTicket(car.LicensePlate); err == nil {
						info.AttendantID = ticket.AttendantID
					}

					whiteCars = append(whiteCars, info)
				}
			}
		}
	}

	return whiteCars, nil
}

// UC13: Find all blue Toyota cars for robbery investigation
func (ps *PoliceService) FindBlueToyotaCars() ([]*VehicleInvestigationInfo, error) {
	var blueToyotas []*VehicleInvestigationInfo

	for _, lot := range ps.parkingService.lots {
		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil {
				car := space.ParkedCar
				if strings.ToLower(car.Color) == "blue" && strings.ToLower(car.Make) == "toyota" {
					info := &VehicleInvestigationInfo{
						Car:      car,
						LotID:    lot.ID,
						SpaceID:  fmt.Sprintf("%d", space.ID), // FIXED: Convert int to string
						ParkedAt: space.ParkedAt,
					}

					// Get attendant info for robbery investigation
					if ticket, err := ps.parkingService.GetActiveTicket(car.LicensePlate); err == nil {
						info.AttendantID = ticket.AttendantID

						// Find attendant name
						if attendant := ps.parkingService.FindAttendantByID(ticket.AttendantID); attendant != nil {
							info.AttendantName = attendant.Name
						}
					}

					blueToyotas = append(blueToyotas, info)
				}
			}
		}
	}

	return blueToyotas, nil
}

// General police query for any color/make combination
func (ps *PoliceService) FindCarsByColorAndMake(color, make string) ([]*VehicleInvestigationInfo, error) {
	var matchingCars []*VehicleInvestigationInfo

	for _, lot := range ps.parkingService.lots {
		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil {
				car := space.ParkedCar
				colorMatch := color == "" || strings.ToLower(car.Color) == strings.ToLower(color)
				makeMatch := make == "" || strings.ToLower(car.Make) == strings.ToLower(make)

				if colorMatch && makeMatch {
					info := &VehicleInvestigationInfo{
						Car:      car,
						LotID:    lot.ID,
						SpaceID:  fmt.Sprintf("%d", space.ID), // FIXED: Convert int to string
						ParkedAt: space.ParkedAt,
					}

					if ticket, err := ps.parkingService.GetActiveTicket(car.LicensePlate); err == nil {
						info.AttendantID = ticket.AttendantID
						if attendant := ps.parkingService.FindAttendantByID(ticket.AttendantID); attendant != nil {
							info.AttendantName = attendant.Name
						}
					}

					matchingCars = append(matchingCars, info)
				}
			}
		}
	}

	return matchingCars, nil
}

func (ps *PoliceService) GenerateInvestigationReport(vehicles []*VehicleInvestigationInfo, caseType string) string {
	report := "=== POLICE INVESTIGATION REPORT ===\n"
	report += "Case Type: " + caseType + "\n"
	report += "Generated: " + time.Now().Format("2006-01-02 15:04:05") + "\n"
	report += fmt.Sprintf("Total Vehicles Found: %d\n\n", len(vehicles))

	for i, vehicle := range vehicles {
		report += fmt.Sprintf("Vehicle %d:\n", i+1)
		report += "  License Plate: " + vehicle.Car.LicensePlate + "\n"
		report += "  Driver: " + vehicle.Car.DriverName + "\n"
		report += "  Color: " + vehicle.Car.Color + "\n"
		report += "  Make: " + vehicle.Car.Make + "\n"
		report += "  Location: Lot " + vehicle.LotID + ", Space " + vehicle.SpaceID + "\n"
		report += "  Parked At: " + vehicle.ParkedAt.Format("2006-01-02 15:04:05") + "\n"
		if vehicle.AttendantID != "" {
			report += "  Attendant ID: " + vehicle.AttendantID + "\n"
			if vehicle.AttendantName != "" {
				report += "  Attendant Name: " + vehicle.AttendantName + "\n"
			}
		}
		report += "\n"
	}

	return report
}
