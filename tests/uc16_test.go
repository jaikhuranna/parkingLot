package tests

import (
	"testing"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"github.com/stretchr/testify/assert"
)

func TestUC16_FindHandicapCarsInRows(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 100)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	// Create handicap cars in different spaces
	handicapCarRowB := models.NewCar("HANDICAP_B", "Driver B")
	handicapCarRowB.SetHandicapStatus(true)
	handicapCarRowB.SetVehicleSize(models.SmallVehicle)

	handicapCarRowA := models.NewCar("HANDICAP_A", "Driver A")
	handicapCarRowA.SetHandicapStatus(true)

	regularCarRowB := models.NewCar("REGULAR_B", "Regular Driver")

	// Park in specific spaces to control row assignment
	lot.Spaces[30].Park(handicapCarRowB) // Space 31 = Row B
	lot.Spaces[10].Park(handicapCarRowA) // Space 11 = Row A
	lot.Spaces[35].Park(regularCarRowB)  // Space 36 = Row B, but not handicap

	// Act
	handicapInBD, err := policeService.FindHandicapCarsInRows([]string{"B", "D"})

	// Assert
	assert.NoError(t, err)
	assert.Len(t, handicapInBD, 1)
	assert.Equal(t, "HANDICAP_B", handicapInBD[0].Car.LicensePlate)
}

func TestUC16_GetVehiclesByLocationCriteria(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 100)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	// Create test vehicle
	smallHandicapCar := models.NewCar("SMALL_H", "Small Handicap Driver")
	smallHandicapCar.SetVehicleSize(models.SmallVehicle)
	smallHandicapCar.SetHandicapStatus(true)

	lot.Spaces[40].Park(smallHandicapCar) // Space 41 = Row B

	// Act
	vehicles, err := policeService.GetVehiclesByLocationCriteria(
		models.SmallVehicle, true, []string{"B"})

	// Assert
	assert.NoError(t, err)
	assert.Len(t, vehicles, 1)
	assert.Equal(t, "SMALL_H", vehicles[0].Car.LicensePlate)
	assert.True(t, vehicles[0].Car.IsHandicap)
}

func TestUC16_ValidateHandicapPermitFraud(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 100)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	// Add handicap cars in suspicious rows
	suspiciousCar := models.NewCar("FRAUD001", "Suspicious Driver")
	suspiciousCar.SetHandicapStatus(true)
	lot.Spaces[45].Park(suspiciousCar) // Row B

	// Act
	validation := policeService.ValidateHandicapPermitFraud()

	// Assert
	assert.Equal(t, 1, validation["totalHandicapVehicles"])
	assert.Equal(t, 1, validation["vehiclesInRowsB_D"])
	assert.True(t, validation["investigationRequired"].(bool))
	assert.Contains(t, validation["fraudRisk"].(string), "Risk")
}

func TestUC16_GenerateHandicapFraudInvestigationReport(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 100)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	fraudCar := models.NewCar("FRAUD_REPORT", "Report Driver")
	fraudCar.SetHandicapStatus(true)
	fraudCar.SetVehicleSize(models.SmallVehicle)
	lot.Spaces[60].Park(fraudCar) // Row C (not in B/D, so shouldn't appear)

	suspiciousCar := models.NewCar("SUSPICIOUS", "Suspicious Driver")
	suspiciousCar.SetHandicapStatus(true)
	lot.Spaces[30].Park(suspiciousCar) // Row B

	// Act
	report := policeService.GenerateHandicapFraudInvestigationReport()

	// Assert
	assert.Contains(t, report, "HANDICAP PERMIT FRAUD INVESTIGATION REPORT")
	assert.Contains(t, report, "Rows B and D")
	assert.Contains(t, report, "SUSPICIOUS")
	assert.NotContains(t, report, "FRAUD_REPORT") // Should not be in B/D rows
	assert.Contains(t, report, "INVESTIGATION RECOMMENDATIONS")
}

func TestUC16_ParkingSpaceRowAssignment(t *testing.T) {
	// Test row assignment logic
	space1 := models.NewParkingSpace(15)  // Should be Row A
	space2 := models.NewParkingSpace(35)  // Should be Row B
	space3 := models.NewParkingSpace(65)  // Should be Row C
	space4 := models.NewParkingSpace(85)  // Should be Row D

	assert.Equal(t, "A", space1.GetRowAssignment())
	assert.Equal(t, "B", space2.GetRowAssignment())
	assert.Equal(t, "C", space3.GetRowAssignment())
	assert.Equal(t, "D", space4.GetRowAssignment())
}

func TestUC16_GetLocationStatistics(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 100)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	// Add vehicles to different rows
	carRowA := models.NewCar("ROW_A", "Driver A")
	carRowA.SetVehicleSize(models.SmallVehicle)
	lot.Spaces[10].Park(carRowA) // Row A

	carRowB := models.NewCar("ROW_B", "Driver B")
	carRowB.SetVehicleSize(models.LargeVehicle)
	carRowB.SetHandicapStatus(true)
	lot.Spaces[30].Park(carRowB) // Row B

	// Act
	stats := policeService.GetLocationStatistics()

	// Assert
	rowCounts := stats["totalVehiclesByRow"].(map[string]int)
	assert.Equal(t, 1, rowCounts["A"])
	assert.Equal(t, 1, rowCounts["B"])
	assert.Equal(t, 0, rowCounts["C"])
	assert.Equal(t, 0, rowCounts["D"])

	handicapCounts := stats["handicapVehiclesByRow"].(map[string]int)
	assert.Equal(t, 0, handicapCounts["A"])
	assert.Equal(t, 1, handicapCounts["B"])
}

func TestUC16_NoHandicapVehiclesInSuspiciousRows(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 100)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	// Add handicap car only in Row A (not suspicious)
	legitimateCar := models.NewCar("LEGIT001", "Legitimate Driver")
	legitimateCar.SetHandicapStatus(true)
	lot.Spaces[10].Park(legitimateCar) // Row A

	// Act
	validation := policeService.ValidateHandicapPermitFraud()

	// Assert
	assert.Equal(t, 1, validation["totalHandicapVehicles"])
	assert.Equal(t, 0, validation["vehiclesInRowsB_D"])
	assert.False(t, validation["investigationRequired"].(bool))
	assert.Contains(t, validation["fraudRisk"].(string), "No Risk")
}
