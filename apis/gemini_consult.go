package apis

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Smy250/backend_app_tend/config"
	"github.com/Smy250/backend_app_tend/models"
	"google.golang.org/genai"
)

const GeminiModel = "gemini-2.5-flash"

func ConsultGemini(gemini_key string, consult string,
	user_id uint64, nConsult uint64, precision uint64) (*genai.GenerateContentResponse, error) {

	// Por los momentos se harán pruebas con los siguientes tres modelos
	// que ofrece la IA Gemini de Google

	// Declararemos un contexto con el fin, para que tome el tiempo
	// necesario, en procesar la información recibida de la API de Gemini.
	//var ctx = context.Background()
	client, err1 := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  gemini_key,
		Backend: genai.BackendGeminiAPI,
	})
	if err1 != nil {
		return nil, err1
	}
	// Si hay un error con la API de Gemini, devolverá el error obtenido.

	// Obtenemos la configuración de la precision de acuerdo a un numero
	// especificado por el usuario (1:Tutor - 2: Investigativo).
	precisionM := precisionModel(precision)
	if precisionM == nil {
		return nil, errors.New("configuración de precisión inválida")
	}

	// La siguiente variable adjunta toda la conversación del usuario
	// almacenada en la base de datos.
	history, err_2 := getRecentMessage(user_id, nConsult, 2)
	if err_2 != nil {
		return nil, err_2
	} // Obtenemos de la bd al menos por ahora 2 mensajes previos para
	// el contexto de la conversación, se puede anexar mas pero lo ideal
	// es no sobrecargar tanto de información la consulta a la API.

	// Creamos el chat con el modelo y precisión previamente obtenidos.
	chat, err_3 := client.Chats.Create(context.Background(), GeminiModel, precisionM, history)
	if err_3 != nil {
		return nil, err_3
	}

	// Se envia la información previa.
	resp, _ := chat.SendMessage(context.Background(), genai.Part{Text: consult})

	return resp, nil
}

func getRecentMessage(user uint64, nConsult uint64, limit int) ([]*genai.Content, error) {
	//Declaramos una variable con la configuración de la BD. para hacer operaciones en ella.
	db, err4 := config.DB_Instance()
	if err4 != nil {
		return []*genai.Content{}, errors.New("db error: no se pudo conectar correctamente con la base de de datos")
	}

	var usr []models.Consultas_AI

	// Importante sin un timeout o context es probable que se ligue
	// las consultas con las respuestas almacenadas en la BD.
	// Confundiendo a la IA.
	db.WithContext(context.Background()).Table("consultas_ais").
		Where("user_id = ? AND consult_uid = ?", user, nConsult).
		Order("rowid DESC").Limit(limit).Find(&usr)

	history := []*genai.Content{}

	for _, elem := range usr {
		history = append([]*genai.Content{
			genai.NewContentFromText(elem.Consult, genai.RoleUser),
			genai.NewContentFromText(elem.Request, genai.RoleModel),
		}, history...)
	}

	return history, nil
}

// Esta funcion devuelve una referencia de tipo GenerateContentConfig la
// cual contiene parametros de temperatura, toP, topK y el máximo de
// tokens que podrá enviar Gemini. dependiendo del nivel de precisión 1
// o 2, se configuraran los parametros para adaptarse lo mas cercano
// posible a esa funcionalidad.
func precisionModel(precisionlvl uint64) *genai.GenerateContentConfig {

	var contentConfig *genai.GenerateContentConfig
	var (
		temp, topP, topK float32
	)

	maxOutputTokens := int32(65536)

	switch precisionlvl {
	// Explicativo/Ejemplos
	case 1:
		temp = 0.65
		topP = 0.78
		topK = 0
	// Investigativo
	case 2:
		temp = 0.3
		topP = 0
		topK = 20
	// Guía/Ejemplos/General
	default:
		temp = 0.79
		topP = 0.95
		topK = 0
	}

	contentConfig = &genai.GenerateContentConfig{
		Temperature:     &temp,
		TopP:            &topP,
		TopK:            &topK,
		MaxOutputTokens: maxOutputTokens,
	}

	return contentConfig
}

// Función para manejar la solicitud de resumen de PDF
func SummarizePDFAPI(gemini_key string, filePath string, prompt string, precision uint64) (string, error) {
	// Creamos un cliente de la API GenAI, especificandole APIKEY
	client, err1 := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  gemini_key,
		Backend: genai.BackendGeminiAPI,
	})
	if err1 != nil {
		return "nil", err1
	}

	//Leemos el el archivo pdf con la función ReadFile
	pdfBytes, err_file := os.ReadFile(filePath)
	if err_file != nil {
		return "", fmt.Errorf("error al leer el archivo PDF en '%s': %w", filePath, err_file)
	}

	if prompt == "" {
		prompt = "Summarize this document in detail. Provide key points and main ideas." // Prompt por defecto si no se proporciona uno.
	}

	// Cargamos los bytes obtenidos del PDF a una estructura parts de la API GenAI, con ello podremos adjuntarle a la consulta y compactar la consulta y PDF en esta estructura.
	parts := []*genai.Part{
		{
			InlineData: &genai.Blob{
				MIMEType: "application/pdf",
				Data:     pdfBytes,
			},
		},
		genai.NewPartFromText(prompt),
	}

	// La siguiente estructura complementa todo el contenido del struct parts y especifica se refiere a información de la estructura parts y dicha informacion se aclara que es del usuario(RoleUser).
	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	// Declaramos y almacenamos en esta variable los parametros de acuerdo
	// a la precisión que estableció el usuario.
	precisionM := precisionModel(precision)
	if precisionM == nil {
		return "", errors.New("configuración de precisión inválida")
	}

	// Creamos el chat con tipo de modelo de gemini (en este caso se trata del 2.5) y precisión previamente obtenidos.
	chat, err_3 := client.Chats.Create(context.Background(), GeminiModel, precisionM, nil)
	if err_3 != nil {
		return "", err_3
	}

	//Llamada a la API de Gemini para consultar a la IA
	resp, err_genai := chat.Models.GenerateContent(
		context.Background(),
		GeminiModel,
		contents,
		nil,
	)
	if err_genai != nil {
		return "", fmt.Errorf("error de la API de Gemini al generar contenido: %w", err_genai)
	}

	// Extraer y enviar la respuesta de la IA
	var summary string
	if resp != nil && len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			summary += part.Text
			//Acceso directo al campo .text de la estructura
		}
	} else {
		return "", errors.New("la IA no proporcionó una respuesta válida o esperada")
	}

	return summary, nil
}
