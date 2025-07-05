package main

import (
	"fmt"
	"parking-lot-system/interfaces"
	"parking-lot-system/models"
	"parking-lot-system/services"
)

func main() {
	fmt.Println("Welcome to Parking Lot System!")
	fmt.Println("UC10-UC11: Advanced Parking Strategies Demo")
	fmt.Println("===========================================")

	// Create multiple parking lots with different capacities
	lot1 := models.NewParkingLot("LOT1", 4) // Smaller lot
	lot2 := models.NewParkingLot("LOT2", 6) // Medium lot
	lot3 := models.NewParkingLot("LOT3", 8) // Larger lot

	service := services.NewParkingService()
	service.AddLot(lot1)
	service.AddLot(lot2)
	service.AddLot(lot3)

	// Add attendants
	attendant := models.NewParkingAttendant("ATT001", "Alice Johnson", "LOT1")
	service.AddAttendant(attendant)

	// Add observers
	owner := interfaces.NewOwnerObserver("Sanjay")
	service.AddObserverToLot("LOT1", owner)
	service.AddObserverToLot("LOT2", owner)
	service.AddObserverToLot("LOT3", owner)

	fmt.Println("\n📊 Initial Lot Status:")
	showDetailedAnalytics(service)

	// Demo vehicles with different properties
	vehicles := []*models.Car{
		createVehicle("ABC123", "John Doe", models.MediumVehicle, false),       // Regular car
		createVehicle("HANDICAP1", "Jane Smith", models.SmallVehicle, true),    // Handicap small car
		createVehicle("LARGE001", "Bob Wilson", models.LargeVehicle, false),    // Large vehicle
		createVehicle("HANDICAP2", "Alice Brown", models.LargeVehicle, true),   // Handicap large car
		createVehicle("TRUCK001", "Charlie Davis", models.LargeVehicle, false), // Another large vehicle
		createVehicle("SMALL001", "Diana Lee", models.SmallVehicle, false),     // Small regular car
	}

	fmt.Println("\n🎯 UC10-UC11 Demo: Advanced Parking Strategies...")

	for i, car := range vehicles {
		fmt.Printf("\n%d. Parking %s (%s, %s%s):\n",
			i+1, car.LicensePlate, car.DriverName,
			car.GetVehicleSizeString(),
			func() string {
				if car.IsHandicap {
					return ", Handicap"
				} else {
					return ""
				}
			}())

		var decision *models.ParkingDecision
		var err error

		// Use smart parking strategy that handles both UC10 and UC11
		decision, err = service.ParkCarSmart(car, "ATT001")

		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
			continue
		}

		fmt.Printf("   ✅ Parked in %s (Space: %s)\n", decision.LotID, decision.SpaceID)
		fmt.Printf("   📋 Strategy: %s\n", decision.Reason)

		// Show current analytics
		fmt.Printf("   📊 Current lot analytics:\n")
		analytics := service.GetDetailedLotAnalytics()
		for lotID, data := range analytics {
			fmt.Printf("      %s: %d/%d spaces, H:%d L:%d S:%d (%.1f%% full)\n",
				lotID,
				data["OccupiedSpaces"], data["TotalSpaces"],
				data["HandicapVehicles"], data["LargeVehicles"], data["SmallVehicles"],
				data["UtilizationRate"])
		}
	}

	fmt.Println("\n📈 Final Detailed Analytics:")
	showDetailedAnalytics(service)

	// Demonstrate specific strategy usage
	fmt.Println("\n🔍 Testing Individual Strategies:")

	// Test handicap priority
	handicapCar := createVehicle("H_TEST", "Handicap Test", models.MediumVehicle, true)
	decision, err := service.ParkHandicapCar(handicapCar, "ATT001")
	if err == nil {
		fmt.Printf("✅ Handicap priority: %s → %s\n", handicapCar.LicensePlate, decision.LotID)
	} else {
		fmt.Printf("❌ Handicap priority failed: %v\n", err)
	}

	// Test large vehicle strategy
	largeCar := createVehicle("L_TEST", "Large Test", models.LargeVehicle, false)
	decision, err = service.ParkLargeVehicle(largeCar, "ATT001")
	if err == nil {
		fmt.Printf("✅ Large vehicle strategy: %s → %s\n", largeCar.LicensePlate, decision.LotID)
	} else {
		fmt.Printf("❌ Large vehicle strategy failed: %v\n", err)
	}
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
