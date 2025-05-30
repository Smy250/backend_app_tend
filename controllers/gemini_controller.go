package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Smy250/backend_app_tend/apis"
	"github.com/Smy250/backend_app_tend/config"
	"github.com/Smy250/backend_app_tend/models"
	"google.golang.org/genai"

	"github.com/gin-gonic/gin"
)

func POST_Consult(ctx *gin.Context) {
	db := config.DB
	usr, _ := ctx.Get("user")
	usr_string := fmt.Sprint(usr)
	usr_ID, _ := strconv.ParseUint(usr_string, 0, 64)

	fmt.Println(usr_ID)

	var jsonData map[string]interface{}
	if err := ctx.BindJSON(&jsonData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener un atributo específico
	consult := jsonData["Consulta"].(string)

	var response *genai.GenerateContentResponse
	var err_2 error
	response, err_2 = apis.ConsultGemini(os.Getenv("GEMINI_API_KEY"), consult)

	if err_2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Modelo de Inteligencia Artificial no Encontrado."})
		return
	}

	// Con respecto al dato usrID al ser de tipo any o cualquiera(generico)
	// Con la aserción podemos transformar un dato any a cualquiera
	// con .(tipo de dato)
	if err_3 := db.Create(&models.Consultas_AI{User_ID: usr_ID, Consult: consult, Request: response.Text()}).Error; err_3 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al insertar el usuario"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Request": response.Text()})
}
