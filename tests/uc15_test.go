package tests

import (
	"testing"
	"time"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"github.com/stretchr/testify/assert"
)

func TestUC15_FindCarsParkedInLastMinutes(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 5)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	// Create and park cars with different timing
	oldCar := models.NewCar("OLD001", "Old Driver")
	recentCar := models.NewCar("RECENT001", "Recent Driver")

	// Park old car first
	lot.ParkCar(oldCar)
	
	// Simulate time passage
	time.Sleep(10 * time.Millisecond)
	
	// Park recent car
	lot.ParkCar(recentCar)

	// Act - Find cars parked in last 1 minute
	recentCars, err := policeService.FindCarsParkedInLastMinutes(1)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, recentCars, 2) // Both cars should be within 1 minute
	
	// Verify recent car is included
	found := false
	for _, vehicle := range recentCars {
		if vehicle.Car.LicensePlate == "RECENT001" {
			found = true
			assert.Equal(t, "LOT1", vehicle.LotID)
			break
		}
	}
	assert.True(t, found)
}

func TestUC15_GetRecentParkingActivity(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	car := models.NewCar("ACTIVITY001", "Activity Driver")
	lot.ParkCar(car)

	// Act
	activity, err := policeService.GetRecentParkingActivity(5 * time.Minute)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, activity, 1)
	assert.Equal(t, "ACTIVITY001", activity[0].Car.LicensePlate)
}

func TestUC15_GenerateTimeBasedInvestigationReport(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	suspiciousCar := models.NewCar("BOMB001", "Suspicious Driver")
	suspiciousCar.SetColor("Red")
	suspiciousCar.SetMake("Van")
	lot.ParkCar(suspiciousCar)

	// Act
	report := policeService.GenerateTimeBasedInvestigationReport(30)

	// Assert
	assert.Contains(t, report, "TIME-BASED INVESTIGATION REPORT")
	assert.Contains(t, report, "Bomb Threat Investigation")
	assert.Contains(t, report, "Last 30 minutes")
	assert.Contains(t, report, "BOMB001")
	assert.Contains(t, report, "Suspicious Driver")
	assert.Contains(t, report, "INVESTIGATION RECOMMENDATIONS")
}

func TestUC15_GetVehicleCountByTimeWindow(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 5)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	// Park multiple cars
	for i := 0; i < 3; i++ {
		car := models.NewCar("COUNT00"+string(rune('1'+i)), "Driver"+string(rune('1'+i)))
		lot.ParkCar(car)
		time.Sleep(1 * time.Millisecond) // Small delay
	}

	// Act
	counts := policeService.GetVehicleCountByTimeWindow(30)

	// Assert
	assert.Equal(t, 3, counts["totalVehicles"])
	assert.Equal(t, "30 minutes", counts["timeWindow"])
	assert.NotNil(t, counts["timestamp"])
}

func TestUC15_NoRecentVehicles(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	// Act - Search for cars in last 0 minutes (should find none)
	recentCars, err := policeService.FindCarsParkedInLastMinutes(0)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, recentCars, 0)
}
