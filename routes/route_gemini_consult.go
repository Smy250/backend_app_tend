package routes

import (
	"github.com/Smy250/backend_app_tend/controllers"
	"github.com/Smy250/backend_app_tend/middleware"
	"github.com/gin-gonic/gin"
)

func Route_Gemini(router_Group *gin.Engine) {
	router_Group.Group("/")
	{
		router_Group.POST("/consult/gemini", middleware.UserAuthentication, controllers.POST_Consult)
	}
}
