package dto

type MenuCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	CategoryID  *uint   `json:"category_id" binding:"required"`
	ImageURL    string  `json:"image_url"`
	IsAvailable bool    `json:"is_available"`
}

type MenuResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  *uint   `json:"category_id"`
	ImageURL    string  `json:"image_url"`
	IsAvailable bool    `json:"is_available"`
}
