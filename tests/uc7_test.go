package tests

import (
	"github.com/stretchr/testify/assert"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"strings"
	"testing"
)

func TestUC7_DriverCanFindTheirCar(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	service.AddLot(lot)

	car := models.NewCar("ABC123", "John Doe")

	// Park the car
	err := service.ParkCar(car)
	assert.NoError(t, err)

	// Act - Find the car
	space, err := service.FindCar("ABC123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, space)
	assert.Equal(t, "ABC123", space.ParkedCar.LicensePlate)
	assert.Equal(t, "John Doe", space.ParkedCar.DriverName)
}

func TestUC7_FindCarWithDetailedLocation(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	service.AddLot(lot)

	car := models.NewCar("XYZ789", "Jane Smith")

	// Park the car
	err := service.ParkCar(car)
	assert.NoError(t, err)

	// Act - Find car with location details
	location, err := service.FindCarWithLocation("XYZ789")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, location)
	assert.Equal(t, "XYZ789", location.Car.LicensePlate)
	assert.Equal(t, "LOT1", location.LotID)
	assert.NotEmpty(t, location.SpaceID)
	assert.False(t, location.ParkedAt.IsZero())
}

func TestUC7_ProvideDirectionsToDriver(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 5)
	service.AddLot(lot)

	car := models.NewCar("DEF456", "Bob Johnson")

	// Park the car
	err := service.ParkCar(car)
	assert.NoError(t, err)

	// Act - Get directions
	directions, err := service.ProvideDirectionsToDriver("DEF456")

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, directions)
	assert.Contains(t, directions, "DEF456")
	assert.Contains(t, directions, "LOT1")
	assert.Contains(t, directions, "Your car")
	assert.Contains(t, directions, "located in")
}

func TestUC7_CarNotFoundScenarios(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	service.AddLot(lot)

	// Test finding non-existent car
	space, err := service.FindCar("NONEXISTENT")
	assert.Error(t, err)
	assert.Nil(t, space)
	assert.Equal(t, "car not found", err.Error())

	// Test directions for non-existent car
	directions, err := service.ProvideDirectionsToDriver("NONEXISTENT")
	assert.Error(t, err)
	assert.Empty(t, directions)
}

func TestUC7_EmptyLicensePlateValidation(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	service.AddLot(lot)

	// Test empty license plate for FindCarWithLocation
	location, err := service.FindCarWithLocation("")
	assert.Error(t, err)
	assert.Nil(t, location)
	assert.Equal(t, "license plate cannot be empty", err.Error())

	// Test empty license plate for directions
	directions, err := service.ProvideDirectionsToDriver("")
	assert.Error(t, err)
	assert.Empty(t, directions)
}

func TestUC7_MultiLotCarFinding(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 2)
	lot2 := models.NewParkingLot("LOT2", 2)
	service.AddLot(lot1)
	service.AddLot(lot2)

	// Fill first lot
	car1 := models.NewCar("ABC123", "John")
	car2 := models.NewCar("XYZ789", "Jane")
	service.ParkCar(car1)
	service.ParkCar(car2)

	// Park car in second lot
	car3 := models.NewCar("DEF456", "Bob")
	err := service.ParkCar(car3)
	assert.NoError(t, err)

	// Act - Find car in second lot
	location, err := service.FindCarWithLocation("DEF456")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, location)
	assert.Equal(t, "LOT2", location.LotID)
	assert.Equal(t, "DEF456", location.Car.LicensePlate)
}

func TestUC7_CarLocationModel(t *testing.T) {
	// Arrange
	car := models.NewCar("ABC123", "John Doe")
	location := models.NewCarLocation(car, "LOT1", "A1", "A", 1, "ATT001")

	// Act
	info := location.GetLocationInfo()

	// Assert
	assert.Equal(t, "ABC123", info["LicensePlate"])
	assert.Equal(t, "LOT1", info["LotID"])
	assert.Equal(t, "A1", info["SpaceID"])
	assert.Equal(t, "A", info["Row"])
	assert.Equal(t, 1, info["Position"])
	assert.Equal(t, "ATT001", info["AttendantID"])
	assert.NotNil(t, info["ParkedAt"])
}

func TestUC7_DirectionsContainCorrectInformation(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 2)
	service.AddLot(lot)

	car := models.NewCar("TEST123", "Test Driver")
	service.ParkCar(car)

	// Act
	directions, err := service.ProvideDirectionsToDriver("TEST123")

	// Assert
	assert.NoError(t, err)

	// Check that directions contain expected information
	expectedPhrases := []string{
		"Your car TEST123",
		"Lot: LOT1",
		"Space:",
		"Row:",
		"Position:",
		"Parked at:",
	}

	for _, phrase := range expectedPhrases {
		assert.True(t, strings.Contains(directions, phrase),
			"Directions should contain: %s", phrase)
	}
}
