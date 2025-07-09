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

// UC16: Find handicap cars in specific rows
func (ps *PoliceService) FindHandicapCarsInRows(rows []string) ([]*VehicleInvestigationInfo, error) {
	var handicapCars []*VehicleInvestigationInfo

	for _, lot := range ps.parkingService.lots {
		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil && space.ParkedCar.IsHandicap {
				spaceRow := space.GetRowAssignment()
				
				// Check if this space is in one of the requested rows
				for _, targetRow := range rows {
					if strings.ToUpper(spaceRow) == strings.ToUpper(targetRow) {
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

						handicapCars = append(handicapCars, info)
						break
					}
				}
			}
		}
	}

	return handicapCars, nil
}

// UC16: Get vehicles by location criteria (size, handicap status, rows)
func (ps *PoliceService) GetVehiclesByLocationCriteria(size models.VehicleSize, handicapOnly bool, rows []string) ([]*VehicleInvestigationInfo, error) {
	var matchingVehicles []*VehicleInvestigationInfo

	for _, lot := range ps.parkingService.lots {
		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil {
				car := space.ParkedCar
				
				// Check criteria
				sizeMatch := (size == models.SmallVehicle && car.Size == models.SmallVehicle) ||
					(size == models.MediumVehicle && car.Size == models.MediumVehicle) ||
					(size == models.LargeVehicle && car.Size == models.LargeVehicle)
				
				handicapMatch := !handicapOnly || car.IsHandicap
				
				spaceRow := space.GetRowAssignment()
				rowMatch := len(rows) == 0
				for _, targetRow := range rows {
					if strings.ToUpper(spaceRow) == strings.ToUpper(targetRow) {
						rowMatch = true
						break
					}
				}

				if sizeMatch && handicapMatch && rowMatch {
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

					matchingVehicles = append(matchingVehicles, info)
				}
			}
		}
	}

	return matchingVehicles, nil
}

// UC16: Validate handicap permit fraud
func (ps *PoliceService) ValidateHandicapPermitFraud() map[string]interface{} {
	validation := make(map[string]interface{})

	// Find all handicap cars
	var allHandicapCars []*VehicleInvestigationInfo
	for _, lot := range ps.parkingService.lots {
		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil && space.ParkedCar.IsHandicap {
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

				allHandicapCars = append(allHandicapCars, info)
			}
		}
	}

	// Find handicap cars in suspicious rows (B and D)
	suspiciousRows, _ := ps.FindHandicapCarsInRows([]string{"B", "D"})

	validation["totalHandicapVehicles"] = len(allHandicapCars)
	validation["vehiclesInRowsB_D"] = len(suspiciousRows)
	validation["fraudRisk"] = func() string {
		if len(allHandicapCars) == 0 {
			return "No Assessment - No handicap vehicles found"
		}
		
		percentage := float64(len(suspiciousRows)) / float64(len(allHandicapCars)) * 100
		
		if percentage > 50 {
			return "High Risk - Majority of handicap vehicles in target rows"
		} else if percentage > 25 {
			return "Medium Risk - Significant number in target rows"
		} else if percentage > 0 {
			return "Low Risk - Some vehicles in target rows require verification"
		} else {
			return "No Risk - No handicap vehicles in suspicious rows"
		}
	}()

	validation["investigationRequired"] = len(suspiciousRows) > 0
	validation["timestamp"] = time.Now().Format("2006-01-02 15:04:05")

	return validation
}

// UC16: Generate handicap permit fraud investigation report
func (ps *PoliceService) GenerateHandicapFraudInvestigationReport() string {
	suspiciousVehicles, err := ps.FindHandicapCarsInRows([]string{"B", "D"})
	if err != nil {
		return "Error generating handicap fraud investigation report: " + err.Error()
	}

	report := "=== HANDICAP PERMIT FRAUD INVESTIGATION REPORT ===\n"
	report += "Investigation Type: Handicap Permit Fraud\n"
	report += "Target Locations: Rows B and D\n"
	report += "Generated: " + time.Now().Format("2006-01-02 15:04:05") + "\n"
	report += fmt.Sprintf("Suspicious Vehicles Found: %d\n\n", len(suspiciousVehicles))

	for i, vehicle := range suspiciousVehicles {
		row := ""
		if space := ps.findParkingSpace(vehicle.LotID, vehicle.SpaceID); space != nil {
			row = space.GetRowAssignment()
		}

		report += fmt.Sprintf("SUSPICIOUS VEHICLE %d:\n", i+1)
		report += "  License Plate: " + vehicle.Car.LicensePlate + "\n"
		report += "  Driver Name: " + vehicle.Car.DriverName + "\n"
		report += "  Vehicle Size: " + vehicle.Car.GetVehicleSizeString() + "\n"
		report += "  Location: Lot " + vehicle.LotID + ", Row " + row + ", Space " + vehicle.SpaceID + "\n"
		report += "  Parked At: " + vehicle.ParkedAt.Format("2006-01-02 15:04:05") + "\n"

		if vehicle.AttendantID != "" {
			report += "  Parking Attendant: " + vehicle.AttendantName + " (ID: " + vehicle.AttendantID + ")\n"
		}

		report += "  Fraud Risk: HIGH - Handicap vehicle in suspicious row\n"
		report += "\n"
	}

	if len(suspiciousVehicles) == 0 {
		report += "NO SUSPICIOUS VEHICLES FOUND\n"
		report += "Status: All handicap vehicles appear to be in legitimate locations\n"
	} else {
		report += "INVESTIGATION RECOMMENDATIONS:\n"
		report += "1. Verify handicap permits for all identified vehicles\n"
		report += "2. Interview parking attendants who processed these vehicles\n"
		report += "3. Cross-reference permit numbers with official database\n"
		report += "4. Review security footage for permit verification process\n"
		report += "5. Conduct field verification of permit authenticity\n"
	}

	return report
}

// UC16: Helper method to find parking space by lot and space ID
func (ps *PoliceService) findParkingSpace(lotID, spaceID string) *models.ParkingSpace {
	for _, lot := range ps.parkingService.lots {
		if lot.ID == lotID {
			for _, space := range lot.Spaces {
				if fmt.Sprintf("%d", space.ID) == spaceID {
					return space
				}
			}
		}
	}
	return nil
}

// UC16: Get location statistics for analysis
func (ps *PoliceService) GetLocationStatistics() map[string]interface{} {
	stats := make(map[string]interface{})
	
	rowCounts := map[string]int{"A": 0, "B": 0, "C": 0, "D": 0}
	handicapByRow := map[string]int{"A": 0, "B": 0, "C": 0, "D": 0}
	sizeByRow := map[string]map[string]int{
		"A": {"Small": 0, "Medium": 0, "Large": 0},
		"B": {"Small": 0, "Medium": 0, "Large": 0},
		"C": {"Small": 0, "Medium": 0, "Large": 0},
		"D": {"Small": 0, "Medium": 0, "Large": 0},
	}

	for _, lot := range ps.parkingService.lots {
		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil {
				row := space.GetRowAssignment()
				rowCounts[row]++

				if space.ParkedCar.IsHandicap {
					handicapByRow[row]++
				}

				sizeStr := space.ParkedCar.GetVehicleSizeString()
				sizeByRow[row][sizeStr]++
			}
		}
	}

	stats["totalVehiclesByRow"] = rowCounts
	stats["handicapVehiclesByRow"] = handicapByRow
	stats["vehicleSizesByRow"] = sizeByRow
	stats["timestamp"] = time.Now().Format("2006-01-02 15:04:05")

	return stats
}

// UC17: Get all cars parked in a specific parking lot
func (ps *PoliceService) GetAllCarsInLot(lotID string) ([]*VehicleInvestigationInfo, error) {
	var allCars []*VehicleInvestigationInfo

	for _, lot := range ps.parkingService.lots {
		if lot.ID == lotID {
			for _, space := range lot.Spaces {
				if space.IsOccupied && space.ParkedCar != nil {
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

					allCars = append(allCars, info)
				}
			}
			return allCars, nil
		}
	}

	return nil, fmt.Errorf("parking lot %s not found", lotID)
}

// UC17: Detect potentially fraudulent license plates
func (ps *PoliceService) DetectFraudulentPlates() ([]*VehicleInvestigationInfo, error) {
	var suspiciousVehicles []*VehicleInvestigationInfo

	for _, lot := range ps.parkingService.lots {
		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil {
				car := space.ParkedCar
				
				// Check for suspicious patterns in license plates
				if ps.isSuspiciousLicensePlate(car.LicensePlate) {
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

					suspiciousVehicles = append(suspiciousVehicles, info)
				}
			}
		}
	}

	return suspiciousVehicles, nil
}

// UC17: Check if a license plate is potentially fraudulent
func (ps *PoliceService) isSuspiciousLicensePlate(plate string) bool {
	plate = strings.ToUpper(strings.TrimSpace(plate))
	
	// Check for suspicious patterns
	patterns := []string{
		"FAKE", "TEST", "TEMP", "XXXX", "0000", "1111", 
		"FRAUD", "STOLEN", "NULL", "ADMIN", "DEBUG",
	}
	
	for _, pattern := range patterns {
		if strings.Contains(plate, pattern) {
			return true
		}
	}
	
	// Check for excessive repetition (e.g., "AAA111", "BBB222")
	if len(plate) >= 6 {
		letters := 0
		numbers := 0
		for _, char := range plate {
			if char >= 'A' && char <= 'Z' {
				letters++
			} else if char >= '0' && char <= '9' {
				numbers++
			}
		}
		
		// Suspicious if all letters or all numbers
		if letters == len(plate) || numbers == len(plate) {
			return true
		}
	}
	
	// Check for sequential patterns
	if ps.hasSequentialPattern(plate) {
		return true
	}
	
	return false
}

// UC17: Check for sequential patterns in license plates
func (ps *PoliceService) hasSequentialPattern(plate string) bool {
	if len(plate) < 3 {
		return false
	}
	
	consecutive := 0
	for i := 1; i < len(plate); i++ {
		if plate[i] == plate[i-1]+1 {
			consecutive++
			if consecutive >= 2 { // Three consecutive characters
				return true
			}
		} else {
			consecutive = 0
		}
	}
	
	return false
}

// UC17: Generate complete parking lot investigation report
func (ps *PoliceService) GenerateCompleteLotInvestigationReport(lotID string) string {
	allCars, err := ps.GetAllCarsInLot(lotID)
	if err != nil {
		return "Error generating complete lot investigation report: " + err.Error()
	}

	fraudulentCars, _ := ps.DetectFraudulentPlates()
	
	// Filter fraudulent cars for this lot
	var lotFraudulentCars []*VehicleInvestigationInfo
	for _, car := range fraudulentCars {
		if car.LotID == lotID {
			lotFraudulentCars = append(lotFraudulentCars, car)
		}
	}

	report := "=== COMPLETE PARKING LOT INVESTIGATION REPORT ===\n"
	report += "Investigation Type: Complete Lot Analysis\n"
	report += "Target Lot: " + lotID + "\n"
	report += "Investigation Focus: Fraudulent License Plates\n"
	report += "Generated: " + time.Now().Format("2006-01-02 15:04:05") + "\n"
	report += fmt.Sprintf("Total Vehicles in Lot: %d\n", len(allCars))
	report += fmt.Sprintf("Suspicious Vehicles Found: %d\n\n", len(lotFraudulentCars))

	// All vehicles section
	report += "=== ALL VEHICLES IN LOT " + lotID + " ===\n"
	for i, vehicle := range allCars {
		report += fmt.Sprintf("VEHICLE %d:\n", i+1)
		report += "  License Plate: " + vehicle.Car.LicensePlate + "\n"
		report += "  Driver Name: " + vehicle.Car.DriverName + "\n"
		report += "  Color: " + vehicle.Car.Color + "\n"
		report += "  Make: " + vehicle.Car.Make + "\n"
		report += "  Vehicle Size: " + vehicle.Car.GetVehicleSizeString() + "\n"
		report += "  Space: " + vehicle.SpaceID + "\n"
		report += "  Parked At: " + vehicle.ParkedAt.Format("2006-01-02 15:04:05") + "\n"

		if vehicle.AttendantID != "" {
			report += "  Parking Attendant: " + vehicle.AttendantName + " (ID: " + vehicle.AttendantID + ")\n"
		}

		if vehicle.Car.IsHandicap {
			report += "  Special Status: Handicap Vehicle\n"
		}

		// Check if this vehicle is suspicious
		isSuspicious := false
		for _, suspicious := range lotFraudulentCars {
			if suspicious.Car.LicensePlate == vehicle.Car.LicensePlate {
				isSuspicious = true
				break
			}
		}

		if isSuspicious {
			report += "  ⚠️  FRAUD ALERT: Suspicious license plate detected\n"
		} else {
			report += "  Status: Normal\n"
		}

		report += "\n"
	}

	// Fraudulent vehicles section
	if len(lotFraudulentCars) > 0 {
		report += "=== FRAUDULENT VEHICLES INVESTIGATION ===\n"
		for i, vehicle := range lotFraudulentCars {
			report += fmt.Sprintf("SUSPICIOUS VEHICLE %d:\n", i+1)
			report += "  License Plate: " + vehicle.Car.LicensePlate + " ⚠️ FLAGGED\n"
			report += "  Driver Name: " + vehicle.Car.DriverName + "\n"
			report += "  Location: Space " + vehicle.SpaceID + "\n"
			report += "  Parked At: " + vehicle.ParkedAt.Format("2006-01-02 15:04:05") + "\n"

			if vehicle.AttendantID != "" {
				report += "  Attendant: " + vehicle.AttendantName + " (ID: " + vehicle.AttendantID + ")\n"
			}

			report += "  Fraud Risk: HIGH - Requires immediate verification\n"
			report += "\n"
		}

		report += "IMMEDIATE ACTIONS REQUIRED:\n"
		report += "1. Verify identity of all flagged drivers\n"
		report += "2. Cross-reference license plates with national database\n"
		report += "3. Interview parking attendants who processed these vehicles\n"
		report += "4. Review security footage for parking time periods\n"
		report += "5. Contact DMV for license plate authenticity verification\n"
		report += "6. Prepare evidence for potential legal proceedings\n"
	} else {
		report += "=== FRAUD ANALYSIS RESULTS ===\n"
		report += "✅ NO FRAUDULENT PLATES DETECTED\n"
		report += "Status: All vehicles in lot appear legitimate\n"
		report += "Recommendation: Continue routine monitoring\n"
	}

	return report
}

// UC17: Get fraud statistics for all lots
func (ps *PoliceService) GetFraudStatistics() map[string]interface{} {
	stats := make(map[string]interface{})

	allFraudulent, err := ps.DetectFraudulentPlates()
	if err != nil {
		stats["error"] = err.Error()
		return stats
	}

	stats["totalSuspiciousVehicles"] = len(allFraudulent)
	stats["timestamp"] = time.Now().Format("2006-01-02 15:04:05")

	// Group by lot
	lotStats := make(map[string]int)
	for _, vehicle := range allFraudulent {
		lotStats[vehicle.LotID]++
	}

	stats["fraudByLot"] = lotStats

	// Calculate fraud rate
	totalVehicles := 0
	for _, lot := range ps.parkingService.lots {
		totalVehicles += lot.GetOccupiedSpaces()
	}

	if totalVehicles > 0 {
		fraudRate := float64(len(allFraudulent)) / float64(totalVehicles) * 100
		stats["fraudRate"] = fraudRate
	} else {
		stats["fraudRate"] = 0.0
	}

	return stats
}

// UC17: Get investigation summary for all implemented use cases
func (ps *PoliceService) GetCompleteInvestigationSummary() map[string]interface{} {
	summary := make(map[string]interface{})

	// UC12: White cars (bomb threat)
	whiteCars, _ := ps.FindWhiteCars()
	summary["whiteCarsCount"] = len(whiteCars)

	// UC13: Blue Toyotas (robbery)
	blueToyotas, _ := ps.FindBlueToyotaCars()
	summary["blueToyotasCount"] = len(blueToyotas)

	// UC14: BMW cars (suspicious activity)
	bmwCars, _ := ps.FindBMWCars()
	summary["bmwCarsCount"] = len(bmwCars)

	// UC15: Recent cars (bomb threat - time-based)
	recentCars, _ := ps.FindCarsParkedInLastMinutes(30)
	summary["recentCarsCount"] = len(recentCars)

	// UC16: Handicap fraud (location-based)
	handicapFraud, _ := ps.FindHandicapCarsInRows([]string{"B", "D"})
	summary["handicapFraudCount"] = len(handicapFraud)

	// UC17: Fraudulent plates (complete investigation)
	fraudulentCars, _ := ps.DetectFraudulentPlates()
	summary["fraudulentPlatesCount"] = len(fraudulentCars)

	summary["investigationTimestamp"] = time.Now().Format("2006-01-02 15:04:05")
	summary["totalInvestigationTypes"] = 6 // UC12-UC17

	return summary
}
