package tests

import (
	"github.com/stretchr/testify/assert"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"testing"
)

func TestUC9_EvenDistributionStrategy(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 2)
	lot2 := models.NewParkingLot("LOT2", 3)
	service.AddLot(lot1)
	service.AddLot(lot2)

	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	service.AddAttendant(attendant)

	car := models.NewCar("ABC123", "John Doe")

	// Act - Use even distribution strategy
	decision, err := service.ParkCarEvenDistribution(car, "ATT001")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, decision)
	assert.Equal(t, "LOT2", decision.LotID) // Should choose lot with more spaces
	assert.Contains(t, decision.Reason, "Even Distribution Strategy")
}

func TestUC9_StrategyInterface(t *testing.T) {
	// Arrange
	strategy := models.NewEvenDistributionStrategy()
	lots := []*models.ParkingLot{
		models.NewParkingLot("LOT1", 2),
		models.NewParkingLot("LOT2", 5),
	}
	car := models.NewCar("XYZ789", "Jane Smith")

	// Act
	selectedLot, err := strategy.FindParkingLot(lots, car)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, selectedLot)
	assert.Equal(t, "LOT2", selectedLot.ID) // Should select lot with more spaces
	assert.Equal(t, "Even Distribution Strategy", strategy.GetStrategyName())
}

func TestUC9_NoAvailableLots(t *testing.T) {
	// Arrange
	strategy := models.NewEvenDistributionStrategy()

	// Create full lots
	lot1 := models.NewParkingLot("LOT1", 1)
	lot1.ParkCar(models.NewCar("FULL1", "Driver1"))

	lot2 := models.NewParkingLot("LOT2", 1)
	lot2.ParkCar(models.NewCar("FULL2", "Driver2"))

	lots := []*models.ParkingLot{lot1, lot2}
	car := models.NewCar("TEST123", "Test Driver")

	// Act
	selectedLot, err := strategy.FindParkingLot(lots, car)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, selectedLot)
	assert.Equal(t, "no available parking spaces in any lot", err.Error())
}

func TestUC9_EmptyLotsArray(t *testing.T) {
	// Arrange
	strategy := models.NewEvenDistributionStrategy()
	lots := []*models.ParkingLot{}
	car := models.NewCar("TEST123", "Test Driver")

	// Act
	selectedLot, err := strategy.FindParkingLot(lots, car)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, selectedLot)
	assert.Equal(t, "no parking lots available", err.Error())
}

func TestUC9_LotUtilizationCalculation(t *testing.T) {
	// Arrange
	lot := models.NewParkingLot("LOT1", 4)

	// Test empty lot
	util := models.CalculateLotUtilization(lot)
	assert.Equal(t, "LOT1", util.LotID)
	assert.Equal(t, 4, util.TotalSpaces)
	assert.Equal(t, 0, util.OccupiedSpaces)
	assert.Equal(t, 4, util.AvailableSpaces)
	assert.Equal(t, 0.0, util.UtilizationRate)

	// Park 2 cars
	lot.ParkCar(models.NewCar("CAR1", "Driver1"))
	lot.ParkCar(models.NewCar("CAR2", "Driver2"))

	// Test half-full lot
	util = models.CalculateLotUtilization(lot)
	assert.Equal(t, 4, util.TotalSpaces)
	assert.Equal(t, 2, util.OccupiedSpaces)
	assert.Equal(t, 2, util.AvailableSpaces)
	assert.Equal(t, 50.0, util.UtilizationRate)
}

func TestUC9_ServiceLevelUtilization(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 2)
	lot2 := models.NewParkingLot("LOT2", 3)
	service.AddLot(lot1)
	service.AddLot(lot2)

	// Park one car in lot1
	lot1.ParkCar(models.NewCar("CAR1", "Driver1"))

	// Act
	utilizations := service.GetLotUtilization()

	// Assert
	assert.Len(t, utilizations, 2)

	// Find LOT1 utilization
	var lot1Util *models.LotUtilization
	for _, util := range utilizations {
		if util.LotID == "LOT1" {
			lot1Util = util
			break
		}
	}

	assert.NotNil(t, lot1Util)
	assert.Equal(t, 1, lot1Util.OccupiedSpaces)
	assert.Equal(t, 1, lot1Util.AvailableSpaces)
	assert.Equal(t, 50.0, lot1Util.UtilizationRate)
}

func TestUC9_AttendantWithStrategy(t *testing.T) {
	// Arrange
	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	lots := []*models.ParkingLot{
		models.NewParkingLot("LOT1", 2),
		models.NewParkingLot("LOT2", 4),
	}
	car := models.NewCar("STRATEGY123", "Strategy Test")
	strategy := models.NewEvenDistributionStrategy()

	// Act
	decision, err := attendant.MakeParkingDecisionWithStrategy(lots, car, strategy)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, decision)
	assert.Equal(t, "LOT2", decision.LotID) // Should choose lot with more spaces
	assert.Equal(t, "ATT001", decision.AttendantID)
	assert.Contains(t, decision.Reason, "Even Distribution Strategy")
}

func TestUC9_InactiveAttendantCannotUseStrategy(t *testing.T) {
	// Arrange
	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	attendant.SetActive(false)

	lots := []*models.ParkingLot{models.NewParkingLot("LOT1", 2)}
	car := models.NewCar("TEST123", "Test")
	strategy := models.NewEvenDistributionStrategy()

	// Act
	decision, err := attendant.MakeParkingDecisionWithStrategy(lots, car, strategy)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, decision)
	assert.Equal(t, "attendant is not active", err.Error())
}

func TestUC9_StrategyFallbackToOriginal(t *testing.T) {
	// Arrange
	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	lots := []*models.ParkingLot{models.NewParkingLot("LOT1", 2)}
	car := models.NewCar("FALLBACK123", "Fallback Test")

	// Act - Pass nil strategy to test fallback
	decision, err := attendant.MakeParkingDecisionWithStrategy(lots, car, nil)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, decision)
	assert.Equal(t, "LOT1", decision.LotID)
	assert.Equal(t, "First available space strategy", decision.Reason)
}

func TestUC9_EvenDistributionWithMultipleCars(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 2)
	lot2 := models.NewParkingLot("LOT2", 2)
	service.AddLot(lot1)
	service.AddLot(lot2)

	attendant := models.NewParkingAttendant("ATT001", "Alice", "LOT1")
	service.AddAttendant(attendant)

	cars := []*models.Car{
		models.NewCar("CAR1", "Driver1"),
		models.NewCar("CAR2", "Driver2"),
		models.NewCar("CAR3", "Driver3"),
	}

	// Act - Park cars with even distribution
	var decisions []*models.ParkingDecision
	for _, car := range cars {
		decision, err := service.ParkCarEvenDistribution(car, "ATT001")
		if err == nil {
			decisions = append(decisions, decision)
		}
	}

	// Assert - Should distribute across lots
	assert.Len(t, decisions, 3)

	// Count cars per lot
	lotCounts := make(map[string]int)
	for _, decision := range decisions {
		lotCounts[decision.LotID]++
	}

	// Should have distributed cars across both lots
	assert.True(t, len(lotCounts) > 1, "Cars should be distributed across multiple lots")
}
