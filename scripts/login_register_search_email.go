package scripts

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FindUserEmail(db *gorm.DB, formEmail string) bool {

	// flag servira para devolver un true o false en caso de que
	// la busqueda del correo en la base de datos sea exitosa, como no.
	var flag bool

	// checksID es de tipo entero uns, compuesta de 64 bits. Almacenará
	// el id resultante en la busqueda del correo, en la base de datos.
	var checksID uint64

	//La siguiente consulta es equivalente a:
	// SELECT id FROM users WHERE email = ? LIMIT 1, checksID recibirá valor
	db.Table("users").Select("id").Where("email = ?", formEmail).Limit(1).Pluck("id", &checksID)
	if checksID == 0 {
		flag = false
	} else {
		flag = true
	}

	return flag
}

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
