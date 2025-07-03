package main

import (
    "fmt"
    "parking-lot-system/models"
    "parking-lot-system/services"
    "parking-lot-system/interfaces"
)

func main() {
    fmt.Println("Welcome to Parking Lot System!")
    fmt.Println("UC10-UC11: Advanced Parking Strategies Demo")
    fmt.Println("===========================================")
    
    // Create multiple parking lots with different capacities
    lot1 := models.NewParkingLot("LOT1", 4)  // Smaller lot
    lot2 := models.NewParkingLot("LOT2", 6)  // Medium lot
    lot3 := models.NewParkingLot("LOT3", 8)  // Larger lot
    
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
        createVehicle("ABC123", "John Doe", models.MediumVehicle, false),      // Regular car
        createVehicle("HANDICAP1", "Jane Smith", models.SmallVehicle, true),   // Handicap small car
        createVehicle("LARGE001", "Bob Wilson", models.LargeVehicle, false),   // Large vehicle
        createVehicle("HANDICAP2", "Alice Brown", models.LargeVehicle, true),  // Handicap large car
        createVehicle("TRUCK001", "Charlie Davis", models.LargeVehicle, false), // Another large vehicle
        createVehicle("SMALL001", "Diana Lee", models.SmallVehicle, false),    // Small regular car
    }
    
    fmt.Println("\n🎯 UC10-UC11 Demo: Advanced Parking Strategies...")
    
    for i, car := range vehicles {
        fmt.Printf("\n%d. Parking %s (%s, %s%s):\n", 
                  i+1, car.LicensePlate, car.DriverName, 
                  car.GetVehicleSizeString(),
                  func() string { if car.IsHandicap { return ", Handicap" } else { return "" } }())
        
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
