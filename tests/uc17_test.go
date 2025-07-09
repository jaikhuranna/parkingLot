package tests

import (
	"testing"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"github.com/stretchr/testify/assert"
)

func TestUC17_GetAllCarsInLot(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 5)
	lot2 := models.NewParkingLot("LOT2", 3)
	parkingService.AddLot(lot1)
	parkingService.AddLot(lot2)

	policeService := services.NewPoliceService(parkingService)

	car1 := models.NewCar("CAR001", "Driver1")
	car2 := models.NewCar("CAR002", "Driver2")
	lot1.ParkCar(car1)
	lot1.ParkCar(car2)

	car3 := models.NewCar("CAR003", "Driver3")
	lot2.ParkCar(car3)

	// Act
	carsInLot1, err := policeService.GetAllCarsInLot("LOT1")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, carsInLot1, 2)
	
	plates := []string{carsInLot1[0].Car.LicensePlate, carsInLot1[1].Car.LicensePlate}
	assert.Contains(t, plates, "CAR001")
	assert.Contains(t, plates, "CAR002")
}

func TestUC17_DetectFraudulentPlates(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 10)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	fraudulentCars := []*models.Car{
		models.NewCar("FAKE123", "Fake Driver"),
		models.NewCar("TEST001", "Test Driver"),
		models.NewCar("XXXX000", "Suspicious Driver"),
	}

	legitimateCars := []*models.Car{
		models.NewCar("ABC1234", "Legit Driver1"),
		models.NewCar("XYZ9876", "Legit Driver2"),
	}

	for _, car := range fraudulentCars {
		lot.ParkCar(car)
	}
	for _, car := range legitimateCars {
		lot.ParkCar(car)
	}

	// Act
	suspiciousVehicles, err := policeService.DetectFraudulentPlates()

	// Assert
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(suspiciousVehicles), 3)

	suspiciousPlates := make([]string, len(suspiciousVehicles))
	for i, vehicle := range suspiciousVehicles {
		suspiciousPlates[i] = vehicle.Car.LicensePlate
	}

	assert.Contains(t, suspiciousPlates, "FAKE123")
	assert.Contains(t, suspiciousPlates, "TEST001")
	assert.Contains(t, suspiciousPlates, "XXXX000")
}

func TestUC17_GenerateCompleteLotInvestigationReport(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 5)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	legitimateCar := models.NewCar("LEGIT123", "Honest Driver")
	fraudCar := models.NewCar("FAKE999", "Fraud Driver")

	lot.ParkCar(legitimateCar)
	lot.ParkCar(fraudCar)

	// Act
	report := policeService.GenerateCompleteLotInvestigationReport("LOT1")

	// Assert
	assert.Contains(t, report, "COMPLETE PARKING LOT INVESTIGATION REPORT")
	assert.Contains(t, report, "Target Lot: LOT1")
	assert.Contains(t, report, "Total Vehicles in Lot: 2")
	assert.Contains(t, report, "LEGIT123")
	assert.Contains(t, report, "FAKE999")
	assert.Contains(t, report, "FRAUD ALERT")
}

func TestUC17_GetCompleteInvestigationSummary(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 10)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	whiteCar := models.NewCar("WHITE01", "White Driver")
	whiteCar.SetColor("White")
	whiteCar.SetMake("Honda")

	blueToyota := models.NewCar("BLUE_TOY", "Toyota Driver")
	blueToyota.SetColor("Blue")
	blueToyota.SetMake("Toyota")

	bmwCar := models.NewCar("BMW_001", "BMW Driver")
	bmwCar.SetMake("BMW")

	fraudCar := models.NewCar("FAKE001", "Fraud Driver")

	lot.ParkCar(whiteCar)
	lot.ParkCar(blueToyota)
	lot.ParkCar(bmwCar)
	lot.ParkCar(fraudCar)

	// Act
	summary := policeService.GetCompleteInvestigationSummary()

	// Assert
	assert.Equal(t, 1, summary["whiteCarsCount"])
	assert.Equal(t, 1, summary["blueToyotasCount"])
	assert.Equal(t, 1, summary["bmwCarsCount"])
	assert.Equal(t, 1, summary["fraudulentPlatesCount"])
	assert.Equal(t, 6, summary["totalInvestigationTypes"])
	assert.NotNil(t, summary["investigationTimestamp"])
}
