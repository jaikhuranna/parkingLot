package tests

import (
	"github.com/stretchr/testify/assert"
	"parking-lot-system/models"
	"parking-lot-system/services"
	"testing"
	"time"
)

func TestUC8_OwnerKnowsWhenCarWasParked(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	service.AddLot(lot)

	car := models.NewCar("ABC123", "John Doe")

	// Act - Park car and get ticket
	ticket, err := service.ParkCarWithTicket(car)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, ticket)
	assert.Equal(t, "ABC123", ticket.LicensePlate)
	assert.Equal(t, "LOT1", ticket.LotID)
	assert.True(t, ticket.IsActive)
	assert.False(t, ticket.ParkedAt.IsZero())
}

func TestUC8_ParkingTicketGeneration(t *testing.T) {
	// Arrange
	car := models.NewCar("XYZ789", "Jane Smith")

	// Act
	ticket := models.NewParkingTicket(car.LicensePlate, "LOT1", "A1")

	// Assert
	assert.NotEmpty(t, ticket.ID)
	assert.Equal(t, "XYZ789", ticket.LicensePlate)
	assert.Equal(t, "LOT1", ticket.LotID)
	assert.Equal(t, "A1", ticket.SpaceID)
	assert.True(t, ticket.IsActive)
	assert.True(t, ticket.UnparkedAt.IsZero())
}

func TestUC8_BillingServiceCalculation(t *testing.T) {
	// Arrange
	billingService := services.NewBillingService(10.0, 5.0) // $10/hour, $5 minimum

	// Test cases
	testCases := []struct {
		duration time.Duration
		expected float64
		name     string
	}{
		{30 * time.Minute, 5.0, "30 minutes - minimum charge"},
		{1 * time.Hour, 10.0, "1 hour - exact hourly rate"},
		{90 * time.Minute, 20.0, "1.5 hours - rounded up to 2 hours"},
		{2*time.Hour + 30*time.Minute, 30.0, "2.5 hours - rounded up to 3 hours"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			fee := billingService.CalculateFee(tc.duration)

			// Assert
			assert.Equal(t, tc.expected, fee)
		})
	}
}

func TestUC8_CompleteUnparkingWithBilling(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	service.AddLot(lot)

	car := models.NewCar("DEF456", "Bob Johnson")

	// Park the car
	ticket, err := service.ParkCarWithTicket(car)
	assert.NoError(t, err)

	// Simulate some parking time
	time.Sleep(10 * time.Millisecond)

	// Act - Unpark with billing
	unparkedCar, bill, err := service.UnparkCarWithBilling("DEF456")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, unparkedCar)
	assert.NotNil(t, bill)
	assert.Equal(t, "DEF456", unparkedCar.LicensePlate)
	assert.Equal(t, ticket.ID, bill.TicketID)
	assert.Equal(t, 5.0, bill.TotalAmount) // Minimum charge
	assert.False(t, ticket.IsActive)
}

func TestUC8_ParkingHistoryTracking(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 5)
	service.AddLot(lot)

	car := models.NewCar("GHI789", "Alice Brown")

	// Park and unpark car
	ticket, _ := service.ParkCarWithTicket(car)
	service.UnparkCarWithBilling("GHI789")

	// Act - Get parking history
	history, err := service.GetParkingHistory("GHI789")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, history, 1)
	assert.Equal(t, ticket.ID, history[0].ID)
}

func TestUC8_ActiveTicketRetrieval(t *testing.T) {
	// Arrange
	service := services.NewParkingService()
	lot := models.NewParkingLot("LOT1", 3)
	service.AddLot(lot)

	car := models.NewCar("JKL012", "Charlie Davis")

	// Park the car
	originalTicket, err := service.ParkCarWithTicket(car)
	assert.NoError(t, err)

	// Act - Get active ticket
	activeTicket, err := service.GetActiveTicket("JKL012")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, activeTicket)
	assert.Equal(t, originalTicket.ID, activeTicket.ID)
	assert.True(t, activeTicket.IsActive)
}

func TestUC8_NoActiveTicketError(t *testing.T) {
	// Arrange
	service := services.NewParkingService()

	// Act - Try to get active ticket for non-existent car
	ticket, err := service.GetActiveTicket("NONEXISTENT")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, ticket)
	assert.Equal(t, "no active ticket found for this vehicle", err.Error())
}

func TestUC8_BillGeneration(t *testing.T) {
	// Arrange
	billingService := services.NewBillingService(12.0, 6.0)

	car := models.NewCar("MNO345", "Test Driver")
	ticket := models.NewParkingTicket(car.LicensePlate, "LOT1", "B2")

	// Simulate parking duration
	time.Sleep(10 * time.Millisecond)
	ticket.CompleteParking()

	// Act
	bill := billingService.GenerateBill(ticket)

	// Assert
	assert.NotNil(t, bill)
	assert.Equal(t, ticket.ID, bill.TicketID)
	assert.Equal(t, "MNO345", bill.LicensePlate)
	assert.Equal(t, 6.0, bill.TotalAmount) // Minimum charge
	assert.Equal(t, 12.0, bill.HourlyRate)
	assert.Equal(t, 6.0, bill.MinimumCharge)
}

func TestUC8_TicketInformationRetrieval(t *testing.T) {
	// Arrange
	car := models.NewCar("PQR678", "Info Test")
	ticket := models.NewParkingTicketWithAttendant(car.LicensePlate, "LOT2", "C3", "ATT001")

	// Act
	info := ticket.GetTicketInfo()

	// Assert
	assert.Equal(t, ticket.ID, info["ID"])
	assert.Equal(t, "PQR678", info["LicensePlate"])
	assert.Equal(t, "LOT2", info["LotID"])
	assert.Equal(t, "C3", info["SpaceID"])
	assert.Equal(t, "ATT001", info["AttendantID"])
	assert.True(t, info["IsActive"].(bool))
	assert.NotNil(t, info["ParkedAt"])
	assert.NotNil(t, info["Duration"])
}

func TestUC8_ParkingDurationCalculation(t *testing.T) {
	// Arrange
	ticket := models.NewParkingTicket("TEST123", "LOT1", "A1")

	// Test active ticket duration
	time.Sleep(10 * time.Millisecond)
	duration1 := ticket.GetParkingDuration()
	assert.True(t, duration1 > 0)

	// Complete parking and test completed duration
	ticket.CompleteParking()
	duration2 := ticket.GetParkingDuration()

	// Duration should be calculated from parked to unparked time
	assert.True(t, duration2 > 0)
	assert.False(t, ticket.IsActive)
}
