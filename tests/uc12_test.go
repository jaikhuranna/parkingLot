package tests

import (
	"github.com/stretchr/testify/assert"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"testing"
)

func TestUC12_FindWhiteCars(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 5)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	// Create test cars
	whiteCar := models.NewCar("WHITE001", "John Smith")
	whiteCar.SetColor("White")
	whiteCar.SetMake("Honda")

	blueCar := models.NewCar("BLUE001", "Jane Doe")
	blueCar.SetColor("Blue")
	blueCar.SetMake("Toyota")

	// Park cars
	lot.ParkCar(whiteCar)
	lot.ParkCar(blueCar)

	// Act
	whiteCars, err := policeService.FindWhiteCars()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, whiteCars, 1)
	assert.Equal(t, "WHITE001", whiteCars[0].Car.LicensePlate)
	assert.Equal(t, "White", whiteCars[0].Car.Color)
	assert.Equal(t, "LOT1", whiteCars[0].LotID)
}

func TestUC12_FindWhiteCarsMultipleLots(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot1 := models.NewParkingLot("LOT1", 3)
	lot2 := models.NewParkingLot("LOT2", 3)
	parkingService.AddLot(lot1)
	parkingService.AddLot(lot2)

	policeService := services.NewPoliceService(parkingService)

	// Create white cars in different lots
	whiteCar1 := models.NewCar("WHITE001", "Driver1")
	whiteCar1.SetColor("White")
	whiteCar1.SetMake("Honda")

	whiteCar2 := models.NewCar("WHITE002", "Driver2")
	whiteCar2.SetColor("White")
	whiteCar2.SetMake("Ford")

	lot1.ParkCar(whiteCar1)
	lot2.ParkCar(whiteCar2)

	// Act
	whiteCars, err := policeService.FindWhiteCars()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, whiteCars, 2)

	// Verify both cars are found
	plates := []string{whiteCars[0].Car.LicensePlate, whiteCars[1].Car.LicensePlate}
	assert.Contains(t, plates, "WHITE001")
	assert.Contains(t, plates, "WHITE002")
}

func TestUC12_NoWhiteCarsFound(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 5)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	// Park only non-white cars
	blueCar := models.NewCar("BLUE001", "Driver1")
	blueCar.SetColor("Blue")
	lot.ParkCar(blueCar)

	redCar := models.NewCar("RED001", "Driver2")
	redCar.SetColor("Red")
	lot.ParkCar(redCar)

	// Act
	whiteCars, err := policeService.FindWhiteCars()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, whiteCars, 0)
}

func TestUC12_CaseInsensitiveColorMatching(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 5)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	// Create cars with different case colors
	car1 := models.NewCar("TEST001", "Driver1")
	car1.SetColor("WHITE") // Uppercase

	car2 := models.NewCar("TEST002", "Driver2")
	car2.SetColor("white") // Lowercase

	car3 := models.NewCar("TEST003", "Driver3")
	car3.SetColor("White") // Title case

	lot.ParkCar(car1)
	lot.ParkCar(car2)
	lot.ParkCar(car3)

	// Act
	whiteCars, err := policeService.FindWhiteCars()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, whiteCars, 3) // All should be found regardless of case
}

func TestUC12_CarDetailsRetrieval(t *testing.T) {
	// Arrange
	car := models.NewCar("TEST123", "Test Driver")
	car.SetColor("Red")
	car.SetMake("BMW")
	car.SetVehicleSize(models.LargeVehicle)
	car.SetHandicapStatus(true)

	// Act
	details := car.GetCarDetails()

	// Assert
	assert.Equal(t, "TEST123", details["LicensePlate"])
	assert.Equal(t, "Test Driver", details["DriverName"])
	assert.Equal(t, "Red", details["Color"])
	assert.Equal(t, "BMW", details["Make"])
	assert.Equal(t, "Large", details["Size"])
	assert.True(t, details["IsHandicap"].(bool))
}

func TestUC12_InvestigationReportGeneration(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	policeService := services.NewPoliceService(parkingService)

	car := models.NewCar("REPORT001", "Report Driver")
	car.SetColor("White")
	car.SetMake("Honda")

	vehicles := []*services.VehicleInvestigationInfo{
		{
			Car:     car,
			LotID:   "LOT1",
			SpaceID: "1",
		},
	}

	// Act
	report := policeService.GenerateInvestigationReport(vehicles, "Test Investigation")

	// Assert
	assert.Contains(t, report, "POLICE INVESTIGATION REPORT")
	assert.Contains(t, report, "Test Investigation")
	assert.Contains(t, report, "REPORT001")
	assert.Contains(t, report, "Report Driver")
	assert.Contains(t, report, "White")
	assert.Contains(t, report, "Honda")
	assert.Contains(t, report, "LOT1")
	assert.Contains(t, report, "Space 1")
}

func TestUC12_GeneralColorMakeSearch(t *testing.T) {
	// Arrange
	parkingService := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 5)
	parkingService.AddLot(lot)

	policeService := services.NewPoliceService(parkingService)

	// Create various cars
	car1 := models.NewCar("RED_BMW", "Driver1")
	car1.SetColor("Red")
	car1.SetMake("BMW")

	car2 := models.NewCar("BLUE_BMW", "Driver2")
	car2.SetColor("Blue")
	car2.SetMake("BMW")

	car3 := models.NewCar("RED_HONDA", "Driver3")
	car3.SetColor("Red")
	car3.SetMake("Honda")

	lot.ParkCar(car1)
	lot.ParkCar(car2)
	lot.ParkCar(car3)

	// Act - Search for BMW cars
	bmwCars, err := policeService.FindCarsByColorAndMake("", "BMW")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, bmwCars, 2)

	// Act - Search for red cars
	redCars, err := policeService.FindCarsByColorAndMake("Red", "")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, redCars, 2)

	// Act - Search for red BMW cars
	redBMWs, err := policeService.FindCarsByColorAndMake("Red", "BMW")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, redBMWs, 1)
	assert.Equal(t, "RED_BMW", redBMWs[0].Car.LicensePlate)
}
