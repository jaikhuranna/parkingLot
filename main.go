package main

import (
    "fmt"
    "parking-lot-system/models"
    "parking-lot-system/services"
    "parking-lot-system/interfaces"
)

func main() {
    fmt.Println("Welcome to Parking Lot System!")
    fmt.Println("UC9: Even Distribution Strategy Demo")
    fmt.Println("===================================")
    
    // Create multiple parking lots with different capacities
    lot1 := models.NewParkingLot("LOT1", 3)
    lot2 := models.NewParkingLot("LOT2", 5)
    lot3 := models.NewParkingLot("LOT3", 4)
    
    service := services.NewParkingService()
    service.AddLot(lot1)
    service.AddLot(lot2)
    service.AddLot(lot3)
    
    // Add attendant
    attendant := models.NewParkingAttendant("ATT001", "Alice Johnson", "LOT1")
    service.AddAttendant(attendant)
    
    // Add observers
    owner := interfaces.NewOwnerObserver("Sanjay")
    service.AddObserverToLot("LOT1", owner)
    service.AddObserverToLot("LOT2", owner)
    service.AddObserverToLot("LOT3", owner)
    
    fmt.Println("\n📊 Initial Lot Status:")
    showLotUtilization(service)
    
    // Demo UC9: Even distribution parking
    cars := []*models.Car{
        models.NewCar("ABC123", "John Doe"),
        models.NewCar("XYZ789", "Jane Smith"),
        models.NewCar("DEF456", "Bob Johnson"),
        models.NewCar("GHI789", "Alice Brown"),
        models.NewCar("JKL012", "Charlie Davis"),
    }
    
    fmt.Println("\n🎯 UC9 Demo: Even Distribution Parking...")
    
    for i, car := range cars {
        fmt.Printf("\n%d. Parking car %s (%s):\n", i+1, car.LicensePlate, car.DriverName)
        
        // Use even distribution strategy
        decision, err := service.ParkCarEvenDistribution(car, "ATT001")
        if err != nil {
            fmt.Printf("   ❌ Error: %v\n", err)
            continue
        }
        
        fmt.Printf("   ✅ Parked in %s (Space: %s)\n", decision.LotID, decision.SpaceID)
        fmt.Printf("   📋 Strategy: %s\n", decision.Reason)
        
        // Show current utilization after each parking
        fmt.Printf("   📊 Current utilization:\n")
        utilizations := service.GetLotUtilization()
        for _, util := range utilizations {
            fmt.Printf("      %s: %d/%d (%.1f%% full)\n", 
                      util.LotID, util.OccupiedSpaces, util.TotalSpaces, util.UtilizationRate)
        }
    }
    
    fmt.Println("\n📈 Final Lot Utilization:")
    showLotUtilization(service)
    
    // Try to park one more car when lots are getting full
    fmt.Println("\n🚗 Attempting to park additional car...")
    extraCar := models.NewCar("MNO345", "Extra Driver")
    decision, err := service.ParkCarEvenDistribution(extraCar, "ATT001")
    if err != nil {
        fmt.Printf("   ❌ Error: %v\n", err)
    } else {
        fmt.Printf("   ✅ Extra car parked in %s\n", decision.LotID)
    }
}

func showLotUtilization(service *services.ParkingService) {
    utilizations := service.GetLotUtilization()
    for _, util := range utilizations {
        fmt.Printf("   %s: %d/%d spaces (%.1f%% utilized)\n", 
                  util.LotID, util.OccupiedSpaces, util.TotalSpaces, util.UtilizationRate)
    }
}
