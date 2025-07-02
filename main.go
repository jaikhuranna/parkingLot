package main

import (
    "fmt"
    "parking-lot-system/models"
    "parking-lot-system/services"
    "parking-lot-system/interfaces"
)

func main() {
    fmt.Println("Welcome to Parking Lot System!")
    fmt.Println("UC7: Car Finding System Demo")
    fmt.Println("============================")
    
    // Create parking lot
    lot := models.NewParkingLot("LOT1", 5)
    service := services.NewParkingService()
    service.AddLot(lot)
    
    // Add observers
    owner := interfaces.NewOwnerObserver("Sanjay")
    service.AddObserverToLot("LOT1", owner)
    
    // Park some cars
    cars := []*models.Car{
        models.NewCar("ABC123", "John Doe"),
        models.NewCar("XYZ789", "Jane Smith"),
        models.NewCar("DEF456", "Bob Johnson"),
    }
    
    for _, car := range cars {
        ticket, err := service.ParkCar(car)
        if err != nil {
            fmt.Printf("Error parking %s: %v\n", car.LicensePlate, err)
            continue
        }
        fmt.Printf("✅ Car %s parked successfully\n", car.LicensePlate)
    }
    
    fmt.Println("\n🎯 UC7 Demo: Finding parked cars...")
    
    // Demonstrate car finding
    for _, car := range cars {
        fmt.Printf("\n🔍 Finding car %s:\n", car.LicensePlate)
        
        // Basic car finding
        space, err := service.FindCar(car.LicensePlate)
        if err != nil {
            fmt.Printf("❌ Error: %v\n", err)
            continue
        }
        
        fmt.Printf("   Found in space: %s\n", space.ID)
        
        // Detailed location finding
        location, err := service.FindCarWithLocation(car.LicensePlate)
        if err != nil {
            fmt.Printf("❌ Error getting location: %v\n", err)
            continue
        }
        
        info := location.GetLocationInfo()
        fmt.Printf("   Detailed info: Lot %v, Space %v, Row %v, Position %v\n", 
                   info["LotID"], info["SpaceID"], info["Row"], info["Position"])
        
        // Provide directions
        directions, err := service.ProvideDirectionsToDriver(car.LicensePlate)
        if err != nil {
            fmt.Printf("❌ Error getting directions: %v\n", err)
            continue
        }
        
        fmt.Println("   Directions:")
        fmt.Println("  ", directions)
    }
}
