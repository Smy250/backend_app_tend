package routes

import (
	"github.com/Smy250/backend_app_tend/controllers"
	"github.com/Smy250/backend_app_tend/middleware"
	"github.com/gin-gonic/gin"
)

func User_Routes(router_Group *gin.Engine) {
	router_Group.Group("/")
	{
		router_Group.POST("user/register", controllers.RegisterUser)
		router_Group.POST("user/login", middleware.UserVerifyLogging, controllers.LoginUser)
		router_Group.GET("user/logout", middleware.UserVerifyLogout, controllers.LogoutUser)
	}
}
