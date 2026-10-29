package model

type Copy struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	BookID      uint   `json:"book_id"`
	InventoryNo string `json:"inventory_no"`
	Status      string `json:"status"`
}

type Availability struct {
	BookID          uint `json:"book_id"`
	TotalCopies     int  `json:"total_copies"`
	AvailableCopies int  `json:"available_copies"`
}
