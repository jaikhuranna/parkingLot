package models

type Car struct {
    LicensePlate string
    DriverName   string
}

func NewCar(licensePlate, driverName string) *Car {
    return &Car{
        LicensePlate: licensePlate,
        DriverName:   driverName,
    }
}

