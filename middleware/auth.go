package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Smy250/backend_app_tend/config"
	"github.com/Smy250/backend_app_tend/models"
	"github.com/Smy250/backend_app_tend/scripts"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func UserAuthentication(ctx *gin.Context) {
	tokenString, err := ctx.Cookie("Authorization")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Acceso no Autorizado": "Token invalido"})
		return
		// OJO: Si err es nulo y llamamos el metodo .Error(), invocará
		// un panic debido a que .Error() hace referencia algo y esta nulo
	} else if tokenString == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"ERROR": "Acceso no autorizado."})
		return
	}

	token, err_2 := scripts.DecryptToken(tokenString)
	if err_2 != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Acceso no Autorizado": err_2.Error()})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//Se verifica que el token aún no se haya expirado.
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		var usr = models.User{}
		db_conn, err_3 := config.DB_Instance()
		if err_3 != nil {
			panic("Error al conectar la base de datos")
		}

		//Encontramos el id del usuario con el atributo "usr" del token.
		db_conn.WithContext(context.Background()).First(&usr, claims["usr"])

		if usr.ID == 0 {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		//Adjuntamos al contexto.
		ctx.Set("user", usr.ID)
		ctx.Set("username", usr.Username)

	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Acceso no Autorizado": "X"})
	}

	ctx.Next()
}
