package services

import (
	"errors"
	"fmt"
	"parking-lot-system/interfaces"
	"parking-lot-system/models"
)

type ParkingService struct {
	lots            []*models.ParkingLot
	securityStaff   []*models.SecurityStaff
	attendants      []*models.ParkingAttendant
	defaultStrategy models.ParkingStrategy
}

func NewParkingService() *ParkingService {
	return &ParkingService{
		lots:            make([]*models.ParkingLot, 0),
		securityStaff:   make([]*models.SecurityStaff, 0),
		attendants:      make([]*models.ParkingAttendant, 0),
		defaultStrategy: models.NewEvenDistributionStrategy(),
	}
}

func (ps *ParkingService) AddLot(lot *models.ParkingLot) {
	ps.lots = append(ps.lots, lot)
}

// Security staff management
func (ps *ParkingService) AddSecurityStaff(staff *models.SecurityStaff) {
	ps.securityStaff = append(ps.securityStaff, staff)
}

func (ps *ParkingService) GetSecurityStaff() []*models.SecurityStaff {
	return ps.securityStaff
}

func (ps *ParkingService) FindSecurityStaffByID(staffID string) *models.SecurityStaff {
	for _, staff := range ps.securityStaff {
		if staff.ID == staffID {
			return staff
		}
	}
	return nil
}

func (ps *ParkingService) AssignSecurityToLot(staffID, lotID string) error {
	staff := ps.FindSecurityStaffByID(staffID)
	if staff == nil {
		return errors.New("security staff not found")
	}

	lot := ps.findLotByID(lotID)
	if lot == nil {
		return errors.New("lot not found")
	}

	staff.AssignToLot(lotID)
	return nil
}

// Attendant management
func (ps *ParkingService) AddAttendant(attendant *models.ParkingAttendant) {
	ps.attendants = append(ps.attendants, attendant)
}

func (ps *ParkingService) GetAttendants() []*models.ParkingAttendant {
	return ps.attendants
}

func (ps *ParkingService) FindAttendantByID(attendantID string) *models.ParkingAttendant {
	for _, attendant := range ps.attendants {
		if attendant.ID == attendantID {
			return attendant
		}
	}
	return nil
}

func (ps *ParkingService) ParkCarWithAttendant(car *models.Car, attendantID string) (*models.ParkingDecision, error) {
	if car == nil {
		return nil, errors.New("car cannot be nil")
	}

	attendant := ps.FindAttendantByID(attendantID)
	if attendant == nil {
		return nil, errors.New("attendant not found")
	}

	decision, err := attendant.MakeParkingDecision(ps.lots, car)
	if err != nil {
		return nil, err
	}

	// Execute the parking decision
	lot := ps.findLotByID(decision.LotID)
	if lot == nil {
		return nil, errors.New("lot specified in decision not found")
	}

	err = lot.ParkCar(car)
	if err != nil {
		return nil, err
	}

	return decision, nil
}

func (ps *ParkingService) AddObserverToLot(lotID string, observer interfaces.ParkingLotObserver) error {
	lot := ps.findLotByID(lotID)
	if lot == nil {
		return errors.New("lot not found")
	}
	lot.AddObserver(observer)
	return nil
}

func (ps *ParkingService) RemoveObserverFromLot(lotID string, observer interfaces.ParkingLotObserver) error {
	lot := ps.findLotByID(lotID)
	if lot == nil {
		return errors.New("lot not found")
	}
	lot.RemoveObserver(observer)
	return nil
}

func (ps *ParkingService) ParkCar(car *models.Car) error {
	if car == nil {
		return errors.New("car cannot be nil")
	}

	for _, lot := range ps.lots {
		if err := lot.ParkCar(car); err == nil {
			return nil
		}
	}

	return errors.New("no available parking space")
}

func (ps *ParkingService) UnparkCar(licensePlate string) (*models.Car, error) {
	if licensePlate == "" {
		return nil, errors.New("license plate cannot be empty")
	}

	for _, lot := range ps.lots {
		if car, err := lot.UnparkCar(licensePlate); err == nil {
			return car, nil
		}
	}

	return nil, errors.New("car not found in any parking lot")
}

func (ps *ParkingService) FindCar(licensePlate string) (*models.ParkingSpace, error) {

	if licensePlate == "" {
		return nil, errors.New("license plate cannot be empty")
	}

	for _, lot := range ps.lots {
		if space := lot.FindCar(licensePlate); space != nil {
			return space, nil
		}
	}

	return nil, errors.New("car not found")
}

func (ps *ParkingService) GetLotStatus(lotID string) (*models.ParkingLot, error) {
	lot := ps.findLotByID(lotID)
	if lot == nil {
		return nil, errors.New("lot not found")
	}
	return lot, nil
}

func (ps *ParkingService) IsAnyLotFull() bool {
	for _, lot := range ps.lots {
		if lot.IsFull() {
			return true
		}
	}
	return false
}

func (ps *ParkingService) findLotByID(lotID string) *models.ParkingLot {
	for _, lot := range ps.lots {
		if lot.ID == lotID {
			return lot
		}
	}
	return nil
}

// UC7: Enhanced car finding functionality

// UC7: Enhanced car finding functionality
func (ps *ParkingService) FindCarWithLocation(licensePlate string) (*models.CarLocation, error) {
	if licensePlate == "" {
		return nil, errors.New("license plate cannot be empty")
	}

	for _, lot := range ps.lots {
		if space := lot.FindCar(licensePlate); space != nil {
			// Convert space.ID from int to string
			spaceIDStr := fmt.Sprintf("%d", space.ID)

			// Extract row and position from space ID if available
			row := ""
			position := 0
			if len(spaceIDStr) > 0 {
				row = string(spaceIDStr[0]) // First character as row
				if len(spaceIDStr) > 1 {
					// Try to parse position from remaining characters
					if pos := spaceIDStr[1:]; len(pos) > 0 {
						position = int(pos[0] - '0')
					}
				}
			}

			return models.NewCarLocation(
				space.ParkedCar,
				lot.ID,
				spaceIDStr, // Now correctly passing string
				row,
				position,
				"", // Attendant ID - could be enhanced later
			), nil
		}
	}

	return nil, errors.New("car not found")
}

func (ps *ParkingService) ProvideDirectionsToDriver(licensePlate string) (string, error) {
	location, err := ps.FindCarWithLocation(licensePlate)
	if err != nil {
		return "", err
	}

	directions := fmt.Sprintf(
		"🗺️  Your car %s is located in:\n"+
			"   📍 Lot: %s\n"+
			"   🅿️  Space: %s\n"+
			"   📝 Row: %s, Position: %d\n"+
			"   ⏰ Parked at: %s",
		licensePlate,
		location.LotID,
		location.SpaceID,
		location.Row,
		location.Position,
		location.ParkedAt.Format("2006-01-02 15:04:05"),
	)

	return directions, nil
}

// UC8: Billing and time tracking functionality
type TicketManager struct {
	tickets map[string]*models.ParkingTicket
}

func NewTicketManager() *TicketManager {
	return &TicketManager{
		tickets: make(map[string]*models.ParkingTicket),
	}
}

var globalTicketManager = NewTicketManager()

func (ps *ParkingService) ParkCarWithTicket(car *models.Car) (*models.ParkingTicket, error) {
	if car == nil {
		return nil, errors.New("car cannot be nil")
	}

	for _, lot := range ps.lots {
		if !lot.IsFull() {
			space := lot.FindAvailableSpace()
			if space != nil {
				err := lot.ParkCar(car)
				if err != nil {
					continue
				}

				// Create and store ticket - Convert space.ID to string
				spaceIDStr := fmt.Sprintf("%d", space.ID)
				ticket := models.NewParkingTicket(car.LicensePlate, lot.ID, spaceIDStr)
				globalTicketManager.tickets[ticket.ID] = ticket

				return ticket, nil
			}
		}
	}

	return nil, errors.New("no available parking space")
}

func (ps *ParkingService) UnparkCarWithBilling(licensePlate string) (*models.Car, *Bill, error) {
	if licensePlate == "" {
		return nil, nil, errors.New("license plate cannot be empty")
	}

	// Find and complete ticket
	var ticket *models.ParkingTicket
	for _, t := range globalTicketManager.tickets {
		if t.LicensePlate == licensePlate && t.IsActive {
			ticket = t
			break
		}
	}

	if ticket == nil {
		return nil, nil, errors.New("active ticket not found for car")
	}

	// Unpark the car
	car, err := ps.UnparkCar(licensePlate)
	if err != nil {
		return nil, nil, err
	}

	// Complete ticket and generate bill
	ticket.CompleteParking()
	billingService := NewBillingService(10.0, 5.0) // $10/hour, $5 minimum
	bill := billingService.GenerateBill(ticket)

	return car, bill, nil
}
func (ps *ParkingService) GetParkingHistory(licensePlate string) ([]*models.ParkingTicket, error) {
	var history []*models.ParkingTicket

	for _, ticket := range globalTicketManager.tickets {
		if ticket.LicensePlate == licensePlate {
			history = append(history, ticket)
		}
	}

	if len(history) == 0 {
		return nil, errors.New("no parking history found for this vehicle")
	}

	return history, nil
}

func (ps *ParkingService) GetActiveTicket(licensePlate string) (*models.ParkingTicket, error) {
	for _, ticket := range globalTicketManager.tickets {
		if ticket.LicensePlate == licensePlate && ticket.IsActive {
			return ticket, nil
		}
	}

	return nil, errors.New("no active ticket found for this vehicle")
}

// UC9: Even distribution parking strategy
func (ps *ParkingService) SetDefaultStrategy(strategy models.ParkingStrategy) {
	ps.defaultStrategy = strategy
}

func (ps *ParkingService) ParkCarWithStrategy(car *models.Car, attendantID string, strategy models.ParkingStrategy) (*models.ParkingDecision, error) {
	if car == nil {
		return nil, errors.New("car cannot be nil")
	}

	attendant := ps.FindAttendantByID(attendantID)
	if attendant == nil {
		return nil, errors.New("attendant not found")
	}

	// Use strategy-based decision making
	decision, err := attendant.MakeParkingDecisionWithStrategy(ps.lots, car, strategy)
	if err != nil {
		return nil, err
	}

	// Execute the parking decision
	lot := ps.findLotByID(decision.LotID)
	if lot == nil {
		return nil, errors.New("lot specified in decision not found")
	}

	err = lot.ParkCar(car)
	if err != nil {
		return nil, err
	}

	return decision, nil
}

func (ps *ParkingService) GetLotUtilization() []*models.LotUtilization {
	var utilizations []*models.LotUtilization

	for _, lot := range ps.lots {
		utilization := models.CalculateLotUtilization(lot)
		utilizations = append(utilizations, utilization)
	}

	return utilizations
}

func (ps *ParkingService) ParkCarEvenDistribution(car *models.Car, attendantID string) (*models.ParkingDecision, error) {
	strategy := models.NewEvenDistributionStrategy()
	return ps.ParkCarWithStrategy(car, attendantID, strategy)
}

// UC10-UC11: Advanced parking strategies
func (ps *ParkingService) ParkHandicapCar(car *models.Car, attendantID string) (*models.ParkingDecision, error) {
	if !car.IsHandicap {
		return nil, errors.New("car is not registered as handicap vehicle")
	}

	strategy := models.NewHandicapPriorityStrategy()
	return ps.ParkCarWithStrategy(car, attendantID, strategy)
}

func (ps *ParkingService) ParkLargeVehicle(car *models.Car, attendantID string) (*models.ParkingDecision, error) {
	if car.Size != models.LargeVehicle {
		return nil, errors.New("car is not classified as large vehicle")
	}

	strategy := models.NewLargeVehicleStrategy()
	return ps.ParkCarWithStrategy(car, attendantID, strategy)
}

func (ps *ParkingService) ParkCarSmart(car *models.Car, attendantID string) (*models.ParkingDecision, error) {
	strategy := models.NewSmartParkingStrategy()
	return ps.ParkCarWithStrategy(car, attendantID, strategy)
}

func (ps *ParkingService) GetHandicapSpacesCount() map[string]int {
	handicapCounts := make(map[string]int)

	for _, lot := range ps.lots {
		count := 0
		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil && space.ParkedCar.IsHandicap {
				count++
			}
		}
		handicapCounts[lot.ID] = count
	}

	return handicapCounts
}

func (ps *ParkingService) GetLargeVehicleSpacesCount() map[string]int {
	largeCounts := make(map[string]int)

	for _, lot := range ps.lots {
		count := 0
		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil && space.ParkedCar.Size == models.LargeVehicle {
				count++
			}
		}
		largeCounts[lot.ID] = count
	}

	return largeCounts
}

func (ps *ParkingService) GetDetailedLotAnalytics() map[string]map[string]interface{} {
	analytics := make(map[string]map[string]interface{})

	for _, lot := range ps.lots {
		util := models.CalculateLotUtilization(lot)

		handicapCount := 0
		largeVehicleCount := 0
		smallVehicleCount := 0

		for _, space := range lot.Spaces {
			if space.IsOccupied && space.ParkedCar != nil {
				if space.ParkedCar.IsHandicap {
					handicapCount++
				}

				switch space.ParkedCar.Size {
				case models.SmallVehicle:
					smallVehicleCount++
				case models.LargeVehicle:
					largeVehicleCount++
				}
			}
		}

		analytics[lot.ID] = map[string]interface{}{
			"TotalSpaces":      util.TotalSpaces,
			"OccupiedSpaces":   util.OccupiedSpaces,
			"AvailableSpaces":  util.AvailableSpaces,
			"UtilizationRate":  util.UtilizationRate,
			"HandicapVehicles": handicapCount,
			"LargeVehicles":    largeVehicleCount,
			"SmallVehicles":    smallVehicleCount,
		}
	}

	return analytics
}

// UC11: Large vehicle management functionality
func (ps *ParkingService) GetBestLotForLargeVehicle() (*models.ParkingLot, int, error) {
	var bestLot *models.ParkingLot
	maxSpaces := -1

	for _, lot := range ps.lots {
		if !lot.IsFull() {
			available := lot.GetAvailableSpaces()
			if available > maxSpaces {
				maxSpaces = available
				bestLot = lot
			}
		}
	}

	if bestLot == nil {
		return nil, 0, errors.New("no lots available for large vehicles")
	}

	return bestLot, maxSpaces, nil
}

func (ps *ParkingService) GetLargeVehicleRecommendations() map[string]interface{} {
	recommendations := make(map[string]interface{})

	bestLot, maxSpaces, err := ps.GetBestLotForLargeVehicle()
	if err != nil {
		recommendations["error"] = err.Error()
		return recommendations
	}

	analytics := ps.GetDetailedLotAnalytics()

	recommendations["recommendedLot"] = bestLot.ID
	recommendations["availableSpaces"] = maxSpaces

	if lotAnalytics, exists := analytics[bestLot.ID]; exists {
		recommendations["currentLargeVehicles"] = lotAnalytics["LargeVehicles"]
		recommendations["utilizationRate"] = lotAnalytics["UtilizationRate"]
	}

	// Find alternative lots suitable for large vehicles
	var alternatives []string
	for _, lot := range ps.lots {
		if lot.ID != bestLot.ID && !lot.IsFull() && lot.GetAvailableSpaces() >= 3 {
			alternatives = append(alternatives, lot.ID)
		}
	}
	recommendations["alternativeLots"] = alternatives

	return recommendations
}

func (ps *ParkingService) ValidateLargeVehicleCapacity() map[string]bool {
	validation := make(map[string]bool)

	for _, lot := range ps.lots {
		utilizationRate := float64(lot.GetOccupiedSpaces()) / float64(lot.Capacity) * 100
		validation[lot.ID] = utilizationRate <= 70.0 // Suitable if 70% or less occupied
	}

	return validation
}

func (ps *ParkingService) GetOptimalLargeVehiclePlacement() map[string]interface{} {
	placement := make(map[string]interface{})

	bestLot, maxSpaces, err := ps.GetBestLotForLargeVehicle()
	if err != nil {
		placement["error"] = err.Error()
		return placement
	}

	// Calculate maneuvering space efficiency
	efficiency := float64(maxSpaces) / float64(bestLot.Capacity) * 100

	placement["optimalLot"] = bestLot.ID
	placement["availableSpaces"] = maxSpaces
	placement["maneuveringEfficiency"] = efficiency
	placement["recommendedForLargeVehicles"] = efficiency >= 30.0 // At least 30% free space

	return placement
}
