package main

import (
	"fmt"
	"parking-lot-system/interfaces"
	"parking-lot-system/models"
	"parking-lot-system/services"
)

func main() {
	fmt.Println("Welcome to Parking Lot System!")
	fmt.Println("UC3: Owner Notification System Demo")
	fmt.Println("=====================================")

	// Create parking lot with small capacity for demo
	lot := models.NewParkingLot("LOT1", 2)
	service := services.NewParkingService()
	service.AddLot(lot)

	// Add owner observer
	owner := interfaces.NewOwnerObserver("Sanjay")
	service.AddObserverToLot("LOT1", owner)

	fmt.Printf("Initial lot status: %d/%d spaces available\n",
		lot.GetAvailableSpaces(), lot.Capacity)

	// Park first car
	car1 := models.NewCar("ABC123", "John Doe")
	err := service.ParkCar(car1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Car %s parked. Available spaces: %d/%d\n",
		car1.LicensePlate, lot.GetAvailableSpaces(), lot.Capacity)

	// Park second car - lot becomes full
	car2 := models.NewCar("XYZ789", "Jane Smith")
	err = service.ParkCar(car2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Car %s parked. Available spaces: %d/%d\n",
		car2.LicensePlate, lot.GetAvailableSpaces(), lot.Capacity)

	// Try to park third car - should fail
	car3 := models.NewCar("DEF456", "Bob Johnson")
	err = service.ParkCar(car3)
	if err != nil {
		fmt.Printf("❌ Could not park car %s: %v\n", car3.LicensePlate, err)
	}

	// Unpark a car - lot becomes available
	fmt.Println("\nUnparking car...")
	unparkedCar, err := service.UnparkCar("ABC123")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Car %s unparked. Available spaces: %d/%d\n",
		unparkedCar.LicensePlate, lot.GetAvailableSpaces(), lot.Capacity)
}
