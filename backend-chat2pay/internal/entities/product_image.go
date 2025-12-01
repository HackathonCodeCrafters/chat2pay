package entities

import "time"

type ProductImage struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint64    `gorm:"not null;index:idx_product_images_product_id" json:"product_id"`
	ImageURL  string    `gorm:"type:text;not null" json:"image_url"`
	IsPrimary bool      `gorm:"type:boolean;not null;default:false" json:"is_primary"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
