package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Smy250/backend_app_tend/config"
	"github.com/Smy250/backend_app_tend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx *gin.Context) {
	var userForm = models.User{} // Declaración del Modelo User
	ctx.BindJSON(&userForm)

	if (userForm.Email == "") ||
		(userForm.Password == "") || (userForm.Username == "") {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Campos faltantes"})
		return
	}

	db := config.DB

	encriptedPassByte, err1 := bcrypt.GenerateFromPassword([]byte(userForm.Password), 14)
	if err1 != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"Error:": "Ha ocurrido un error al generar la contraseña del usuario"})
	}

	userForm.Password = string(encriptedPassByte)

	db.Create(&userForm)

	ctx.JSON(http.StatusOK, gin.H{"El usuario se ha registrado correctamente:": userForm})

}

func LoginUser(ctx *gin.Context) {
	var userForm, userAux models.User // Declaración del Modelo User

	ctx.BindJSON(&userForm)
	if userForm.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Campos faltantes"})
		return
	}

	db := config.DB

	db.Where("email = ?", userForm.Email).Find(&userAux)
	if userAux.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"Error:": "El correo ingresado no se encuentra registrado."})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userAux.Password), []byte(userForm.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"ERRORX:": err.Error()})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usr": userAux.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err_2 := token.SignedString([]byte(os.Getenv("SECRET")))
	if err_2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"ERROR:": err_2.Error()})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	ctx.JSON(http.StatusAccepted, gin.H{"Logueado": "Ha iniciado sesión correctamente."})

}

func LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("Authorization", "", -1, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"EXITO": "Se ha cerrado la sesion correctamente."})
}
