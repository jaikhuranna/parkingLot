package main

import (
	"fmt"
	"time"
	"parking-lot-system/interfaces"
	"parking-lot-system/models"
	"parking-lot-system/services"
)

func main() {
	fmt.Println("Welcome to Parking Lot System!")
	fmt.Println("UC15-UC16: Time-based Queries and Location-based Searches")
	fmt.Println("========================================================")

	// Create parking lots
	lot1 := models.NewParkingLot("LOT1", 100)
	lot2 := models.NewParkingLot("LOT2", 100)
	service := services.NewParkingService()
	service.AddLot(lot1)
	service.AddLot(lot2)

	// Add attendants
	attendant := models.NewParkingAttendant("ATT001", "Alice Johnson", "LOT1")
	service.AddAttendant(attendant)

	// Add observers
	owner := interfaces.NewOwnerObserver("Sanjay")
	service.AddObserverToLot("LOT1", owner)
	service.AddObserverToLot("LOT2", owner)

	// Create police service
	policeService := services.NewPoliceService(service)

	// UC15-UC16 Demo: Add test vehicles with specific timing and locations
	fmt.Println("\n🚗 Setting up test scenario for UC15 & UC16...")

	// Recent vehicles for UC15 testing
	recentCars := []*models.Car{
		func() *models.Car {
			car := models.NewCar("RECENT001", "John Smith")
			car.SetColor("Red")
			car.SetMake("Honda")
			return car
		}(),
		func() *models.Car {
			car := models.NewCar("RECENT002", "Jane Doe")
			car.SetColor("Blue")
			car.SetMake("Toyota")
			return car
		}(),
	}

	// Handicap vehicles for UC16 testing (placed in specific rows)
	handicapCars := []*models.Car{
		func() *models.Car {
			car := models.NewCar("HANDICAP_B1", "Handicap Driver B1")
			car.SetHandicapStatus(true)
			car.SetVehicleSize(models.SmallVehicle)
			car.SetColor("White")
			return car
		}(),
		func() *models.Car {
			car := models.NewCar("HANDICAP_D1", "Handicap Driver D1")
			car.SetHandicapStatus(true)
			car.SetVehicleSize(models.SmallVehicle)
			car.SetColor("Silver")
			return car
		}(),
		func() *models.Car {
			car := models.NewCar("HANDICAP_A1", "Legitimate Handicap")
			car.SetHandicapStatus(true)
			car.SetVehicleSize(models.MediumVehicle)
			car.SetColor("Black")
			return car
		}(),
	}

	// Park recent cars
	fmt.Println("\n📍 Parking recent vehicles...")
	for i, car := range recentCars {
		_, err := service.ParkCarWithTicket(car)
		if err != nil {
			fmt.Printf("Warning: Could not park %s: %v\n", car.LicensePlate, err)
		} else {
			fmt.Printf("✅ Parked recent car %d: %s\n", i+1, car.LicensePlate)
		}
		time.Sleep(100 * time.Millisecond) // Small delay between parking
	}

	// Park handicap cars in specific spaces to control row assignment
	fmt.Println("\n📍 Parking handicap vehicles in specific rows...")
	
	// Park in Row B (spaces 26-50)
	lot1.Spaces[30].Park(handicapCars[0]) // Space 31 = Row B
	fmt.Printf("✅ Parked %s in Row B (Space 31)\n", handicapCars[0].LicensePlate)
	
	// Park in Row D (spaces 76-100)
	lot1.Spaces[80].Park(handicapCars[1]) // Space 81 = Row D
	fmt.Printf("✅ Parked %s in Row D (Space 81)\n", handicapCars[1].LicensePlate)
	
	// Park in Row A (spaces 1-25) - legitimate location
	lot1.Spaces[10].Park(handicapCars[2]) // Space 11 = Row A
	fmt.Printf("✅ Parked %s in Row A (Space 11)\n", handicapCars[2].LicensePlate)

	// UC15 Demo: Time-based parking queries
	fmt.Println("\n🔍 UC15: Time-based Parking Queries Demo")
	fmt.Println("==========================================")

	// Find cars parked in last 30 minutes
	recentVehicles, err := policeService.FindCarsParkedInLastMinutes(30)
	if err != nil {
		fmt.Printf("❌ Error finding recent vehicles: %v\n", err)
	} else {
		fmt.Printf("🚨 Found %d vehicles parked in last 30 minutes:\n", len(recentVehicles))
		for i, vehicle := range recentVehicles {
			fmt.Printf("   %d. %s (%s %s) - Lot: %s, Space: %s\n", 
				      i+1, vehicle.Car.LicensePlate, vehicle.Car.Color, 
				      vehicle.Car.Make, vehicle.LotID, vehicle.SpaceID)
		}
	}

	// Generate time-based investigation report
	fmt.Println("\n📋 UC15: Time-based Investigation Report:")
	timeReport := policeService.GenerateTimeBasedInvestigationReport(30)
	fmt.Print(timeReport)

	// UC15: Vehicle count by time window
	fmt.Println("\n📊 UC15: Vehicle Count Analysis:")
	counts := policeService.GetVehicleCountByTimeWindow(30)
	if errMsg, hasError := counts["error"]; hasError {
		fmt.Printf("❌ Count analysis error: %v\n", errMsg)
	} else {
		fmt.Printf("Time Window: %v\n", counts["timeWindow"])
		fmt.Printf("Total Vehicles: %v\n", counts["totalVehicles"])
		fmt.Printf("Last 15 minutes: %v\n", counts["last15Minutes"])
		fmt.Printf("Last 30 minutes: %v\n", counts["last30Minutes"])
	}

	// UC16 Demo: Location-based vehicle searches
	fmt.Println("\n🔍 UC16: Location-based Vehicle Searches Demo")
	fmt.Println("==============================================")

	// Find handicap cars in suspicious rows (B and D)
	suspiciousHandicapCars, err := policeService.FindHandicapCarsInRows([]string{"B", "D"})
	if err != nil {
		fmt.Printf("❌ Error finding handicap cars in rows B/D: %v\n", err)
	} else {
		fmt.Printf("🚨 Found %d handicap vehicles in suspicious rows B/D:\n", len(suspiciousHandicapCars))
		for i, vehicle := range suspiciousHandicapCars {
			fmt.Printf("   %d. %s (%s) - Lot: %s, Space: %s\n", 
				      i+1, vehicle.Car.LicensePlate, vehicle.Car.DriverName,
				      vehicle.LotID, vehicle.SpaceID)
		}
	}

	// Generate handicap fraud investigation report
	fmt.Println("\n📋 UC16: Handicap Fraud Investigation Report:")
	fraudReport := policeService.GenerateHandicapFraudInvestigationReport()
	fmt.Print(fraudReport)

	// UC16: Validate handicap permit fraud
	fmt.Println("\n🔒 UC16: Handicap Permit Fraud Validation:")
	validation := policeService.ValidateHandicapPermitFraud()
	fmt.Printf("Total Handicap Vehicles: %v\n", validation["totalHandicapVehicles"])
	fmt.Printf("Vehicles in Rows B/D: %v\n", validation["vehiclesInRowsB_D"])
	fmt.Printf("Fraud Risk Assessment: %v\n", validation["fraudRisk"])
	fmt.Printf("Investigation Required: %v\n", validation["investigationRequired"])

	// UC16: Location statistics
	fmt.Println("\n📊 UC16: Parking Location Statistics:")
	stats := policeService.GetLocationStatistics()
	rowCounts := stats["totalVehiclesByRow"].(map[string]int)
	handicapCounts := stats["handicapVehiclesByRow"].(map[string]int)

	fmt.Println("Vehicles by Row:")
	for row := 'A'; row <= 'D'; row++ {
		rowStr := string(row)
		fmt.Printf("   Row %s: %d total, %d handicap\n", 
			  rowStr, rowCounts[rowStr], handicapCounts[rowStr])
	}

	// Combined UC15-UC16 Analysis
	fmt.Println("\n🔍 Combined UC15-UC16 Security Analysis")
	fmt.Println("======================================")

	// Find small handicap cars in suspicious rows that were parked recently
	smallHandicapInBD, err := policeService.GetVehiclesByLocationCriteria(
		models.SmallVehicle, true, []string{"B", "D"})
	if err != nil {
		fmt.Printf("❌ Error in combined analysis: %v\n", err)
	} else {
		fmt.Printf("🚨 Small handicap vehicles in rows B/D: %d\n", len(smallHandicapInBD))
		for i, vehicle := range smallHandicapInBD {
			timeSinceParked := time.Since(vehicle.ParkedAt)
			fmt.Printf("   %d. %s - Parked %.0f minutes ago\n", 
				      i+1, vehicle.Car.LicensePlate, timeSinceParked.Minutes())
		}
	}

	fmt.Println("\n🎉 UC15 & UC16 Implementation Complete!")
	fmt.Println("======================================")
	fmt.Println("✅ Time-based parking queries operational")
	fmt.Println("✅ Location-based vehicle searches operational")
	fmt.Println("✅ Handicap permit fraud detection active")
	fmt.Println("✅ Row-based parking space mapping functional")
	fmt.Println("✅ Integrated police investigation system ready")
}
func createVehicle(plate, driver string, size models.VehicleSize, handicap bool) *models.Car {
	car := models.NewCar(plate, driver)
	car.SetVehicleSize(size)
	car.SetHandicapStatus(handicap)
	return car
}

func showDetailedAnalytics(service *services.ParkingService) {
	analytics := service.GetDetailedLotAnalytics()
	for lotID, data := range analytics {
		fmt.Printf("   %s: %d/%d spaces (%.1f%% utilized)\n",
			lotID, data["OccupiedSpaces"], data["TotalSpaces"], data["UtilizationRate"])
		fmt.Printf("      └─ Handicap: %d, Large: %d, Small: %d\n",
			data["HandicapVehicles"], data["LargeVehicles"], data["SmallVehicles"])
	}
}

    // UC11 Large Vehicle Management Demo
    fmt.Println("\n🚛 UC11 Large Vehicle Management:")
    
    // Test large vehicle strategy
    testLargeCar := models.NewCar("LARGE_TEST", "Large Vehicle Test")
    testLargeCar.SetVehicleSize(models.LargeVehicle)
    
    largeDecision, err := service.ParkLargeVehicle(testLargeCar, "ATT001")
    if err == nil {
        fmt.Printf("✅ Large vehicle parked in: %s\n", largeDecision.LotID)
    } else {
        fmt.Printf("❌ Large vehicle parking failed: %v\n", err)
    }
    
    // Show recommendations
    recommendations := service.GetLargeVehicleRecommendations()
    if errMsg, hasError := recommendations["error"]; hasError {
        fmt.Printf("❌ No recommendations: %v\n", errMsg)
    } else {
        fmt.Printf("📋 Best lot for large vehicles: %v\n", recommendations["recommendedLot"])
        fmt.Printf("   Available spaces: %v\n", recommendations["availableSpaces"])
        fmt.Printf("   Current large vehicles: %v\n", recommendations["currentLargeVehicles"])
        fmt.Printf("   Utilization rate: %.1f%%\n", recommendations["utilizationRate"])
        
        if alternatives, ok := recommendations["alternativeLots"].([]string); ok && len(alternatives) > 0 {
            fmt.Printf("   Alternative lots: %v\n", alternatives)
        }
    }
    
    // Show capacity validation
    fmt.Println("\n🔍 Large Vehicle Capacity Validation:")
    validation := service.ValidateLargeVehicleCapacity()
    for lotID, suitable := range validation {
        status := "❌ Over capacity"
        if suitable {
            status = "✅ Suitable"
        }
        fmt.Printf("   %s: %s\n", lotID, status)
    }
    
    // Show optimal placement analysis
    fmt.Println("\n📊 Optimal Large Vehicle Placement:")
    placement := service.GetOptimalLargeVehiclePlacement()
    if errMsg, hasError := placement["error"]; hasError {
        fmt.Printf("❌ No placement analysis: %v\n", errMsg)
    } else {
        fmt.Printf("   Optimal lot: %v\n", placement["optimalLot"])
        fmt.Printf("   Maneuvering efficiency: %.1f%%\n", placement["maneuveringEfficiency"])
        fmt.Printf("   Recommended: %v\n", placement["recommendedForLargeVehicles"])
    }

    // UC12-UC13 Police Investigation Demo
    fmt.Println("\n🚔 UC12-UC13 Police Investigation Features:")
    
    // Create police service
    policeService := services.NewPoliceService(service)
    
    // Add some test cars with colors and makes
    testCars := []*models.Car{
        func() *models.Car {
            car := models.NewCar("WHITE001", "John Smith")
            car.SetColor("White")
            car.SetMake("Honda")
            return car
        }(),
        func() *models.Car {
            car := models.NewCar("BLUE_TOY", "Jane Doe")
            car.SetColor("Blue")
            car.SetMake("Toyota")
            return car
        }(),
        func() *models.Car {
            car := models.NewCar("WHITE002", "Bob Wilson")
            car.SetColor("White")
            car.SetMake("Ford")
            return car
        }(),
    }
    
    // Park test cars
    for _, car := range testCars {
        _, err := service.ParkCarWithTicket(car)
        if err != nil {
            fmt.Printf("Warning: Could not park %s: %v\n", car.LicensePlate, err)
        }
    }
    
    // UC12: Find white cars for bomb threat investigation
    fmt.Println("\n🔍 UC12: Investigating bomb threat - Finding white cars...")
    whiteCars, err := policeService.FindWhiteCars()
    if err != nil {
        fmt.Printf("❌ Error finding white cars: %v\n", err)
    } else {
        fmt.Printf("✅ Found %d white vehicles:\n", len(whiteCars))
        for i, vehicle := range whiteCars {
            fmt.Printf("   %d. %s (%s %s) - Lot: %s, Space: %s\n", 
                      i+1, vehicle.Car.LicensePlate, vehicle.Car.Color, 
                      vehicle.Car.Make, vehicle.LotID, vehicle.SpaceID)
        }
    }
    
    // UC13: Find blue Toyota cars for robbery investigation
    fmt.Println("\n🔍 UC13: Investigating robbery - Finding blue Toyota cars...")
    blueToyotas, err := policeService.FindBlueToyotaCars()
    if err != nil {
        fmt.Printf("❌ Error finding blue Toyotas: %v\n", err)
    } else {
        fmt.Printf("✅ Found %d blue Toyota vehicles:\n", len(blueToyotas))
        for i, vehicle := range blueToyotas {
            fmt.Printf("   %d. %s (Driver: %s) - Lot: %s, Space: %s\n", 
                      i+1, vehicle.Car.LicensePlate, vehicle.Car.DriverName,
                      vehicle.LotID, vehicle.SpaceID)
            if vehicle.AttendantName != "" {
                fmt.Printf("      Attendant: %s (ID: %s)\n", vehicle.AttendantName, vehicle.AttendantID)
            }
        }
    }
    
    // Generate investigation reports
    if len(whiteCars) > 0 {
        fmt.Println("\n📋 White Cars Investigation Report:")
        report := policeService.GenerateInvestigationReport(whiteCars, "Bomb Threat Investigation")
        fmt.Print(report)
    }

    // UC12-UC13 Police Investigation Demo
    fmt.Println("\n🚔 UC12-UC13 Police Investigation Features:")
    
    // Create police service
    policeService := services.NewPoliceService(service)
    
    // Add some test cars with colors and makes
    testCars := []*models.Car{
        func() *models.Car {
            car := models.NewCar("WHITE001", "John Smith")
            car.SetColor("White")
            car.SetMake("Honda")
            return car
        }(),
        func() *models.Car {
            car := models.NewCar("BLUE_TOY", "Jane Doe")
            car.SetColor("Blue")
            car.SetMake("Toyota")
            return car
        }(),
        func() *models.Car {
            car := models.NewCar("WHITE002", "Bob Wilson")
            car.SetColor("White")
            car.SetMake("Ford")
            return car
        }(),
    }
    
    // Park test cars
    for _, car := range testCars {
        err := service.ParkCar(car)
        if err != nil {
            fmt.Printf("Warning: Could not park %s: %v\n", car.LicensePlate, err)
        }
    }
    
    // UC12: Find white cars for bomb threat investigation
    fmt.Println("\n🔍 UC12: Investigating bomb threat - Finding white cars...")
    whiteCars, err := policeService.FindWhiteCars()
    if err != nil {
        fmt.Printf("❌ Error finding white cars: %v\n", err)
    } else {
        fmt.Printf("✅ Found %d white vehicles:\n", len(whiteCars))
        for i, vehicle := range whiteCars {
            fmt.Printf("   %d. %s (%s %s) - Lot: %s, Space: %s\n", 
                      i+1, vehicle.Car.LicensePlate, vehicle.Car.Color, 
                      vehicle.Car.Make, vehicle.LotID, vehicle.SpaceID)
        }
    }
    
    // UC13: Find blue Toyota cars for robbery investigation
    fmt.Println("\n🔍 UC13: Investigating robbery - Finding blue Toyota cars...")
    blueToyotas, err := policeService.FindBlueToyotaCars()
    if err != nil {
        fmt.Printf("❌ Error finding blue Toyotas: %v\n", err)
    } else {
        fmt.Printf("✅ Found %d blue Toyota vehicles:\n", len(blueToyotas))
        for i, vehicle := range blueToyotas {
            fmt.Printf("   %d. %s (Driver: %s) - Lot: %s, Space: %s\n", 
                      i+1, vehicle.Car.LicensePlate, vehicle.Car.DriverName,
                      vehicle.LotID, vehicle.SpaceID)
            if vehicle.AttendantName != "" {
                fmt.Printf("      Attendant: %s (ID: %s)\n", vehicle.AttendantName, vehicle.AttendantID)
            }
        }
    }
    
    // Generate investigation reports
    if len(whiteCars) > 0 {
        fmt.Println("\n📋 White Cars Investigation Report:")
        report := policeService.GenerateInvestigationReport(whiteCars, "Bomb Threat Investigation")
        fmt.Print(report)
    }
    // UC12-UC13 Police Investigation Demo
    fmt.Println("\n🚔 UC12-UC13 Police Investigation Features:")
    
    // Create police service
    policeService := services.NewPoliceService(service)
    
    // Add some test cars with colors and makes
    testCars := []*models.Car{
        func() *models.Car {
            car := models.NewCar("WHITE001", "John Smith")
            car.SetColor("White")
            car.SetMake("Honda")
            return car
        }(),
        func() *models.Car {
            car := models.NewCar("BLUE_TOY", "Jane Doe")
            car.SetColor("Blue")
            car.SetMake("Toyota")
            return car
        }(),
        func() *models.Car {
            car := models.NewCar("WHITE002", "Bob Wilson")
            car.SetColor("White")
            car.SetMake("Ford")
            return car
        }(),
    }
    
    // Park test cars
    for _, car := range testCars {
        err := service.ParkCar(car)
        if err != nil {
            fmt.Printf("Warning: Could not park %s: %v\n", car.LicensePlate, err)
        }
    }
    
    // UC12: Find white cars for bomb threat investigation
    fmt.Println("\n🔍 UC12: Investigating bomb threat - Finding white cars...")
    whiteCars, err := policeService.FindWhiteCars()
    if err != nil {
        fmt.Printf("❌ Error finding white cars: %v\n", err)
    } else {
        fmt.Printf("✅ Found %d white vehicles:\n", len(whiteCars))
        for i, vehicle := range whiteCars {
            fmt.Printf("   %d. %s (%s %s) - Lot: %s, Space: %s\n", 
                      i+1, vehicle.Car.LicensePlate, vehicle.Car.Color, 
                      vehicle.Car.Make, vehicle.LotID, vehicle.SpaceID)
        }
    }
    
    // UC13: Find blue Toyota cars for robbery investigation
    fmt.Println("\n🔍 UC13: Investigating robbery - Finding blue Toyota cars...")
    blueToyotas, err := policeService.FindBlueToyotaCars()
    if err != nil {
        fmt.Printf("❌ Error finding blue Toyotas: %v\n", err)
    } else {
        fmt.Printf("✅ Found %d blue Toyota vehicles:\n", len(blueToyotas))
        for i, vehicle := range blueToyotas {
            fmt.Printf("   %d. %s (Driver: %s) - Lot: %s, Space: %s\n", 
                      i+1, vehicle.Car.LicensePlate, vehicle.Car.DriverName,
                      vehicle.LotID, vehicle.SpaceID)
            if vehicle.AttendantName != "" {
                fmt.Printf("      Attendant: %s (ID: %s)\n", vehicle.AttendantName, vehicle.AttendantID)
            }
        }
    }
    
    // Generate investigation reports
    if len(whiteCars) > 0 {
        fmt.Println("\n📋 White Cars Investigation Report:")
        report := policeService.GenerateInvestigationReport(whiteCars, "Bomb Threat Investigation")
        fmt.Print(report)
    }
    
    // UC13: Enhanced robbery investigation with attendant details
    fmt.Println("\n🚔 UC13: Blue Toyota Robbery Investigation")
    
    blueToyotaCount := policeService.GetBlueToyotaCount()
    fmt.Printf("📊 Total blue Toyota vehicles found: %d\n", blueToyotaCount)
    
    if blueToyotaCount > 0 {
        fmt.Println("\n📋 Generating comprehensive robbery investigation report...")
        robberyReport := policeService.GenerateRobberyInvestigationReport("Armed suspect, approximately 5'10\", wearing dark clothing")
        fmt.Print(robberyReport)
        
        fmt.Println("\n🔍 Evidence validation:")
        evidence := policeService.ValidateRobberyEvidence()
        fmt.Printf("   Case strength: %v\n", evidence["caseStrength"])
        fmt.Printf("   Evidence quality: %v\n", evidence["evidenceQuality"])
        fmt.Printf("   Attendant witnesses: %v\n", evidence["attendantWitnesses"])
    } else {
        fmt.Println("📋 No blue Toyota vehicles found for robbery investigation")
        fmt.Println("   Recommendation: Expand search to nearby areas")
    }
