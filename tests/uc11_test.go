package tests

import (
	"github.com/stretchr/testify/assert"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"testing"
)

func TestUC11_LargeVehicleStrategy(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 2)
	lot2 := models.NewParkingLot("LOT2", 5)
	service.AddLot(lot1)
	service.AddLot(lot2)

	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	service.AddAttendant(attendant)

	largeCar := models.NewCar("TRUCK001", "Bob Wilson")
	largeCar.SetVehicleSize(models.LargeVehicle)

	// Act
	decision, err := service.ParkLargeVehicle(largeCar, "ATT001")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, decision)
	assert.Equal(t, "LOT2", decision.LotID) // Should choose lot with more spaces
	assert.Equal(t, "ATT001", decision.AttendantID)
	assert.Contains(t, decision.Reason, "Large Vehicle Strategy")
}

func TestUC11_NonLargeVehicleRejection(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	service.AddLot(lot)

	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	service.AddAttendant(attendant)

	smallCar := models.NewCar("SMALL001", "Jane Smith")
	smallCar.SetVehicleSize(models.SmallVehicle)

	// Act
	decision, err := service.ParkLargeVehicle(smallCar, "ATT001")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, decision)
	assert.Equal(t, "car is not classified as large vehicle", err.Error())
}

func TestUC11_BestLotForLargeVehicle(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 3)
	lot2 := models.NewParkingLot("LOT2", 7) // Most spaces
	lot3 := models.NewParkingLot("LOT3", 5)

	service.AddLot(lot1)
	service.AddLot(lot2)
	service.AddLot(lot3)

	// Act
	bestLot, maxSpaces, err := service.GetBestLotForLargeVehicle()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, bestLot)
	assert.Equal(t, "LOT2", bestLot.ID)
	assert.Equal(t, 7, maxSpaces)
}

func TestUC11_LargeVehicleRecommendations(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 3)
	lot2 := models.NewParkingLot("LOT2", 8)
	service.AddLot(lot1)
	service.AddLot(lot2)

	// Act
	recommendations := service.GetLargeVehicleRecommendations()

	// Assert
	assert.Equal(t, "LOT2", recommendations["recommendedLot"])
	assert.Equal(t, 8, recommendations["availableSpaces"])
	assert.Equal(t, 0, recommendations["currentLargeVehicles"])
	assert.Equal(t, 0.0, recommendations["utilizationRate"])
	assert.NotNil(t, recommendations["alternativeLots"])
}

func TestUC11_LargeVehicleCapacityValidation(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 10)
	lot2 := models.NewParkingLot("LOT2", 10)
	service.AddLot(lot1)
	service.AddLot(lot2)

	// Fill lot1 beyond 70% (8 out of 10 spaces)
	for i := 0; i < 8; i++ {
		car := models.NewCar("CAR"+string(rune('0'+i)), "Driver")
		lot1.ParkCar(car)
	}

	// Leave lot2 with low utilization (2 out of 10 spaces)
	for i := 0; i < 2; i++ {
		car := models.NewCar("CAR2_"+string(rune('0'+i)), "Driver2")
		lot2.ParkCar(car)
	}

	// Act
	validation := service.ValidateLargeVehicleCapacity()

	// Assert
	assert.False(t, validation["LOT1"]) // Over 70% capacity
	assert.True(t, validation["LOT2"])  // Under 70% capacity
}

func TestUC11_OptimalLargeVehiclePlacement(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 10)
	service.AddLot(lot)

	// Act
	placement := service.GetOptimalLargeVehiclePlacement()

	// Assert
	assert.Equal(t, "LOT1", placement["optimalLot"])
	assert.Equal(t, 10, placement["availableSpaces"])
	assert.Equal(t, 100.0, placement["maneuveringEfficiency"])
	assert.True(t, placement["recommendedForLargeVehicles"].(bool))
}

func TestUC11_LargeVehicleAnalytics(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 5)
	service.AddLot(lot)

	// Park different vehicle sizes
	largeCar1 := models.NewCar("L1", "Large 1")
	largeCar1.SetVehicleSize(models.LargeVehicle)

	largeCar2 := models.NewCar("L2", "Large 2")
	largeCar2.SetVehicleSize(models.LargeVehicle)

	smallCar := models.NewCar("S1", "Small 1")
	smallCar.SetVehicleSize(models.SmallVehicle)

	lot.ParkCar(largeCar1)
	lot.ParkCar(largeCar2)
	lot.ParkCar(smallCar)

	// Act
	largeCounts := service.GetLargeVehicleSpacesCount()
	analytics := service.GetDetailedLotAnalytics()

	// Assert
	assert.Equal(t, 2, largeCounts["LOT1"])
	assert.Equal(t, 2, analytics["LOT1"]["LargeVehicles"])
	assert.Equal(t, 1, analytics["LOT1"]["SmallVehicles"])
	assert.Equal(t, 3, analytics["LOT1"]["OccupiedSpaces"])
}

func TestUC11_NoLotsAvailableForLargeVehicles(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 1)
	lot1.ParkCar(models.NewCar("FULL1", "Driver1"))
	service.AddLot(lot1)

	// Act
	bestLot, maxSpaces, err := service.GetBestLotForLargeVehicle()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, bestLot)
	assert.Equal(t, 0, maxSpaces)
	assert.Equal(t, "no lots available for large vehicles", err.Error())
}

func TestUC11_LargeVehicleWithSmartStrategy(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 3)
	lot2 := models.NewParkingLot("LOT2", 7)
	service.AddLot(lot1)
	service.AddLot(lot2)

	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	service.AddAttendant(attendant)

	largeCar := models.NewCar("SMART_L1", "Smart Large")
	largeCar.SetVehicleSize(models.LargeVehicle)

	// Act
	decision, err := service.ParkCarSmart(largeCar, "ATT001")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, decision)
	assert.Equal(t, "LOT2", decision.LotID) // Should choose largest lot
	assert.Contains(t, decision.Reason, "Smart Parking Strategy")
}
