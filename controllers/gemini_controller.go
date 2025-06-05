package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Smy250/backend_app_tend/apis"
	"github.com/Smy250/backend_app_tend/config"
	"github.com/Smy250/backend_app_tend/models"
	"google.golang.org/genai"

	"github.com/gin-gonic/gin"
)

func POST_Consult(ctx *gin.Context) {

	//Obtenemos el id y username desde el middleware.
	var usr_username = ""
	var usr_ID uint64

	if userID, ok := ctx.Get("user"); ok {
		if idUint, ok := userID.(uint64); ok {
			usr_ID = idUint
		} else {
			// Si el tipo de dato no es uint64
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar la consulta."})
			return
		}
	}

	if username, ok := ctx.Get("username"); ok {
		usr_username = username.(string)
	}
	// Refactorizar

	// Obtenemos los datos de la respuesta del usuario recibido.
	var jsonData map[string]any
	if err := ctx.BindJSON(&jsonData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener un atributo específico
	consult := jsonData["Consulta"].(string)
	nrConsult := jsonData["ConsultUID"].(string)
	nrConsultInt, _ := strconv.ParseUint(nrConsult, 10, 64)

	consult = fmt.Sprintf("`%s` (Consulta del Usuario:%s)", consult, usr_username)

	var response *genai.GenerateContentResponse
	var err2 error

	response, err2 = apis.ConsultGemini(os.Getenv("GEMINI_API_KEY"), consult, usr_ID, nrConsultInt)

	if err2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Modelo de Inteligencia Artificial no Encontrado."})
		return
	}

	db, err4 := config.DB_Instance()
	if err4 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err4.Error()})
		return
	}

	// Con respecto al dato usrID al ser de tipo any o cualquiera(generico)
	// Con la aserción podemos transformar un dato any a cualquiera
	// con .(tipo de dato)
	if err4 = db.Create(&models.Consultas_AI{User_ID: usr_ID, Consult: consult, ConsultUID: nrConsultInt, Request: response.Text()}).Error; err4 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar la consulta"})
		return
	}

	var test string = strings.ReplaceAll(response.Text(), "\n", " ")

	ctx.JSON(http.StatusOK, gin.H{"Request": test})
}
