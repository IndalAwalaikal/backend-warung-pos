package model

type TransactionItem struct {
    ID            uint    `gorm:"primaryKey" json:"id"`
    TransactionID uint    `json:"transaction_id"`
    MenuID        *uint   `json:"menu_id"`
    Quantity      int     `json:"quantity"`
    Price         float64 `json:"price"`
    Menu          Menu    `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
}
