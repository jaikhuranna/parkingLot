package models

type VehicleSize int

const (
	SmallVehicle  VehicleSize = iota // 0
	MediumVehicle                    // 1
	LargeVehicle                     // 2
)

type Car struct {
	LicensePlate string
	DriverName   string
	Size         VehicleSize
	IsHandicap   bool
}

func NewCar(licensePlate, driverName string) *Car {
	return &Car{
		LicensePlate: licensePlate,
		DriverName:   driverName,
		Size:         MediumVehicle, // FIXED: Default to MediumVehicle (1), not SmallVehicle (0)
		IsHandicap:   false,
	}
}

func (c *Car) SetVehicleSize(size VehicleSize) {
	c.Size = size
}

func (c *Car) SetHandicapStatus(isHandicap bool) {
	c.IsHandicap = isHandicap
}

func (c *Car) GetVehicleSizeString() string {
	switch c.Size {
	case SmallVehicle:
		return "Small"
	case MediumVehicle:
		return "Medium" // FIXED: Correct return value for MediumVehicle
	case LargeVehicle:
		return "Large"
	default:
		return "Medium" // FIXED: Default fallback to Medium
	}
}
