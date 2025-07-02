package services

import (
    "parking-lot-system/models"
    "time"
    "math"
    "fmt"
)

type BillingService struct {
    HourlyRate    float64
    MinimumCharge float64
}

func NewBillingService(hourlyRate, minimumCharge float64) *BillingService {
    return &BillingService{
        HourlyRate:    hourlyRate,
        MinimumCharge: minimumCharge,
    }
}

type Bill struct {
    TicketID        string
    LicensePlate    string
    ParkedAt        time.Time
    UnparkedAt      time.Time
    Duration        time.Duration
    HourlyRate      float64
    TotalAmount     float64
    MinimumCharge   float64
}

func (bs *BillingService) CalculateFee(duration time.Duration) float64 {
    hours := duration.Hours()
    
    if hours < 1.0 {
        return bs.MinimumCharge
    }
    
    // Round up to next hour for billing
    billingHours := math.Ceil(hours)
    total := billingHours * bs.HourlyRate
    
    if total < bs.MinimumCharge {
        return bs.MinimumCharge
    }
    
    return total
}

func (bs *BillingService) GenerateBill(ticket *models.ParkingTicket) *Bill {
    duration := ticket.GetParkingDuration()
    totalAmount := bs.CalculateFee(duration)
    
    return &Bill{
        TicketID:      ticket.ID,
        LicensePlate:  ticket.LicensePlate,
        ParkedAt:      ticket.ParkedAt,
        UnparkedAt:    ticket.UnparkedAt,
        Duration:      duration,
        HourlyRate:    bs.HourlyRate,
        TotalAmount:   totalAmount,
        MinimumCharge: bs.MinimumCharge,
    }
}

func (b *Bill) GetBillSummary() map[string]interface{} {
    return map[string]interface{}{
        "TicketID":      b.TicketID,
        "LicensePlate":  b.LicensePlate,
        "ParkedAt":      b.ParkedAt.Format("2006-01-02 15:04:05"),
        "UnparkedAt":    b.UnparkedAt.Format("2006-01-02 15:04:05"),
        "Duration":      b.Duration.String(),
        "HourlyRate":    b.HourlyRate,
        "TotalAmount":   b.TotalAmount,
        "MinimumCharge": b.MinimumCharge,
    }
}

func (b *Bill) PrintBill() string {
    return fmt.Sprintf(`
=================================
         PARKING BILL
=================================
Ticket ID: %s
License Plate: %s
Parked At: %s
Unparked At: %s
Duration: %s
Hourly Rate: $%.2f
Total Amount: $%.2f
=================================
Thank you for using our parking!
`, 
        b.TicketID,
        b.LicensePlate,
        b.ParkedAt.Format("2006-01-02 15:04:05"),
        b.UnparkedAt.Format("2006-01-02 15:04:05"),
        b.Duration.String(),
        b.HourlyRate,
        b.TotalAmount,
    )
}
