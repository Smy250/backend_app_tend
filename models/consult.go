package models

type Consultas_AI struct {
	ID      uint64 `gorm:"primaryKey"`
	Consult string
	Request string
}
