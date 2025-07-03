package tests

import (
	"github.com/stretchr/testify/assert"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"testing"
)

func TestUC10_HandicapPriorityStrategy(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 2)
	lot2 := models.NewParkingLot("LOT2", 3)
	service.AddLot(lot1)
	service.AddLot(lot2)

	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	service.AddAttendant(attendant)

	handicapCar := models.NewCar("HANDICAP1", "John Doe")
	handicapCar.SetHandicapStatus(true)

	// Act
	decision, err := service.ParkHandicapCar(handicapCar, "ATT001")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, decision)
	assert.Equal(t, "ATT001", decision.AttendantID)
	assert.Contains(t, decision.Reason, "Handicap Priority Strategy")
}

func TestUC10_NonHandicapCarRejection(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	service.AddLot(lot)

	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	service.AddAttendant(attendant)

	regularCar := models.NewCar("REGULAR1", "Jane Smith")

	// Act
	decision, err := service.ParkHandicapCar(regularCar, "ATT001")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, decision)
	assert.Equal(t, "car is not registered as handicap vehicle", err.Error())
}

func TestUC10_VehicleSizeProperties(t *testing.T) {
	// Arrange
	car := models.NewCar("TEST123", "Test Driver")

	// Test default values - NOW CORRECTLY EXPECTS MediumVehicle
	assert.Equal(t, models.MediumVehicle, car.Size) // Should be 1 (MediumVehicle)
	assert.False(t, car.IsHandicap)
	assert.Equal(t, "Medium", car.GetVehicleSizeString()) // Should be "Medium"

	// Test setting properties
	car.SetVehicleSize(models.SmallVehicle)
	car.SetHandicapStatus(true)

	assert.Equal(t, models.SmallVehicle, car.Size)
	assert.True(t, car.IsHandicap)
	assert.Equal(t, "Small", car.GetVehicleSizeString())

	// Test large vehicle
	car.SetVehicleSize(models.LargeVehicle)
	assert.Equal(t, models.LargeVehicle, car.Size)
	assert.Equal(t, "Large", car.GetVehicleSizeString())
}

func TestUC10_HandicapCarWithSmartStrategy(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 2)
	lot2 := models.NewParkingLot("LOT2", 4)
	service.AddLot(lot1)
	service.AddLot(lot2)

	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	service.AddAttendant(attendant)

	handicapCar := models.NewCar("SMART_H1", "Smart Handicap")
	handicapCar.SetHandicapStatus(true)
	handicapCar.SetVehicleSize(models.SmallVehicle)

	// Act
	decision, err := service.ParkCarSmart(handicapCar, "ATT001")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, decision)
	assert.Contains(t, decision.Reason, "Smart Parking Strategy")
}

func TestUC10_HandicapAnalytics(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 5)
	service.AddLot(lot)

	handicapCar1 := models.NewCar("H1", "Handicap 1")
	handicapCar1.SetHandicapStatus(true)

	handicapCar2 := models.NewCar("H2", "Handicap 2")
	handicapCar2.SetHandicapStatus(true)

	regularCar := models.NewCar("R1", "Regular 1")

	lot.ParkCar(handicapCar1)
	lot.ParkCar(handicapCar2)
	lot.ParkCar(regularCar)

	// Act
	handicapCounts := service.GetHandicapSpacesCount()
	analytics := service.GetDetailedLotAnalytics()

	// Assert
	assert.Equal(t, 2, handicapCounts["LOT1"])
	assert.Equal(t, 2, analytics["LOT1"]["HandicapVehicles"])
	assert.Equal(t, 3, analytics["LOT1"]["OccupiedSpaces"])
}

func TestUC10_VehicleSizeEnumValues(t *testing.T) {
	// Test that enum values are correct
	assert.Equal(t, 0, int(models.SmallVehicle))
	assert.Equal(t, 1, int(models.MediumVehicle))
	assert.Equal(t, 2, int(models.LargeVehicle))
}
