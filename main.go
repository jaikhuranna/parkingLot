package main

import (
	"fmt"
	"parking-lot-system/interfaces"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"time"
)

func main() {
	fmt.Println("Welcome to Parking Lot System!")
	fmt.Println("UC8: Time Tracking and Billing Demo")
	fmt.Println("===================================")

	// Create parking lot
	lot := models.NewParkingLot("LOT1", 5)
	service := services.NewParkingService()
	service.AddLot(lot)

	// Add observers
	owner := interfaces.NewOwnerObserver("Sanjay")
	service.AddObserverToLot("LOT1", owner)

	// Demo UC8: Park car with ticket
	car := models.NewCar("ABC123", "John Doe")
	fmt.Printf("🚗 Parking car %s...\n", car.LicensePlate)

	ticket, err := service.ParkCarWithTicket(car)
	if err != nil {
		fmt.Printf("Error parking car: %v\n", err)
		return
	}

	fmt.Printf("✅ Car parked successfully!\n")
	fmt.Printf("   Ticket ID: %s\n", ticket.ID)
	fmt.Printf("   Parked at: %s\n", ticket.ParkedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Lot: %s, Space: %s\n", ticket.LotID, ticket.SpaceID)

	// Simulate some parking time
	fmt.Printf("\n⏰ Simulating parking duration...\n")
	time.Sleep(100 * time.Millisecond) // Simulate time passage

	// Demo billing on unpark
	fmt.Printf("\n💰 Unparking with billing...\n")
	unparkedCar, bill, err := service.UnparkCarWithBilling(car.LicensePlate)
	if err != nil {
		fmt.Printf("Error unparking: %v\n", err)
		return
	}

	fmt.Printf("✅ Car %s unparked successfully!\n", unparkedCar.LicensePlate)
	fmt.Printf("📄 Billing Information:\n")
	fmt.Printf("   Duration: %s\n", bill.Duration.String())
	fmt.Printf("   Rate: $%.2f/hour\n", bill.HourlyRate)
	fmt.Printf("   Amount: $%.2f\n", bill.TotalAmount)

	// Show detailed bill
	fmt.Printf("\n🧾 Detailed Bill:\n")
	fmt.Print(bill.PrintBill())

	// Demo parking history
	fmt.Printf("\n📊 Parking History for %s:\n", car.LicensePlate)
	history, err := service.GetParkingHistory(car.LicensePlate)
	if err != nil {
		fmt.Printf("Error getting history: %v\n", err)
	} else {
		for i, h := range history {
			info := h.GetTicketInfo()
			fmt.Printf("   %d. %s - %s (Duration: %v)\n",
				i+1, info["ParkedAt"], info["UnparkedAt"], info["Duration"])
		}
	}
}
