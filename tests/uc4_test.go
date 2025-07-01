package tests

import (
	"github.com/stretchr/testify/assert"
	"parking-lot-system/interfaces"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"testing"
)

// Mock security observer for testing
type MockSecurityObserver struct {
	NotifiedFull      bool
	NotifiedAvailable bool
	LastLotID         string
}

func NewMockSecurityObserver() *MockSecurityObserver {
	return &MockSecurityObserver{
		NotifiedFull:      false,
		NotifiedAvailable: false,
		LastLotID:         "",
	}
}

func (m *MockSecurityObserver) OnLotFull(lotID string) {
	m.NotifiedFull = true
	m.LastLotID = lotID
}

func (m *MockSecurityObserver) OnLotAvailable(lotID string) {
	m.NotifiedAvailable = true
	m.LastLotID = lotID
}

func (m *MockSecurityObserver) Reset() {
	m.NotifiedFull = false
	m.NotifiedAvailable = false
	m.LastLotID = ""
}

func TestUC4_SecurityNotifiedWhenLotBecomesFull(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 2)
	service.AddLot(lot)

	mockSecurity := NewMockSecurityObserver()
	service.AddObserverToLot("LOT1", mockSecurity)

	car1 := models.NewCar("ABC123", "John Doe")
	car2 := models.NewCar("XYZ789", "Jane Smith")

	// Act - Park first car (should not notify)
	err1 := service.ParkCar(car1)
	assert.NoError(t, err1)
	assert.False(t, mockSecurity.NotifiedFull)

	// Act - Park second car (lot becomes full)
	err2 := service.ParkCar(car2)
	assert.NoError(t, err2)

	// Assert - Security should be notified
	assert.True(t, mockSecurity.NotifiedFull)
	assert.Equal(t, "LOT1", mockSecurity.LastLotID)
}

func TestUC4_SecurityStaffManagement(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 2)
	service.AddLot(lot)

	staff := models.NewSecurityStaff("SEC001", "Officer Johnson", "Traffic Control")
	service.AddSecurityStaff(staff)

	// Act
	err := service.AssignSecurityToLot("SEC001", "LOT1")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "LOT1", staff.AssignedLot)

	// Test staff info
	info := staff.GetInfo()
	assert.Equal(t, "SEC001", info["ID"])
	assert.Equal(t, "Officer Johnson", info["Name"])
	assert.Equal(t, "LOT1", info["AssignedLot"])
	assert.True(t, info["IsActive"].(bool))
}

func TestUC4_SecurityObserverWithRealImplementation(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 1)
	service.AddLot(lot)

	security := interfaces.NewSecurityObserver("Officer Smith", "SEC002")
	service.AddObserverToLot("LOT1", security)

	car := models.NewCar("ABC123", "John Doe")

	// Act - This should trigger notification (not tested for output, just no errors)
	err := service.ParkCar(car)

	// Assert
	assert.NoError(t, err)
	assert.True(t, lot.IsFull())
}

func TestUC4_InvalidSecurityOperations(t *testing.T) {
	// Arrange
	service := services.NewParkingService()

	// Test assigning non-existent staff
	err1 := service.AssignSecurityToLot("NONEXISTENT", "LOT1")
	assert.Error(t, err1)
	assert.Equal(t, "security staff not found", err1.Error())

	// Test assigning to non-existent lot
	staff := models.NewSecurityStaff("SEC001", "Officer", "Position")
	service.AddSecurityStaff(staff)

	err2 := service.AssignSecurityToLot("SEC001", "NONEXISTENT")
	assert.Error(t, err2)
	assert.Equal(t, "lot not found", err2.Error())
}
