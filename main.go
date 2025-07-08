// Paquete
package main

import (
	"strings"

	"github.com/Smy250/backend_app_tend/config"
	"github.com/Smy250/backend_app_tend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Inicialización de:
// (Variables de Entorno, Inicialización y Chequeo en la Base de Datos,
// Migración de Tablas a BD en Caso de no tenerlas creadas.
func init_background() {
	config.LoadEnv()
	config.Check_Migration()
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	// Llamada a la función para iniciar el .env y la base de datos.
	init_background()

	//Inicializamos el framework Gin para gestionar el servidor HTTP
	var gin_Router *gin.Engine = gin.Default()

	//Especificamos el Cors, que por ahora admitirá cualquier conexión
	gin_Router.Use(cors.New(cors.Config{ // tu frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://") && strings.Contains(origin, ":5173")
		},
	}))

	//Especificamos las rutas que se usarán en el servidor.
	// Contiene las rutas para la consulta de Gemini, Middlew y usuario
	routes.Route_Gemini(gin_Router)
	routes.Route_Middleware(gin_Router)
	routes.User_Routes(gin_Router)

	// Le especificamos que corra el servidor y este atento
	// a escuchar a cualquier dirección en el puerto 8081.
	gin_Router.Run(":8081")
}
