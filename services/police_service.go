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

// UC13: Robbery investigation specific methods (CLEAN IMPLEMENTATION)
func (ps *PoliceService) GenerateRobberyInvestigationReport(suspectDescription string) string {
	blueToyotas, err := ps.FindBlueToyotaCars()
	if err != nil {
		return "Error generating robbery investigation report: " + err.Error()
	}

	report := "=== ROBBERY INVESTIGATION REPORT ===\n"
	report += "Case Type: Armed Robbery\n"
	report += "Target Vehicle: Blue Toyota\n"
	if suspectDescription != "" {
		report += "Suspect Description: " + suspectDescription + "\n"
	}
	report += "Generated: " + time.Now().Format("2006-01-02 15:04:05") + "\n"
	report += fmt.Sprintf("Blue Toyota Vehicles Found: %d\n\n", len(blueToyotas))

	for i, vehicle := range blueToyotas {
		report += fmt.Sprintf("SUSPECT VEHICLE %d:\n", i+1)
		report += "  License Plate: " + vehicle.Car.LicensePlate + "\n"
		report += "  Driver Name: " + vehicle.Car.DriverName + "\n"
		report += "  Location: Lot " + vehicle.LotID + ", Space " + vehicle.SpaceID + "\n"
		report += "  Time Parked: " + vehicle.ParkedAt.Format("2006-01-02 15:04:05") + "\n"

		if vehicle.AttendantID != "" {
			report += "  Parking Attendant: " + vehicle.AttendantName + " (ID: " + vehicle.AttendantID + ")\n"
			report += "  Attendant Status: Active\n"
		}

		report += "  Vehicle Details:\n"
		report += "    Size: " + vehicle.Car.GetVehicleSizeString() + "\n"
		if vehicle.Car.IsHandicap {
			report += "    Special Status: Handicap Vehicle\n"
		}

		report += "\n"
	}

	if len(blueToyotas) == 0 {
		report += "NO MATCHING VEHICLES FOUND\n"
		report += "Recommendation: Expand search parameters\n"
	} else {
		report += "INVESTIGATION RECOMMENDATIONS:\n"
		report += "1. Interview all identified drivers\n"
		report += "2. Review parking attendant logs\n"
		report += "3. Check security footage for time periods listed above\n"
		report += "4. Verify attendant identification and background\n"
	}

	return report
}

func (ps *PoliceService) GetBlueToyotaCount() int {
	blueToyotas, err := ps.FindBlueToyotaCars()
	if err != nil {
		return 0
	}
	return len(blueToyotas)
}

func (ps *PoliceService) GetSuspectVehicleDetails(licensePlate string) map[string]interface{} {
	details := make(map[string]interface{})

	blueToyotas, err := ps.FindBlueToyotaCars()
	if err != nil {
		details["error"] = err.Error()
		return details
	}

	for _, vehicle := range blueToyotas {
		if vehicle.Car.LicensePlate == licensePlate {
			details["found"] = true
			details["licensePlate"] = vehicle.Car.LicensePlate
			details["driverName"] = vehicle.Car.DriverName
			details["color"] = vehicle.Car.Color
			details["make"] = vehicle.Car.Make
			details["lotID"] = vehicle.LotID
			details["spaceID"] = vehicle.SpaceID
			details["parkedAt"] = vehicle.ParkedAt.Format("2006-01-02 15:04:05")
			details["attendantID"] = vehicle.AttendantID
			details["attendantName"] = vehicle.AttendantName
			details["vehicleSize"] = vehicle.Car.GetVehicleSizeString()
			details["isHandicap"] = vehicle.Car.IsHandicap
			return details
		}
	}

	details["found"] = false
	details["message"] = "Vehicle not found in blue Toyota suspects"
	return details
}

func (ps *PoliceService) ValidateRobberyEvidence() map[string]interface{} {
	evidence := make(map[string]interface{})

	blueToyotas, err := ps.FindBlueToyotaCars()
	if err != nil {
		evidence["error"] = err.Error()
		return evidence
	}

	evidence["totalSuspectVehicles"] = len(blueToyotas)
	evidence["caseStrength"] = func() string {
		if len(blueToyotas) == 0 {
			return "Weak - No suspect vehicles found"
		} else if len(blueToyotas) == 1 {
			return "Strong - Single suspect vehicle identified"
		} else {
			return "Moderate - Multiple suspect vehicles require investigation"
		}
	}()

	var attendantCount int
	for _, vehicle := range blueToyotas {
		if vehicle.AttendantID != "" {
			attendantCount++
		}
	}

	evidence["attendantWitnesses"] = attendantCount
	evidence["evidenceQuality"] = func() string {
		if attendantCount == len(blueToyotas) {
			return "High - All vehicles have attendant witnesses"
		} else if attendantCount > 0 {
			return "Medium - Some attendant witnesses available"
		} else {
			return "Low - No attendant witnesses"
		}
	}()

	return evidence
}
