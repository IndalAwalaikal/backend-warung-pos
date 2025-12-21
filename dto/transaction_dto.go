package dto

type TransactionItemDTO struct {
	MenuID   uint    `json:"menu_id" binding:"required"`
	Quantity int     `json:"quantity" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
}

type TransactionCreateRequest struct {
	Items         []TransactionItemDTO `json:"items" binding:"required,dive,required"`
	Subtotal      float64               `json:"subtotal"`
	Tax           float64               `json:"tax"`
	Discount      float64               `json:"discount"`
	Total         float64               `json:"total"`
	PaymentMethod string                `json:"payment_method"`
	AmountPaid    float64               `json:"amount_paid"`
	CashierID     uint                  `json:"cashier_id"`
}

type TransactionResponse struct {
	ID    uint        `json:"id"`
	Items interface{} `json:"items"`
	Total float64     `json:"total"`
}
