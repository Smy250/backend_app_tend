package routes

import (
	"github.com/Smy250/backend_app_tend/controllers"
	"github.com/gin-gonic/gin"
)

func Route_Gemini(router_Group *gin.Engine) {
	gemini := router_Group.Group("/")
	{
		gemini.POST("/consult/gemini", controllers.POST_Consult)
	}
}
