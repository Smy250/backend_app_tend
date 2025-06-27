package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `gorm:"primaryKey" json:"ID"` // Identificador del User
	Email     string         `gorm:"unique" json:"EMail"`  // Correo
	Username  string         `json:"Username"`             // Nombre de Usuario
	Password  string         `json:"Password"`             // Contraseña
	CreatedAt time.Time      // Tiempo en que se creo.
	UpdatedAt time.Time      // Tiempo que se actualizo.
	DeletedAt gorm.DeletedAt `gorm:"index"` // Eliminado o no (de forma logica)

	Consultas_AI []Consultas_AI `gorm:"foreignKey:User_ID"`
	// Un usuario esta ligado a varias Consultas de AI 1:M
	// Una consulta AI esta ligada a un usuario 1:1 ---> (1:M)
}

type Consultas_AI struct {
	ID         uint64 `gorm:"primaryKey"`
	User_ID    uint64 `json:"User_ID" gorm:"not null"`
	ConsultUID uint64 `json:"ConsultUID" gorm:"not null"`
	Consult    string `json:"Consult"`
	Request    string `json:"Request"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// Referencia al id del usuario
// Numero de consultas del usuario.
// Contenido de la consulta del usuario
// Contenido de la respuesta de Gemini
// Tiempo en que fue creado el registro
// Tiempo en que fue elminado el registro
// Eliminado o no (de forma logica)
