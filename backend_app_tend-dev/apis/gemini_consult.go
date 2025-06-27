package apis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Smy250/backend_app_tend/config"
	"github.com/Smy250/backend_app_tend/models"
	"google.golang.org/genai"
)

// Una instancia global del cliente de Gemini para incializarlo
// una sola vez
var aiClient *genai.Client

// Lista de modelos de IA disponibles
var aiModels = []string{
	"gemini-2.0-flash-lite-001",
	"gemini-2.0-flash",
	"gemini-2.5-flash-preview-04-17",
}

// Función para inicializar el cliente de la API de Gemini
func InitGeminiClient(ctx context.Context, apiKey string) error {
	var err error
	aiClient, err = genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return fmt.Errorf(" error fatal: No se pudo inicializar el cliente de Gemini. Verifique su API Key. Error: %w", err)

	}
	log.Println("Cliente de Gemini incializado correctamente.")
	return nil
}

// Consulta con la Inteligencia Artificial (IA) de Gemini
// Se usa el cliente como global.
func ConsultGemini(gemini_key string, consult string,
	user_id uint64, nConsult uint64, model uint64) (*genai.GenerateContentResponse, error) {

	if aiClient == nil {
		var err error
		aiClient, err = genai.NewClient(context.Background(), &genai.ClientConfig{
			APIKey:  gemini_key, // Usa la clave proporcionada si el cliente global no existe
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			return nil, fmt.Errorf("error al crear cliente de Gemini en ConsultGemini: %w", err)
		}
	}
	if model >= uint64(len(aiModels)) {
		return nil, errors.New("error en el modelo: Solo puede elegir tres modelos, [3 - Avanzado, 2 - Intermedio, 1 - Basico]")
	}

	chat, err_2 := aiClient.Chats.Create(context.Background(), aiModels[model], &genai.GenerateContentConfig{ResponseMIMEType: "application/json"}, nil)
	if err_2 != nil {
		return nil, err_2
	}

	history, err_3 := getRecentMessage(user_id, int(nConsult), 2)
	if err_3 != nil {
		return nil, err_3
	}

	history = append(history, genai.Part{Text: consult})

	resp, _ := chat.SendMessage(context.Background(), history...)

	return resp, nil
}

func getRecentMessage(user uint64, nConsult, limit int) ([]genai.Part, error) {
	db, err := config.DB_Instance()
	if err != nil {
		return []genai.Part{}, errors.New("db error: no se pudo conectar correctamente con la base de de datos")
	}

	var usr []models.Consultas_AI

	// Importante sin un timeout o context es probable que se ligue
	// las consultas con las respuestas almacenadas en la BD.
	// Confundiendo a la IA.
	db.WithContext(context.Background()).Table("consultas_ais").
		Where("user_id = ? AND consult_uid = ?", user, nConsult).
		Order("rowid DESC").Limit(limit).Find(&usr)

	history := []genai.Part{}

	for _, elem := range usr {
		history = append([]genai.Part{{Text: elem.Consult}}, history...)
		history = append([]genai.Part{{Text: elem.Request}}, history...)
	}

	return history, nil
}

// Estructura que define el formato del JSON de entrada para este endpoint
type SummarizePDFRequest struct {
	FilePath string `json:"file_path" binding:"required"`
	Prompt   string `json:"prompt"`
	ModelID  uint64 `json:"model_id"`
}

// Función para manejar la solicitud de resumen de PDF
func SummarizePDFAPI(filePath string, prompt string, modelIndex uint64) (string, error) {
	if aiClient == nil {
		return "", errors.New("cliente de IA no inicializado en el paquete 'apis'")
	}
	//Validacion del indice del modelo
	modelToUse := "gemini-2.5-flash-preview-04-17"
	if modelIndex < uint64(len(aiModels)) {
		modelToUse = aiModels[modelIndex]
	} else {
		log.Printf("Advertencia: El model_id %d está fuera de rango. Usando modelo por defecto: %s", modelIndex, modelToUse)
	}
	//Leer el archivo pdf
	pdfBytes, err_file := os.ReadFile(filePath)
	if err_file != nil {
		return "", fmt.Errorf("error al leer el archivo PDF en '%s': %w", filePath, err_file)
	}

	if prompt == "" {
		prompt = "Summarize this document in detail. Provide key points and main ideas." // Prompt por defecto si no se proporciona uno.
	}
	/*Slice de punteros
	parts := []*genai.Part{
		&genai.Part{
			InlineData: &genai.Blob{
				MIMEType: "application/pdf",
				Data:     pdfBytes,
			},
		},
		genai.NewPartFromText(promptText),
	}*/

	// Usando la sintaxis de la documentación (slice de punteros)
	parts := []*genai.Part{
		{
			InlineData: &genai.Blob{
				MIMEType: "application/pdf",
				Data:     pdfBytes,
			},
		},
		genai.NewPartFromText(prompt),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	//Llamada al modelo de Gemini para generar el resumen
	resp, err_genai := aiClient.Models.GenerateContent(
		context.Background(),
		modelToUse,
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
			//summary += fmt.Sprintf("%v", part)
			summary += part.Text //Acceso directo al campo .text de la estructura
		}
	} else {
		return "", errors.New("la IA no proporcionó una respuesta válida o esperada")
	}

	return summary, nil
}

/*
func temperatureParameter(temp float32, topP float32, topK float32, maxOutputTokens int32) *genai.GenerateContentConfig {

	//Parametro por defecto si todos los parametros estan en cero.
	if temp == 0 || topP == 0 || topK == 0 || maxOutputTokens == 0 {
		temp = 0.9
		topP = 0.5
		topK = 20.0
		maxOutputTokens = 100
	}
	// Devolvemos un struct de tipo genai.GenerateContentConfig:
	return &genai.GenerateContentConfig{
		Temperature:      &temp,
		TopP:             &topP,
		TopK:             &topK,
		MaxOutputTokens:  maxOutputTokens,
		ResponseMIMEType: "application/json",
	}
}*/
