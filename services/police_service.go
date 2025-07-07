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

// UC14: Find all BMW cars for suspicious activity monitoring
func (ps *PoliceService) FindBMWCars() ([]*VehicleInvestigationInfo, error) {
	var bmwCars []*VehicleInvestigationInfo

	for _, lot := range ps.parkingService.lots {
		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil {
				car := space.ParkedCar
				if strings.ToLower(car.Make) == "bmw" {
					info := &VehicleInvestigationInfo{
						Car:      car,
						LotID:    lot.ID,
						SpaceID:  fmt.Sprintf("%d", space.ID),
						ParkedAt: space.ParkedAt,
					}

					// Get attendant info if available
					if ticket, err := ps.parkingService.GetActiveTicket(car.LicensePlate); err == nil {
						info.AttendantID = ticket.AttendantID
						if attendant := ps.parkingService.FindAttendantByID(ticket.AttendantID); attendant != nil {
							info.AttendantName = attendant.Name
						}
					}

					bmwCars = append(bmwCars, info)
				}
			}
		}
	}

	return bmwCars, nil
}

// UC14: Generate BMW security monitoring report
func (ps *PoliceService) GenerateBMWSecurityReport() string {
	bmwCars, err := ps.FindBMWCars()
	if err != nil {
		return "Error generating BMW security report: " + err.Error()
	}

	report := "=== BMW SECURITY MONITORING REPORT ===\n"
	report += "Alert Type: Suspicious Activity Monitoring\n"
	report += "Target Vehicle: BMW (All Models)\n"
	report += "Security Level: Enhanced\n"
	report += "Generated: " + time.Now().Format("2006-01-02 15:04:05") + "\n"
	report += fmt.Sprintf("BMW Vehicles Found: %d\n\n", len(bmwCars))

	for i, vehicle := range bmwCars {
		report += fmt.Sprintf("BMW VEHICLE %d:\n", i+1)
		report += "  License Plate: " + vehicle.Car.LicensePlate + "\n"
		report += "  Driver Name: " + vehicle.Car.DriverName + "\n"
		report += "  Color: " + vehicle.Car.Color + "\n"
		report += "  Location: Lot " + vehicle.LotID + ", Space " + vehicle.SpaceID + "\n"
		report += "  Time Parked: " + vehicle.ParkedAt.Format("2006-01-02 15:04:05") + "\n"
		report += "  Vehicle Size: " + vehicle.Car.GetVehicleSizeString() + "\n"

		if vehicle.AttendantID != "" {
			report += "  Parking Attendant: " + vehicle.AttendantName + " (ID: " + vehicle.AttendantID + ")\n"
		}

		if vehicle.Car.IsHandicap {
			report += "  Special Status: Handicap Vehicle\n"
		}

		report += "  Security Priority: HIGH\n"
		report += "\n"
	}

	if len(bmwCars) == 0 {
		report += "NO BMW VEHICLES FOUND\n"
		report += "Status: All Clear - No BMW vehicles in parking system\n"
	} else {
		report += "SECURITY RECOMMENDATIONS:\n"
		report += "1. Increase patrol frequency around BMW vehicle locations\n"
		report += "2. Monitor for unusual activity near these vehicles\n"
		report += "3. Alert security staff to enhanced surveillance\n"
		report += "4. Cross-reference with incident reports\n"
		report += "5. Document all BMW vehicle movements\n"
	}

	return report
}

// UC14: Get BMW count for security dashboard
func (ps *PoliceService) GetBMWCount() int {
	bmwCars, err := ps.FindBMWCars()
	if err != nil {
		return 0
	}
	return len(bmwCars)
}

// UC14: Get BMW vehicles by security priority
func (ps *PoliceService) GetBMWVehiclesByPriority() map[string]interface{} {
	bmwCars, err := ps.FindBMWCars()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	result := make(map[string]interface{})
	result["totalBMWVehicles"] = len(bmwCars)
	result["securityLevel"] = "HIGH"
	result["requiresEnhancedSecurity"] = len(bmwCars) > 0

	var highPriority, mediumPriority, lowPriority int
	for _, vehicle := range bmwCars {
		if vehicle.Car.IsHandicap || vehicle.Car.Size == models.LargeVehicle {
			highPriority++
		} else if vehicle.Car.Size == models.MediumVehicle {
			mediumPriority++
		} else {
			lowPriority++
		}
	}

	result["highPriority"] = highPriority
	result["mediumPriority"] = mediumPriority
	result["lowPriority"] = lowPriority

	return result
}

// UC14: Validate BMW security protocols
func (ps *PoliceService) ValidateBMWSecurityProtocols() map[string]interface{} {
	validation := make(map[string]interface{})

	bmwCars, err := ps.FindBMWCars()
	if err != nil {
		validation["error"] = err.Error()
		return validation
	}

	validation["totalBMWVehicles"] = len(bmwCars)
	validation["securityProtocolActive"] = len(bmwCars) > 0

	var attendantCoverage int
	for _, vehicle := range bmwCars {
		if vehicle.AttendantID != "" {
			attendantCoverage++
		}
	}

	validation["attendantCoverage"] = attendantCoverage
	validation["coverageQuality"] = func() string {
		if len(bmwCars) == 0 {
			return "Not Applicable"
		}
		coverage := float64(attendantCoverage) / float64(len(bmwCars)) * 100
		if coverage == 100 {
			return "Excellent - All BMW vehicles have attendant records"
		} else if coverage >= 80 {
			return "Good - Most BMW vehicles have attendant records"
		} else {
			return "Needs Improvement - Limited attendant coverage"
		}
	}()

	return validation
}

// General police query methods (UC12-UC13 support)
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
						SpaceID:  fmt.Sprintf("%d", space.ID),
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

// UC13: Robbery investigation specific methods
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

// UC15: Find cars parked in the last specified minutes
func (ps *PoliceService) FindCarsParkedInLastMinutes(minutes int) ([]*VehicleInvestigationInfo, error) {
	var recentCars []*VehicleInvestigationInfo
	cutoffTime := time.Now().Add(-time.Duration(minutes) * time.Minute)

	for _, lot := range ps.parkingService.lots {
		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil {
				if space.ParkedAt.After(cutoffTime) {
					info := &VehicleInvestigationInfo{
						Car:      space.ParkedCar,
						LotID:    lot.ID,
						SpaceID:  fmt.Sprintf("%d", space.ID),
						ParkedAt: space.ParkedAt,
					}

					if ticket, err := ps.parkingService.GetActiveTicket(space.ParkedCar.LicensePlate); err == nil {
						info.AttendantID = ticket.AttendantID
						if attendant := ps.parkingService.FindAttendantByID(ticket.AttendantID); attendant != nil {
							info.AttendantName = attendant.Name
						}
					}

					recentCars = append(recentCars, info)
				}
			}
		}
	}

	return recentCars, nil
}

// UC15: Get recent parking activity with flexible time range
func (ps *PoliceService) GetRecentParkingActivity(timeRange time.Duration) ([]*VehicleInvestigationInfo, error) {
	var recentActivity []*VehicleInvestigationInfo
	cutoffTime := time.Now().Add(-timeRange)

	for _, lot := range ps.parkingService.lots {
		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil {
				if space.ParkedAt.After(cutoffTime) {
					info := &VehicleInvestigationInfo{
						Car:      space.ParkedCar,
						LotID:    lot.ID,
						SpaceID:  fmt.Sprintf("%d", space.ID),
						ParkedAt: space.ParkedAt,
					}

					if ticket, err := ps.parkingService.GetActiveTicket(space.ParkedCar.LicensePlate); err == nil {
						info.AttendantID = ticket.AttendantID
						if attendant := ps.parkingService.FindAttendantByID(ticket.AttendantID); attendant != nil {
							info.AttendantName = attendant.Name
						}
					}

					recentActivity = append(recentActivity, info)
				}
			}
		}
	}

	return recentActivity, nil
}

// UC15: Generate time-based investigation report for bomb threats
func (ps *PoliceService) GenerateTimeBasedInvestigationReport(minutes int) string {
	recentCars, err := ps.FindCarsParkedInLastMinutes(minutes)
	if err != nil {
		return "Error generating time-based investigation report: " + err.Error()
	}

	report := "=== TIME-BASED INVESTIGATION REPORT ===\n"
	report += "Alert Type: Bomb Threat Investigation\n"
	report += "Time Range: Last " + fmt.Sprintf("%d", minutes) + " minutes\n"
	report += "Generated: " + time.Now().Format("2006-01-02 15:04:05") + "\n"
	report += fmt.Sprintf("Vehicles Found: %d\n\n", len(recentCars))

	for i, vehicle := range recentCars {
		timeSinceParked := time.Since(vehicle.ParkedAt)
		report += fmt.Sprintf("RECENT VEHICLE %d:\n", i+1)
		report += "  License Plate: " + vehicle.Car.LicensePlate + "\n"
		report += "  Driver Name: " + vehicle.Car.DriverName + "\n"
		report += "  Color: " + vehicle.Car.Color + "\n"
		report += "  Make: " + vehicle.Car.Make + "\n"
		report += "  Location: Lot " + vehicle.LotID + ", Space " + vehicle.SpaceID + "\n"
		report += "  Parked At: " + vehicle.ParkedAt.Format("2006-01-02 15:04:05") + "\n"
		report += "  Time Since Parked: " + fmt.Sprintf("%.0f minutes ago", timeSinceParked.Minutes()) + "\n"

		if vehicle.AttendantID != "" {
			report += "  Parking Attendant: " + vehicle.AttendantName + " (ID: " + vehicle.AttendantID + ")\n"
		}

		if vehicle.Car.IsHandicap {
			report += "  Special Status: Handicap Vehicle\n"
		}

		report += "  Security Priority: HIGH\n"
		report += "\n"
	}

	if len(recentCars) == 0 {
		report += "NO RECENT VEHICLES FOUND\n"
		report += "Status: All vehicles parked more than " + fmt.Sprintf("%d", minutes) + " minutes ago\n"
	} else {
		report += "INVESTIGATION RECOMMENDATIONS:\n"
		report += "1. Immediately inspect all recent vehicle arrivals\n"
		report += "2. Review security footage for the specified time range\n"
		report += "3. Interview parking attendants on duty\n"
		report += "4. Coordinate with bomb disposal unit if threats are credible\n"
		report += "5. Evacuate areas if necessary based on threat assessment\n"
	}

	return report
}

// UC15: Get vehicle count for specific time windows
func (ps *PoliceService) GetVehicleCountByTimeWindow(minutes int) map[string]interface{} {
	result := make(map[string]interface{})

	recentCars, err := ps.FindCarsParkedInLastMinutes(minutes)
	if err != nil {
		result["error"] = err.Error()
		return result
	}

	result["timeWindow"] = fmt.Sprintf("%d minutes", minutes)
	result["totalVehicles"] = len(recentCars)
	result["timestamp"] = time.Now().Format("2006-01-02 15:04:05")

	// Categorize by time buckets
	var last15min, last30min, older int
	now := time.Now()

	for _, vehicle := range recentCars {
		minutesAgo := int(now.Sub(vehicle.ParkedAt).Minutes())
		if minutesAgo <= 15 {
			last15min++
		} else if minutesAgo <= 30 {
			last30min++
		} else {
			older++
		}
	}

	result["last15Minutes"] = last15min
	result["last30Minutes"] = last30min
	result["olderThanRequested"] = older

	return result
}
