package tests

import (
    "testing"
    "parking-lot-system/models"
    "parking-lot-system/services"
    "github.com/stretchr/testify/assert"
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

