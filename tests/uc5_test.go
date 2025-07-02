package tests

import (
	"github.com/stretchr/testify/assert"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"testing"
)

func TestUC5_OwnerNotifiedWhenSpaceBecomesAvailable(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 2)
	service.AddLot(lot)

	mockOwner := NewMockOwnerObserver()
	service.AddObserverToLot("LOT1", mockOwner)

	// Fill the lot first
	car1 := models.NewCar("ABC123", "John Doe")
	car2 := models.NewCar("XYZ789", "Jane Smith")
	service.ParkCar(car1)
	service.ParkCar(car2)

	// Verify lot is full
	assert.True(t, lot.IsFull())
	assert.True(t, mockOwner.NotifiedFull)

	// Reset mock to test availability notification
	mockOwner.Reset()

	// Act - Unpark one car (space becomes available)
	unparkedCar, err := service.UnparkCar("ABC123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, unparkedCar)
	assert.True(t, mockOwner.NotifiedAvailable)
	assert.Equal(t, "LOT1", mockOwner.LastLotID)
	assert.False(t, lot.IsFull())
	assert.Equal(t, 1, lot.GetAvailableSpaces())
}

func TestUC5_CanParkNewCarAfterSpaceAvailable(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 2)
	service.AddLot(lot)

	// Fill the lot
	car1 := models.NewCar("ABC123", "John Doe")
	car2 := models.NewCar("XYZ789", "Jane Smith")
	service.ParkCar(car1)
	service.ParkCar(car2)

	// Verify lot is full
	car3 := models.NewCar("DEF456", "Bob Johnson")
	err := service.ParkCar(car3)
	assert.Error(t, err)

	// Act - Unpark one car
	_, err = service.UnparkCar("ABC123")
	assert.NoError(t, err)

	// Assert - New car can now park
	err = service.ParkCar(car3)
	assert.NoError(t, err)

	// Verify car is actually parked
	foundSpace, err := service.FindCar("DEF456")
	assert.NoError(t, err)
	assert.NotNil(t, foundSpace)
	assert.Equal(t, "DEF456", foundSpace.ParkedCar.LicensePlate)
}

func TestUC5_MultipleSpacesBecomingAvailable(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	service.AddLot(lot)

	mockOwner := NewMockOwnerObserver()
	service.AddObserverToLot("LOT1", mockOwner)

	// Fill the lot completely
	cars := []*models.Car{
		models.NewCar("ABC123", "John"),
		models.NewCar("XYZ789", "Jane"),
		models.NewCar("DEF456", "Bob"),
	}

	for _, car := range cars {
		service.ParkCar(car)
	}

	assert.True(t, lot.IsFull())
	assert.True(t, mockOwner.NotifiedFull)

	// Reset and unpark first car
	mockOwner.Reset()
	_, err := service.UnparkCar("ABC123")
	assert.NoError(t, err)

	// Should be notified when first space becomes available
	assert.True(t, mockOwner.NotifiedAvailable)
	assert.Equal(t, 1, lot.GetAvailableSpaces())

	// Reset and unpark second car
	mockOwner.Reset()
	_, err = service.UnparkCar("XYZ789")
	assert.NoError(t, err)

	// Should NOT be notified again (already available)
	assert.False(t, mockOwner.NotifiedAvailable)
	assert.Equal(t, 2, lot.GetAvailableSpaces())
}

func TestUC5_NoNotificationWhenAlreadyAvailable(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	service.AddLot(lot)

	mockOwner := NewMockOwnerObserver()
	service.AddObserverToLot("LOT1", mockOwner)

	// Park only one car (lot not full)
	car1 := models.NewCar("ABC123", "John Doe")
	service.ParkCar(car1)

	// Verify no full notification
	assert.False(t, mockOwner.NotifiedFull)
	assert.False(t, lot.IsFull())

	// Act - Unpark the car
	_, err := service.UnparkCar("ABC123")
	assert.NoError(t, err)

	// Assert - No availability notification (was never full)
	assert.False(t, mockOwner.NotifiedAvailable)
}
