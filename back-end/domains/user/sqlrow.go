package user

import (
	"time"

	"github.com/google/uuid"
)

type InvalidToken struct {
	Token   string
	Expires time.Time
}

func (InvalidToken) TableName() string {
	return "invalid_token"
}

type Admin struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string
	Password  string
	Email     string `gorm:"unique"`
	AvatarUrl string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Admin) TableName() string {
	return "admin"
}

type Customer struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string
	Password    string
	Email       string `gorm:"unique"`
	AvatarUrl   string
	GoogleID    string    `gorm:"uniqueIndex"`
	HasPassword bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (Customer) TableName() string {
	return "customer"
}
