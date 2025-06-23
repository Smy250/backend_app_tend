package controllers

import (
	"net/http"

	"github.com/Smy250/backend_app_tend/config"
	"github.com/Smy250/backend_app_tend/models"
	"github.com/Smy250/backend_app_tend/scripts"
	"github.com/gin-gonic/gin"
)

func GetUserHistory(ctx *gin.Context) {
	db, err := config.DB_Instance()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var numberMessages []uint64 = scripts.FindConversationsNumbers(ctx, db)

	if len(numberMessages) == 0 {
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"Conversaciones": numberMessages})
	}
}

func GetUserHistoryID(ctx *gin.Context) {
	/* type historyMessages struct {
		consult string
		request string
	} */
	// messages := []historyMessages{}

	db, err := config.DB_Instance()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var usrMessages []models.ConsultaRespuesta = scripts.FindConversationHistoryID(ctx, db)

	if len(usrMessages) == 0 {
		ctx.AbortWithStatus(http.StatusOK)
	}

	ctx.JSON(http.StatusOK, gin.H{"Historial": usrMessages})
}
