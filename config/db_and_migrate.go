package config

import (
	"errors"
	"os"

	"github.com/Smy250/backend_app_tend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Instancia de la la base de datos, para evitar llamarla globalmente
// y que ocurran casos raros.
func DB_Instance() (*gorm.DB, error) {
	var db_location = os.Getenv("DB_LOCATION")
	if db_location == "" {
		db_location = "./db/test.sqlite3"
	}

	var db *gorm.DB
	var err error

	db, err = gorm.Open(sqlite.Open(db_location), &gorm.Config{})
	if err != nil {
		return nil, errors.New("db error: no se pudo conectar correctamente con la base de de datos")
	}

	return db, err
}

// Chequeo y migración de la base de datos.
func Check_Migration() {
	db, err := DB_Instance()
	if err != nil {
		panic("Error al conectar la base de datos")
	}
	db.AutoMigrate(&models.User{}, &models.Consultas_AI{})
}
