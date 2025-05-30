package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey" json:"ID"`
	Email     string `gorm:"unique" json:"EMail"`
	Username  string `json:"Username"`
	Password  string `json:"Password"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Consultas_AI []Consultas_AI `gorm:"foreignKey:User_ID"`
	// Un usuario esta ligado a varias Consultas de AI 1:M
	// Una consulta AI esta ligada a un usuario 1:1 ---> (1:M)
}

type Consultas_AI struct {
	ID        uint64 `gorm:"primaryKey"`
	User_ID   uint64 `json:"User_ID"`
	Consult   string `json:"Consult"`
	Request   string `json:"Request"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
