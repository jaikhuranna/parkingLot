package main

import (
    "fmt"
    "parking-lot-system/models"
    "parking-lot-system/services"
)

func main() {
    fmt.Println("Welcome to Parking Lot System!")
    fmt.Println("UC2: Driver Can Unpark Car Demo")
    
    // Create parking lot
    lot := models.NewParkingLot("LOT1", 100)
    service := services.NewParkingService()
    service.AddLot(lot)
    
    // Park a car
    car := models.NewCar("ABC123", "John Doe")
    err := service.ParkCar(car)
    
    if err != nil {
        fmt.Printf("Error parking car: %v\n", err)
        return
    }
    
    fmt.Printf("✅ Car %s parked successfully\n", car.LicensePlate)
    
    // Unpark the car
    unparkedCar, err := service.UnparkCar("ABC123")
    
    if err != nil {
        fmt.Printf("Error unparking car: %v\n", err)
        return
    }
    
    fmt.Printf("✅ Car %s unparked successfully\n", unparkedCar.LicensePlate)
    fmt.Printf("Driver %s can now go home!\n", unparkedCar.DriverName)
}

