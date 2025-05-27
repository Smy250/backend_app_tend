package controllers

import (
	"net/http"
	"os"

	"github.com/Smy250/backend_app_tend/apis"
	"github.com/Smy250/backend_app_tend/config"
	"github.com/Smy250/backend_app_tend/models"
	"google.golang.org/genai"

	"github.com/gin-gonic/gin"
)

func POST_Consult(c *gin.Context) {
	db := config.DB

	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener un atributo específico
	consult := jsonData["Consulta"].(string)

	var response *genai.GenerateContentResponse
	var err_2 error
	response, err_2 = apis.ConsultGemini(os.Getenv("GEMINI_API_KEY"), consult)

	if err_2 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Modelo de Inteligencia Artificial no Encontrado."})
		return
	}

	if err_3 := db.Create(&models.Consultas_AI{Consult: consult, Request: response.Text()}).Error; err_3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al insertar el usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Request": response.Text()})
}
