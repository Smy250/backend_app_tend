package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Smy250/backend_app_tend/apis"
	"github.com/Smy250/backend_app_tend/config"
	"github.com/Smy250/backend_app_tend/models"
	"github.com/Smy250/backend_app_tend/scripts"
	"google.golang.org/genai"

	"github.com/gin-gonic/gin"
)

func POST_Consult(ctx *gin.Context) {

	db, err := config.DB_Instance()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var usr_ID uint64 = scripts.FindUserID(ctx, db)
	if usr_ID == 0 {
		// Si el tipo de dato no es entero devolverá un error.
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar la consulta."})
		return
	}

	var consult = models.ConsultaGemini{}
	if err2 := ctx.ShouldBindJSON(&consult); err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}

	prompt := scripts.PromptPrecision(consult.Precision, consult.Consulta)

	// Declaramos una estructura que contiene la respuesta de Gemini.
	var response *genai.GenerateContentResponse
	var err3 error

	response, err3 = apis.ConsultGemini(os.Getenv("GEMINI_API_KEY"), prompt, usr_ID, consult.ConsultUID, consult.Precision)

	if response == nil ||
		len(response.Candidates) == 0 ||
		response.Candidates[0].Content == nil ||
		len(response.Candidates[0].Content.Parts) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "La IA no generó respuesta. Intenta con otra precisión o revisa tu consulta."})
		return
	}

	if err3 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Modelo de Inteligencia Artificial no Encontrado."})
		return
	}

	// Con respecto al dato usrID al ser de tipo any o cualquiera(generico)
	// Con la aserción podemos transformar un dato any a cualquiera
	// con .(tipo de dato)
	if err4 := db.Create(&models.Consultas_AI{User_ID: usr_ID, ConsultUID: consult.ConsultUID, Consult: consult.Consulta, Request: response.Text()}).Error; err4 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar la consulta"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Request": response.Text()})
}

// POST_SummarizePDF es el controlador que maneja la solicitud de resumen de PDF.
func POST_SummarizePDF(ctx *gin.Context) {

	var req models.SummarizePDFRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Formato de JSON inválido o campos faltantes: %v", err)})
		return
	}

	consulta := scripts.PromptPrecision(req.Precision, req.Consulta)

	// Llama a la nueva función del paquete apis para el resumen
	summary, err := apis.SummarizePDFAPI(os.Getenv("GEMINI_API_KEY"), req.FilePath, consulta, req.Precision)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al procesar el resumen del PDF: %v", err.Error())})
		return
	}

	db, err2 := config.DB_Instance()
	if err2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err2.Error()})
		return
	}

	var usr_ID uint64 = scripts.FindUserID(ctx, db)
	if usr_ID == 0 {
		// Si el tipo de dato no es entero devolverá un error.
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar la consulta."})
		return
	}

	if err3 := db.Create(&models.Consultas_AI{User_ID: usr_ID, ConsultUID: req.ConsultUID, Consult: req.Consulta, Request: summary}).Error; err3 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar la consulta"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Request": summary})
}
