package models

import (
	"time"

	"gorm.io/gorm"
)

/*
	Estructura del modelo usuario

Lo que ves al lado del tipo de dato en vuelto en 'gorm...' es la forma
de especificar que por ejemplo este dato puede ser representado
tanto en json, como un atributo clave en el orm de gorm y ser
mapeado en este como una tabla. En el caso de ID en gorm es
especifica que es una primary key y como json sera
identificado como el campo ID.
*/
type User struct {
	// Identificador del usuario.
	ID uint64 `gorm:"primaryKey" json:"ID"`
	// Correo.
	Email string `gorm:"unique" json:"EMail"`
	// Nombre del usuario.
	Username string `json:"Username"`
	// Contraseña
	Password string `json:"Password"`
	// Estado del login
	Login uint8
	// Tiempo en que se creo.
	CreatedAt time.Time
	// Tiempo que se actualizo.
	UpdatedAt time.Time
	// Borrado lógico en caso de aplicarse.
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Consultas_AI []Consultas_AI `gorm:"foreignKey:User_ID"`
	// Un usuario esta ligado a varias Consultas de AI 1:M
	// Una consulta AI esta ligada a un usuario 1:1 ---> (1:M)
}

// Estructura del modelo Consultas_AI, tanto JSON, como tabla en GORM
type Consultas_AI struct {
	// Identificador de la consulta.
	ID uint64 `gorm:"primaryKey"`
	// Identificador del usuario.
	User_ID uint64 `json:"User_ID" gorm:"not null"`
	// N° de identifador de conversaciones.
	ConsultUID uint64 `json:"ConsultUID" gorm:"not null"`
	// Contenido de la consulta del usuario.
	Consult string `json:"Consult"`
	// Contenido de la respuesta de Gemini.
	Request string `json:"Request"`
	// Tiempo en que fue creado el registro.
	CreatedAt time.Time
	// Tiempo en que fue elminado el registro.
	UpdatedAt time.Time
	// Borrado lógico en caso de aplicarse.
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Estructura de la consulta en Gemini. Contiene el texto de la consulta
// ID de la consulta del usuario y la precision o mejor dicho el tipo de
// respuesta que el usuario necesita el usuario.
type ConsultaGemini struct {
	Consulta   string `json:"Consulta"`
	ConsultUID uint64 `json:"ConsultUID" binding:"required"`
	Precision  uint64 `json:"Precision"`
}

// Estructura ConsultaRespuesta. Solo contendrá la información de la
// consulta del usuario y respuesta Gemini.
type ConsultaRespuesta struct {
	Consulta  string `json:"Consulta"`
	Respuesta string `json:"Respuesta"`
}

// SummarizePDFRequest define el formato del JSON de entrada para este endpoint.
type SummarizePDFRequest struct {
	ConsultaGemini
	// Anidamiento de los campos del struct ConsultaGemini.
	// Es decir, 'Consulta, ConsultUID, Precision' conforman el struct.
	FilePath string `json:"FilePath" binding:"required"` // Ruta al archivo PDF en el servidor
}
