package scripts

import (
	"strconv"

	"github.com/Smy250/backend_app_tend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Los siguientes scripts son usandos en la lógica del controlador.
// Relacionadas a encontrar el correo del usuario, id del usuario
// las conversaciones que ha tenido y la información de las conversaciones.

// Encontrar el usuario de acuerdo a un email.
func FindUserEmail(db *gorm.DB, formEmail string) bool {

	// flag servira para devolver un true o false en caso de que
	// la busqueda del correo en la base de datos sea exitosa, como no.
	var flag bool

	// checksID es de tipo entero uns, compuesta de 64 bits. Almacenará
	// el id resultante en la busqueda del correo, en la base de datos.
	var checksID uint64

	//La siguiente consulta es equivalente a:
	// SELECT id FROM users WHERE email = ? LIMIT 1, checksID recibirá el valor
	db.Table("users").Select("id").Where("email = ?", formEmail).Limit(1).Pluck("id", &checksID)
	if checksID == 0 {
		flag = false
	} else {
		flag = true
	}

	return flag
}

// Encontrar el usuario de acuerdo a un ID.
func FindUserID(ctx *gin.Context, db *gorm.DB) uint64 {
	var usr_ID uint64 = 0
	var value uint64 = 0

	if userID, ok := ctx.Get("user"); ok {
		if idUint, ok := userID.(uint64); ok {
			usr_ID = idUint
		}
	}

	db.Table("users").Select("id").Where("id = ?", usr_ID).Limit(1).Pluck("id", &value)

	return usr_ID
}

// Encontrar el numero de conversaciones del usuario.
func FindConversationsNumbers(ctx *gin.Context, db *gorm.DB) []uint64 {
	var usr_ID uint64 = 0
	var consult_values = make([]uint64, 0)

	if userID, ok := ctx.Get("user"); ok {
		if idUint, ok := userID.(uint64); ok {
			usr_ID = idUint
		}
	}

	db.Model(&models.Consultas_AI{}).
		Select("min(consult_uid) as consult_uid").
		Where("user_id = ?", usr_ID).
		Group("consult_uid").
		Scan(&consult_values)

	return consult_values
}

// Esta consulta a la base de datos devuelve toda la información de
// una conversación especifica de acuerdo al ID del usuario e id de la consulta.
func FindConversationHistoryID(ctx *gin.Context, db *gorm.DB) []models.ConsultaRespuesta {
	var usr_ID uint64 = 0
	var consult_values = []models.ConsultaRespuesta{}

	if userID, ok := ctx.Get("user"); ok {
		if idUint, ok := userID.(uint64); ok {
			usr_ID = idUint
		}
	}

	consultUID := ctx.Param("consult_uid")
	consultUID_N, err := strconv.ParseUint(consultUID, 10, 64)
	if err != nil {
		return consult_values
	}

	db.Model(models.Consultas_AI{}).
		Select("consult AS Consulta, request AS Respuesta").
		Where("user_id = ? AND consult_uid = ?", usr_ID, consultUID_N).
		Scan(&consult_values)

	return consult_values
}

func GetNextConversationID(ctx *gin.Context, db *gorm.DB) uint64 {
	var usr_ID uint64 = 0
	var conversationID uint64 = 0

	if userID, ok := ctx.Get("user"); ok {
		if idUint, ok := userID.(uint64); ok {
			usr_ID = idUint
		}
	}

	err := db.Model(models.Consultas_AI{}).
		Select("consult_uid AS conversationID").
		Where("user_id = ?", usr_ID).Last(&conversationID).Error

	if err != nil {
		return uint64(0)
	}

	return conversationID
}
