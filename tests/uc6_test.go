package tests

import (
	"github.com/stretchr/testify/assert"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"testing"
)

func TestUC6_AttendantCanParkCars(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	service.AddLot(lot)

	attendant := models.NewParkingAttendant("ATT001", "Alice Johnson", "LOT1")
	service.AddAttendant(attendant)

	car := models.NewCar("ABC123", "John Doe")

	// Act
	decision, err := service.ParkCarWithAttendant(car, "ATT001")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, decision)
	assert.Equal(t, "ATT001", decision.AttendantID)
	assert.Equal(t, "LOT1", decision.LotID)
	assert.Equal(t, "First available space strategy", decision.Reason)

	// Verify car is actually parked
	space, err := service.FindCar("ABC123")
	assert.NoError(t, err)
	assert.NotNil(t, space)
	assert.Equal(t, "ABC123", space.ParkedCar.LicensePlate)
}

func TestUC6_AttendantMakesParkingDecision(t *testing.T) {
	// Arrange
	lots := []*models.ParkingLot{
		models.NewParkingLot("LOT1", 2),
		models.NewParkingLot("LOT2", 2),
	}

	// Fill first lot
	car1 := models.NewCar("ABC123", "John")
	car2 := models.NewCar("XYZ789", "Jane")
	lots[0].ParkCar(car1)
	lots[0].ParkCar(car2)

	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	car3 := models.NewCar("DEF456", "Bob")

	// Act
	decision, err := attendant.MakeParkingDecision(lots, car3)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, decision)
	assert.Equal(t, "LOT2", decision.LotID) // Should choose available lot
	assert.Equal(t, "ATT001", decision.AttendantID)
}

func TestUC6_AttendantInfo(t *testing.T) {
	// Arrange
	attendant := models.NewParkingAttendant("ATT001", "Alice Johnson", "LOT1")

	// Act
	info := attendant.GetInfo()

	// Assert
	assert.Equal(t, "ATT001", info["ID"])
	assert.Equal(t, "Alice Johnson", info["Name"])
	assert.Equal(t, "LOT1", info["LotID"])
	assert.True(t, info["IsActive"].(bool))

	// Test deactivation
	attendant.SetActive(false)
	info = attendant.GetInfo()
	assert.False(t, info["IsActive"].(bool))
}

func TestUC6_InactiveAttendantCannotMakeDecisions(t *testing.T) {
	// Arrange
	lots := []*models.ParkingLot{models.NewParkingLot("LOT1", 2)}
	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	attendant.SetActive(false)

	car := models.NewCar("ABC123", "John")

	// Act
	decision, err := attendant.MakeParkingDecision(lots, car)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, decision)
	assert.Equal(t, "attendant is not active", err.Error())
}

func TestUC6_NoAvailableSpacesForAttendant(t *testing.T) {
	// Arrange
	lots := []*models.ParkingLot{models.NewParkingLot("LOT1", 1)}

	// Fill the lot
	car1 := models.NewCar("ABC123", "John")
	lots[0].ParkCar(car1)

	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	car2 := models.NewCar("XYZ789", "Jane")

	// Act
	decision, err := attendant.MakeParkingDecision(lots, car2)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, decision)
	assert.Equal(t, "no available parking spaces", err.Error())
}

func TestUC6_ServiceLevelAttendantOperations(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")

	// Test adding attendant
	service.AddAttendant(attendant)
	attendants := service.GetAttendants()
	assert.Len(t, attendants, 1)
	assert.Equal(t, "ATT001", attendants[0].ID)

	// Test finding attendant
	found := service.FindAttendantByID("ATT001")
	assert.NotNil(t, found)
	assert.Equal(t, "Alice", found.Name)

	// Test finding non-existent attendant
	notFound := service.FindAttendantByID("NONEXISTENT")
	assert.Nil(t, notFound)
}
