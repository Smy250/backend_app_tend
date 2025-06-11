package apis

import (
	"context"
	"errors"

	"github.com/Smy250/backend_app_tend/config"
	"github.com/Smy250/backend_app_tend/models"
	"google.golang.org/genai"
)

func ConsultGemini(gemini_key string, consult string,
	user_id uint64, nConsult uint64, model uint64) (*genai.GenerateContentResponse, error) {

	if model > 3 {
		return nil, errors.New("error en el modelo: Solo puede elegir tres modelos, [3 - Avanzado, 2 - Intermedio, 1 - Basico]")
	}

	// Por los momentos se harán pruebas con los siguientes tres modelos
	// que ofrece la IA Gemini de Google
	var ai_Models = []string{
		"gemini-2.0-flash-lite-001",
		"gemini-2.0-flash",
		"gemini-2.5-flash-preview-04-17",
	}

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
	// Si hay un error al crear un cliente, se verificará

	// La siguiente variable se refiere a los parametros de generación
	// de contenido. en ellos esta temperatura, maximos de token por resp.
	//var config *genai.GenerateContentConfig = temperatureParameter(0, 0, 0, 0)
	chat, err_2 := client.Chats.Create(context.Background(), ai_Models[model], &genai.GenerateContentConfig{ResponseMIMEType: "application/json"}, nil)
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
