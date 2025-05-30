package config

import (
	"os"

	"github.com/Smy250/backend_app_tend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Load_DB() {
	var err error
	var db_location = os.Getenv("DB_LOCATION")
	if db_location == "" {
		db_location = "./db/test.sqlite3"
	}

	DB, err = gorm.Open(sqlite.Open(db_location), &gorm.Config{})
	if err != nil {
		panic("Error al conectar la base de datos")
	}
}

func Check_Migration() {
	Load_DB()
	DB.AutoMigrate(&models.User{}, &models.Consultas_AI{})
}
