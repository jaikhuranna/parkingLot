package interfaces

import "fmt"

// Observer pattern for parking lot notifications
type ParkingLotObserver interface {
	OnLotFull(lotID string)
	OnLotAvailable(lotID string)
}

// Owner observer implementation
type OwnerObserver struct {
	OwnerName string
}

func NewOwnerObserver(ownerName string) *OwnerObserver {
	return &OwnerObserver{
		OwnerName: ownerName,
	}
}

func (o *OwnerObserver) OnLotFull(lotID string) {
	fmt.Printf("📢 OWNER NOTIFICATION: Parking lot %s is now FULL! Put out the full sign.\n", lotID)
}

func (o *OwnerObserver) OnLotAvailable(lotID string) {
	fmt.Printf("📢 OWNER NOTIFICATION: Parking lot %s has space available! Remove the full sign.\n", lotID)
}

// NEW: Security observer implementation
type SecurityObserver struct {
	SecurityStaffName string
	StaffID           string
}

func NewSecurityObserver(staffName, staffID string) *SecurityObserver {
	return &SecurityObserver{
		SecurityStaffName: staffName,
		StaffID:           staffID,
	}
}

func (s *SecurityObserver) OnLotFull(lotID string) {
	fmt.Printf("🚨 SECURITY ALERT: Parking lot %s is FULL! Redirect security staff to manage overflow traffic.\n", lotID)
	fmt.Printf("Security Staff %s (ID: %s) - Deploy to alternate parking areas.\n", s.SecurityStaffName, s.StaffID)
}

func (s *SecurityObserver) OnLotAvailable(lotID string) {
	fmt.Printf("✅ SECURITY UPDATE: Parking lot %s has available spaces. Normal traffic flow resumed.\n", lotID)
	fmt.Printf("Security Staff %s (ID: %s) - Return to regular positions.\n", s.SecurityStaffName, s.StaffID)
}
