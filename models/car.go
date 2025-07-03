package models

type Car struct {
    Size       VehicleSize
    IsHandicap bool
	LicensePlate string
	DriverName   string
}

func NewCar(licensePlate, driverName string) *Car {
	return &Car{
		LicensePlate: licensePlate,
		DriverName:   driverName,
	}
}

// Enhanced car properties for advanced parking strategies
type VehicleSize int

const (
    SmallVehicle VehicleSize = iota
    MediumVehicle
    LargeVehicle
)

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
        return "Medium"
    case LargeVehicle:
        return "Large"
    default:
        return "Medium"
    }
}
