package config

import (
	"errors"
	"os"

	"github.com/Smy250/backend_app_tend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DB_Instance() (*gorm.DB, error) {
	var db_location = os.Getenv("DB_LOCATION")
	if db_location == "" {
		db_location = "./db/test.sqlite3"
	}

	db, err := gorm.Open(sqlite.Open(db_location), &gorm.Config{})
	if err != nil {
		return nil, errors.New("db error: no se pudo conectar correctamente con la base de de datos")
	}

	return db, err
}

func Check_Migration() {
	db, err := DB_Instance()
	if err != nil {
		panic("Error al conectar la base de datos")
	}
	db.AutoMigrate(&models.User{}, &models.Consultas_AI{})
}
