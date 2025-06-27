package routes

import (
	"github.com/Smy250/backend_app_tend/controllers"
	"github.com/Smy250/backend_app_tend/middleware"
	"github.com/gin-gonic/gin"
)

// Lo siguiente engloba todas las rutas que usará el servidor HTTP.
// Relacionadas a la lógica de recepción y envio de información a la IA de Gemini, como prueba sin verificación y middleware de verificación
func Route_Gemini(router_Group *gin.Engine) {
	router_Group.Group("/")
	{
		router_Group.POST("/consult/gemini", middleware.UserAuthentication, controllers.POST_Consult)
	}
}
