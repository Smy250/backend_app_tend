package routes

import (
	"github.com/Smy250/backend_app_tend/controllers"
	"github.com/Smy250/backend_app_tend/middleware"
	"github.com/gin-gonic/gin"
)

// Lo siguiente engloba todas las rutas que usará el servidor HTTP.
// Relacionadas a la lógica del usuario.
func User_Routes(router_Group *gin.Engine) {
	router_Group.Group("/")
	{
		router_Group.POST("user/register", controllers.RegisterUser)
		router_Group.POST("user/login", middleware.UserVerifyLogging, controllers.LoginUser)
		router_Group.GET("user/logout", middleware.UserVerifyLogout, controllers.LogoutUser)
		router_Group.GET("user/number_conversations", middleware.UserAuthentication, controllers.GetUserHistory)
		router_Group.GET("user/conversation/:consult_uid", middleware.UserAuthentication, controllers.GetUserHistoryID)
		router_Group.GET("user/new_conversation", middleware.UserAuthentication, controllers.NewConversationUser)
	}
}
