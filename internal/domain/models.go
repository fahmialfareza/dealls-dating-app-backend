package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Name      string `json:"name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
	IsPremium bool   `json:"is_premium" gorm:"default:false"`

	// Profile
	Profile *Profile `json:"profile" gorm:"foreignKey:UserID"`
}

type Profile struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	UserID uint   `json:"user_id" gorm:"unique"`
	Bio    string `json:"bio"`
	Image  string `json:"image"`

	// Belongs to
	User *User `json:"user" gorm:"foreignKey:UserID"`
}

type Swipe struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	SwiperID uint      `json:"swiper_id" gorm:"index:swipe_unique_idx,unique"`
	SwipedID uint      `json:"swiped_id" gorm:"index:swipe_unique_idx,unique"`
	Type     string    `json:"type" gorm:"size:10;check:type IN ('like','pass')"`
	Date     time.Time `json:"date" gorm:"type:date;index:swipe_unique_idx,unique"`

	// Belongs to
	Swiper User `gorm:"foreignKey:SwiperID"`
	Swiped User `gorm:"foreignKey:SwipedID"`
}

type Purchase struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	UserID       uint      `json:"user_id" gorm:"unique"`
	PackageType  string    `json:"package_type" gorm:"size:10;check:package_type IN ('premium')"`
	PurchaseDate time.Time `json:"purchase_date"`

	// Belongs to
	User User `gorm:"foreignKey:UserID"`
}
