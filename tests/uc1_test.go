package tests

import (
	"github.com/stretchr/testify/assert"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"testing"
)

func TestUC1_DriverCanParkCar(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 100)
	service.AddLot(lot)

	car := models.NewCar("ABC123", "John Doe")

	// Act
	err := service.ParkCar(car)

	// Assert
	assert.NoError(t, err)
	assert.True(t, lot.Spaces[0].IsOccupied)
	assert.Equal(t, "ABC123", lot.Spaces[0].ParkedCar.LicensePlate)
}
