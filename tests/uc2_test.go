package tests

import (
    "testing"
    "parking-lot-system/models"
    "parking-lot-system/services"
    "github.com/stretchr/testify/assert"
)

func TestUC2_DriverCanUnparkCar(t *testing.T) {
    // Arrange
    service := services.NewParkingService()
    lot := models.NewParkingLot("LOT1", 100)
    service.AddLot(lot)
    
    car := models.NewCar("ABC123", "John Doe")
    
    // Park the car first
    err := service.ParkCar(car)
    assert.NoError(t, err)
    
    // Act - Unpark the car
    unparkedCar, err := service.UnparkCar("ABC123")
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, unparkedCar)
    assert.Equal(t, "ABC123", unparkedCar.LicensePlate)
    assert.Equal(t, "John Doe", unparkedCar.DriverName)
}

func TestUnparkCar_CarNotFound(t *testing.T) {
    // Arrange
    service := services.NewParkingService()
    lot := models.NewParkingLot("LOT1", 100)
    service.AddLot(lot)
    
    // Act - Try to unpark non-existent car
    unparkedCar, err := service.UnparkCar("NONEXISTENT")
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, unparkedCar)
    assert.Equal(t, "car not found in any parking lot", err.Error())
}

func TestUnparkCar_EmptyLicensePlate(t *testing.T) {
    // Arrange
    service := services.NewParkingService()
    lot := models.NewParkingLot("LOT1", 100)
    service.AddLot(lot)
    
    // Act - Try to unpark with empty license plate
    unparkedCar, err := service.UnparkCar("")
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, unparkedCar)
    assert.Equal(t, "license plate cannot be empty", err.Error())
}

func TestUnparkCar_SpaceBecomesAvailable(t *testing.T) {
    // Arrange
    service := services.NewParkingService()
    lot := models.NewParkingLot("LOT1", 2)
    service.AddLot(lot)
    
    car1 := models.NewCar("ABC123", "John Doe")
    car2 := models.NewCar("XYZ789", "Jane Smith")
    
    // Park both cars
    err1 := service.ParkCar(car1)
    err2 := service.ParkCar(car2)
    assert.NoError(t, err1)
    assert.NoError(t, err2)
    
    // Verify lot is full
    car3 := models.NewCar("DEF456", "Bob Johnson")
    err3 := service.ParkCar(car3)
    assert.Error(t, err3)
    
    // Act - Unpark one car
    unparkedCar, err := service.UnparkCar("ABC123")
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, unparkedCar)
    
    // Verify space is now available
    err4 := service.ParkCar(car3)
    assert.NoError(t, err4)
}
