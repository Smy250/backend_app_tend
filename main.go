// Paquete
package main

import (
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
	// Llamada a la función para iniciar el .env y la base de datos.
	init_background()

	//Inicializamos el framework Gin para gestionar el servidor HTTP
	var gin_Router *gin.Engine = gin.Default()

	//Especificamos el Cors, que por ahora admitirá cualquier conexión
	gin_Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // tu frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
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

/*
	Cosas por Implementar:
	* 1. Login-Registro y Autenticar (Creación y Registro de Usuarios).
	[LISTO]
	* 2. Protección de rutas para evitar el acceso a información sensible.
	[LISTO]
	* 4. Modelar En Base A lo Que se Quiere del Proyecto [LISTO]
	* 5. IA de Gemini y sigan sus conversaciones para cada uno, sin que se mezcle el contexto de la conversación de un usuario con el de otro. [LISTO]
	* 8. Con lo anterior también evitar que se sobrepase del limite de consultas por minuto, que de acuerdo al plan que estable Google con su IA Gemini, es de 15.[EN PROGRESO]
	*9. Endpoint para devolver los historiales de conversación del usuario y toda la conversación específica selecionada. [LISTO]
*/
