package main

import (
	"fmt"
	"parking-lot-system/interfaces"
	"parking-lot-system/models"
	"parking-lot-system/services"
)

func main() {
	fmt.Println("Welcome to Parking Lot System!")
	fmt.Println("UC5: Space Available Notification Demo")
	fmt.Println("=====================================")

	// Create parking lot with small capacity for demo
	lot := models.NewParkingLot("LOT1", 2)
	service := services.NewParkingService()
	service.AddLot(lot)

	// Create security staff
	securityStaff := models.NewSecurityStaff("SEC001", "Officer Johnson", "Traffic Control")
	service.AddSecurityStaff(securityStaff)
	service.AssignSecurityToLot("SEC001", "LOT1")

	// Add observers
	owner := interfaces.NewOwnerObserver("Sanjay")
	security := interfaces.NewSecurityObserver("Officer Johnson", "SEC001")

	service.AddObserverToLot("LOT1", owner)
	service.AddObserverToLot("LOT1", security)

	fmt.Printf("Initial lot status: %d/%d spaces available\n",
		lot.GetAvailableSpaces(), lot.Capacity)

	// Fill the lot completely
	car1 := models.NewCar("ABC123", "John Doe")
	car2 := models.NewCar("XYZ789", "Jane Smith")

	service.ParkCar(car1)
	fmt.Printf("✅ Car %s parked. Available spaces: %d/%d\n",
		car1.LicensePlate, lot.GetAvailableSpaces(), lot.Capacity)

	service.ParkCar(car2)
	fmt.Printf("✅ Car %s parked. Available spaces: %d/%d\n",
		car2.LicensePlate, lot.GetAvailableSpaces(), lot.Capacity)

	// Demonstrate UC5: Space becomes available
	fmt.Println("\n🎯 UC5 Demo: Unparking car to make space available...")
	unparkedCar, err := service.UnparkCar("ABC123")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Car %s unparked. Available spaces: %d/%d\n",
		unparkedCar.LicensePlate, lot.GetAvailableSpaces(), lot.Capacity)

	// Show that new cars can now park
	car3 := models.NewCar("DEF456", "Bob Johnson")
	err = service.ParkCar(car3)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✅ New car %s successfully parked in available space!\n", car3.LicensePlate)
	}
}
