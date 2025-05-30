package routes

import (
	"github.com/Smy250/backend_app_tend/controllers"
	"github.com/gin-gonic/gin"
)

func User_Routes(router_Group *gin.Engine) {
	router_Group.Group("/")
	{
		router_Group.POST("user/register", controllers.RegisterUser)
		router_Group.POST("user/login", controllers.LoginUser)
		router_Group.GET("user/logout", controllers.LogoutUser)
	}
}
