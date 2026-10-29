package model

type Book struct {
	ID         uint     `json:"id" gorm:"primaryKey"`
	Title      string   `json:"title"`
	ISBN       string   `json:"isbn"`
	CategoryID uint     `json:"category_id"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`
	Authors    []Author `json:"authors" gorm:"many2many:book_authors;"`
	Copies     []Copy   `json:"copies,omitempty" gorm:"foreignKey:BookID"`
}

type BookInput struct {
	Title      string `json:"title"`
	ISBN       string `json:"isbn"`
	CategoryID uint   `json:"category_id"`
	AuthorIDs  []uint `json:"author_ids"`
}
