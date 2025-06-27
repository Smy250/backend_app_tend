package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("Please create a .env file in the root directory of the project\n", "Es necesario la creacion de un archivo .env en su directorio, agreguelo a su directorio raiz e inicie nuevamente.")
	}
}
