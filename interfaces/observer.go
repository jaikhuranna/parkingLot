package interfaces

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
    println("📢 OWNER NOTIFICATION: Parking lot", lotID, "is now FULL! Put out the full sign.")
}

func (o *OwnerObserver) OnLotAvailable(lotID string) {
    println("📢 OWNER NOTIFICATION: Parking lot", lotID, "has space available! Remove the full sign.")
}

