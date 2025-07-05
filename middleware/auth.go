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
		ctx.Abort()
		return
	}

	token, err_2 := scripts.DecryptToken(tokenString)
	if err_2 != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Acceso no Autorizado": err_2.Error()})
		ctx.Abort()
		return
	}

	// Lo siguiente se refiere a verificar los datos del token previamente
	// desencriptado.
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//Se verifica que el token aún no se haya expirado.
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		var usr = models.User{}
		db, err_3 := config.DB_Instance()
		if err_3 != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error en la base de datos": err_3.Error()})
			ctx.Abort()
			return
		}

		// Encontramos el id del usuario con el atributo "usr" del token.
		db.First(&usr, claims["usr"])

		if usr.ID == 0 {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Adjuntamos al contexto del GIN ID y nombre de usuario, para ser
		// usado luego en la función a la cual requiere este middleware.
		ctx.Set("user", usr.ID)
		ctx.Set("username", usr.Username)

	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Acceso no Autorizado": "X"})
		ctx.Abort()
		return
	}
	// Con el metodo Next pasamos a la función pendiente, la cual necesita esta implementación para verficar que realmente el usuario esta autenticado
	ctx.Next()
}

func UserVerifyLogging(ctx *gin.Context) {
	// Verificamos si ya se ha logueado el usuario
	tokenString, _ := ctx.Cookie("Authorization")
	var token = &jwt.Token{}

	token, _ = scripts.DecryptToken(tokenString)

	// Lo siguiente se resume en que si el token del usuario previamente
	// desncriptado, no es nulo se procederá a verificar en la base de
	// datos el id. De ser un token nulo se procederá a loguearse y si no
	// le indicara que ya se encuentra logueado.
	if token != nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {

			var usr = &models.User{}
			db_conn, err := config.DB_Instance()
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error en la base de datos": err.Error()})
				return
			}

			//Encontramos el id del usuario con el atributo "usr" del token.
			db_conn.WithContext(context.Background()).First(&usr, claims["usr"])

			// si usr es distinto a nulo, devolvemos el mensaje de
			// que no hay necesidad de hacerlo nuevamente
			if usr.Login == 1 {
				ctx.AbortWithStatus(http.StatusConflict)
				return
			}
		}
	}

	ctx.Next()
}

func UserVerifyLogout(ctx *gin.Context) {
	// Obtenemos en un string, el apartado autohrization de la cookie.
	tokenString, _ := ctx.Cookie("Authorization")

	if tokenString != "" {
		token, err := scripts.DecryptToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Acceso no Autorizado": err.Error()})
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			var usr = &models.User{}
			db, err_2 := config.DB_Instance()
			if err_2 != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error en la base de datos": err_2.Error()})
				ctx.Abort()
				return
			}

			// Encontramos el id del usuario con el atributo "usr" del token.
			db.First(&usr, claims["usr"])

			db.Model(&models.User{}).
				Where("id = ?", usr.ID).
				Update("login", 0)
		}
		ctx.Next()
	} else {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Usted no tiene sesión activa."})
		return
	}
	// Lo anterior se resume en que si hay informacion en el campo cookie
	// de autorización se procederá a cerrar sesion. De lo contrario se
	// indicara que el usuario no tiene una sesión activa.
}
