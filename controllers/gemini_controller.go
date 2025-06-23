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
	"github.com/Smy250/backend_app_tend/scripts"
	"google.golang.org/genai"

	"github.com/gin-gonic/gin"
)

func POST_Consult(ctx *gin.Context) {

	//Obtenemos el id y username desde el middleware.
	var usr_username = ""

	if username, ok := ctx.Get("username"); ok {
		usr_username = username.(string)
	}

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

	// Obtenemos los datos de la respuesta del usuario recibido.
	var jsonData map[string]any
	if err2 := ctx.BindJSON(&jsonData); err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}

	// Formulamos lo recibido del JSON a variables.
	consult := jsonData["Consulta"].(string)

	nrConsult := jsonData["ConsultUID"].(string)
	nrConsultInt, _ := strconv.ParseUint(nrConsult, 10, 64)

	modelo_tx := jsonData["Modelo"].(string)
	modelo, _ := strconv.ParseUint(modelo_tx, 10, 64)

	// Consulta del usuario.
	consult = fmt.Sprintf("\"%s\" (Consulta del Usuario: %s)", consult, usr_username)

	// Declaramos una estructura que contiene la respuesta de Gemini.
	var response *genai.GenerateContentResponse
	var err3 error

	response, err3 = apis.ConsultGemini(os.Getenv("GEMINI_API_KEY"), consult, usr_ID, nrConsultInt, modelo, 1)

	if err3 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Modelo de Inteligencia Artificial no Encontrado."})
		return
	}

	//fmt.Printf("%v - %v - %v", usr_ID, nrConsultInt, consult)

	// Con respecto al dato usrID al ser de tipo any o cualquiera(generico)
	// Con la aserción podemos transformar un dato any a cualquiera
	// con .(tipo de dato)
	if err4 := db.Create(&models.Consultas_AI{User_ID: usr_ID, ConsultUID: nrConsultInt, Consult: consult, Request: response.Text()}).Error; err4 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar la consulta"})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err4.Error()})
		return
	}

	// Lo siguiente fue para mostrar visualmente agradable en consola la
	// respuesta, será eliminado posteriormente.
	var test string = strings.ReplaceAll(response.Text(), "\n", " ")

	ctx.JSON(http.StatusOK, gin.H{"Request": test})
}

func POST_Consult_NoAuth(ctx *gin.Context) {
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

	modelo_tx := jsonData["Modelo"].(string)
	modelo, _ := strconv.ParseUint(modelo_tx, 10, 64)

	consult = fmt.Sprintf("`%s` (Consulta del Usuario: Invitado)", consult)

	var response *genai.GenerateContentResponse
	var err2 error

	response, err2 = apis.ConsultGemini(os.Getenv("GEMINI_API_KEY"), consult, 0, nrConsultInt, modelo, 2)

	if err2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Modelo de Inteligencia Artificial no Encontrado."})
		return
	}

	db, err3 := config.DB_Instance()
	if err3 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err3.Error()})
		return
	}

	// Con respecto al dato usrID al ser de tipo any o cualquiera(generico)
	// Con la aserción podemos transformar un dato any a cualquiera
	// con .(tipo de dato)
	if err4 := db.Create(&models.Consultas_AI{User_ID: 0, Consult: consult, ConsultUID: nrConsultInt, Request: response.Text()}).Error; err4 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar la consulta"})
		return
	}

	//var test string = strings.ReplaceAll(response.Text(), "\n", " ")

	ctx.JSON(http.StatusOK, gin.H{"Request": response.Text()})
}
