package routes

import (
	"net/http"

	"github.com/Smy250/backend_app_tend/middleware"
	"github.com/gin-gonic/gin"
)

func Route_Middleware(router_Group *gin.Engine) {
	middleware_route := router_Group.Group("/")
	{
		middleware_route.GET("/test_2", middleware.UserAuthentication, func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"X": "OK"})
		})
	}
}
