package model

import "time"

type Transaction struct {
    ID        uint              `gorm:"primaryKey" json:"id"`
    Total       float64           `json:"total"`
    Subtotal    float64           `json:"subtotal"`
    Tax         float64           `json:"tax"`
    Discount    float64           `json:"discount"`
    PaymentMethod string          `json:"payment_method"`
    AmountPaid  float64           `json:"amount_paid"`
    CashierID   *uint             `json:"cashier_id"`
    Items       []TransactionItem `gorm:"foreignKey:TransactionID" json:"items"`
    CreatedAt   time.Time         `json:"created_at"`
}
