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
	init_background()

	var gin_Router *gin.Engine = gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://google.com"}
	gin_Router.Use(cors.New(config))

	routes.Route_Gemini(gin_Router)
	routes.Route_Middleware(gin_Router)
	routes.User_Routes(gin_Router)

	gin_Router.Run("0.0.0.0:8081")
}

/*
	Cosas por Implementar:
	* 1. Login-Registro y Autenticar (Creación y Registro de Usuarios).
	[LISTO]
	* 2. Protección de rutas para evitar el acceso a información sensible.
	[En Progreso]
	* 3. Arreglar la request o respuesta del json para evitar los "\n"
	en las respuestas. [En Progreso]
	* 4. Modelar En Base A lo Que se Quiere del Proyecto
	* 5. Posiblemente al escalar, organizar aún mejor. Apegandose mas al
	proyecto, la estructura de los directorios del proyecto.
	* 6. Ligado al punto 5, seguramente halla una reestructuración del código del proyecto para adaptarse a dicho requerimiento.
	* 7. Tratar y asegurarse de que cada usuario tenga su propio chat con la
	IA de Gemini y sigan sus conversaciones para cada uno, sin que se mezcle el contexto de la conversación de un usuario con el de otro.
	* 8. Con lo anterior también evitar que se sobrepase del limite de consultas por minuto, que de acuerdo al plan que estable Google con su IA Gemini, es de 15.
*/
