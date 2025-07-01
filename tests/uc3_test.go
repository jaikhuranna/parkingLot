package tests

import (
	"github.com/stretchr/testify/assert"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"testing"
)

// Mock observer for testing notifications
type MockOwnerObserver struct {
	NotifiedFull      bool
	NotifiedAvailable bool
	LastLotID         string
}

func NewMockOwnerObserver() *MockOwnerObserver {
	return &MockOwnerObserver{
		NotifiedFull:      false,
		NotifiedAvailable: false,
		LastLotID:         "",
	}
}

func (m *MockOwnerObserver) OnLotFull(lotID string) {
	m.NotifiedFull = true
	m.LastLotID = lotID
}

func (m *MockOwnerObserver) OnLotAvailable(lotID string) {
	m.NotifiedAvailable = true
	m.LastLotID = lotID
}

func (m *MockOwnerObserver) Reset() {
	m.NotifiedFull = false
	m.NotifiedAvailable = false
	m.LastLotID = ""
}

func TestUC3_OwnerNotifiedWhenLotBecomesFull(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 2) // Small capacity for easy testing
	service.AddLot(lot)

	mockOwner := NewMockOwnerObserver()
	service.AddObserverToLot("LOT1", mockOwner)

	car1 := models.NewCar("ABC123", "John Doe")
	car2 := models.NewCar("XYZ789", "Jane Smith")

	// Act - Park first car (should not notify)
	err1 := service.ParkCar(car1)
	assert.NoError(t, err1)

	// Assert - No notification yet
	assert.False(t, mockOwner.NotifiedFull)

	// Act - Park second car (lot becomes full)
	err2 := service.ParkCar(car2)
	assert.NoError(t, err2)

	// Assert - Owner should be notified
	assert.True(t, mockOwner.NotifiedFull)
	assert.Equal(t, "LOT1", mockOwner.LastLotID)
	assert.True(t, lot.IsFull())
}

func TestUC3_OwnerNotifiedWhenLotBecomesAvailable(t *testing.T) {
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

	// Reset mock to test availability notification
	mockOwner.Reset()

	// Act - Unpark one car (lot becomes available)
	unparkedCar, err := service.UnparkCar("ABC123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, unparkedCar)
	assert.True(t, mockOwner.NotifiedAvailable)
	assert.Equal(t, "LOT1", mockOwner.LastLotID)
	assert.False(t, lot.IsFull())
	assert.Equal(t, 1, lot.GetAvailableSpaces())
}

func TestUC3_NoDuplicateNotifications(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 1) // Single space lot
	service.AddLot(lot)

	mockOwner := NewMockOwnerObserver()
	service.AddObserverToLot("LOT1", mockOwner)

	car1 := models.NewCar("ABC123", "John Doe")
	car2 := models.NewCar("XYZ789", "Jane Smith")

	// Act - Park car to fill lot
	err1 := service.ParkCar(car1)
	assert.NoError(t, err1)
	assert.True(t, mockOwner.NotifiedFull)

	// Reset notification flag
	mockOwner.NotifiedFull = false

	// Act - Try to park another car (should fail, no new notification)
	err2 := service.ParkCar(car2)

	// Assert - No duplicate notification
	assert.Error(t, err2)
	assert.False(t, mockOwner.NotifiedFull)
}

func TestUC3_MultipleObserversNotified(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 1)
	service.AddLot(lot)

	mockOwner1 := NewMockOwnerObserver()
	mockOwner2 := NewMockOwnerObserver()
	service.AddObserverToLot("LOT1", mockOwner1)
	service.AddObserverToLot("LOT1", mockOwner2)

	car := models.NewCar("ABC123", "John Doe")

	// Act - Park car to fill lot
	err := service.ParkCar(car)

	// Assert - Both observers notified
	assert.NoError(t, err)
	assert.True(t, mockOwner1.NotifiedFull)
	assert.True(t, mockOwner2.NotifiedFull)
	assert.Equal(t, "LOT1", mockOwner1.LastLotID)
	assert.Equal(t, "LOT1", mockOwner2.LastLotID)
}

func TestUC3_ObserverCanBeRemoved(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 1)
	service.AddLot(lot)

	mockOwner := NewMockOwnerObserver()
	service.AddObserverToLot("LOT1", mockOwner)

	// Remove observer
	service.RemoveObserverFromLot("LOT1", mockOwner)

	car := models.NewCar("ABC123", "John Doe")

	// Act - Park car to fill lot
	err := service.ParkCar(car)

	// Assert - Observer not notified after removal
	assert.NoError(t, err)
	assert.False(t, mockOwner.NotifiedFull)
	assert.True(t, lot.IsFull())
}

func TestUC3_LotStatusMethods(t *testing.T) {
	// Arrange
	lot := models.NewParkingLot("LOT1", 3)
	car1 := models.NewCar("ABC123", "John Doe")
	car2 := models.NewCar("XYZ789", "Jane Smith")

	// Test initial state
	assert.False(t, lot.IsFull())
	assert.Equal(t, 3, lot.GetAvailableSpaces())
	assert.Equal(t, 0, lot.GetOccupiedSpaces())

	// Park one car
	err1 := lot.ParkCar(car1)
	assert.NoError(t, err1)
	assert.False(t, lot.IsFull())
	assert.Equal(t, 2, lot.GetAvailableSpaces())
	assert.Equal(t, 1, lot.GetOccupiedSpaces())

	// Park second car
	err2 := lot.ParkCar(car2)
	assert.NoError(t, err2)
	assert.False(t, lot.IsFull())
	assert.Equal(t, 1, lot.GetAvailableSpaces())
	assert.Equal(t, 2, lot.GetOccupiedSpaces())
}

func TestUC3_ServiceLevelLotStatus(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 1)
	lot2 := models.NewParkingLot("LOT2", 2)
	service.AddLot(lot1)
	service.AddLot(lot2)

	car1 := models.NewCar("ABC123", "John Doe")

	// Initially no lots are full
	assert.False(t, service.IsAnyLotFull())

	// Fill first lot
	err1 := service.ParkCar(car1) // Should go to LOT1
	assert.NoError(t, err1)
	assert.True(t, service.IsAnyLotFull())

	// Get specific lot status
	lotStatus, err := service.GetLotStatus("LOT1")
	assert.NoError(t, err)
	assert.True(t, lotStatus.IsFull())

	lotStatus2, err := service.GetLotStatus("LOT2")
	assert.NoError(t, err)
	assert.False(t, lotStatus2.IsFull())
}

func TestUC3_InvalidLotOperations(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	mockOwner := NewMockOwnerObserver()

	// Test adding observer to non-existent lot
	err1 := service.AddObserverToLot("NONEXISTENT", mockOwner)
	assert.Error(t, err1)
	assert.Equal(t, "lot not found", err1.Error())

	// Test removing observer from non-existent lot
	err2 := service.RemoveObserverFromLot("NONEXISTENT", mockOwner)
	assert.Error(t, err2)
	assert.Equal(t, "lot not found", err2.Error())

	// Test getting status of non-existent lot
	_, err3 := service.GetLotStatus("NONEXISTENT")
	assert.Error(t, err3)
	assert.Equal(t, "lot not found", err3.Error())
}
