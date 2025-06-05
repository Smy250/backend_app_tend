package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Smy250/backend_app_tend/config"
	"github.com/Smy250/backend_app_tend/models"
	"github.com/Smy250/backend_app_tend/scripts"
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

	db, err1 := config.DB_Instance()
	if err1 != nil {
		panic("Error al conectar la base de datos")
	}

	//Verificamos si el correo existe.
	if exists := scripts.FindUserEmail(db, userForm.Email); exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error:": "El correo ya esta en uso."})
		return
	}

	// Encriptamos la contraseña que el usuario envio por el formulario JSON
	encriptedPassByte, err2 := bcrypt.GenerateFromPassword([]byte(userForm.Password), 14)
	if err2 != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"Error:": "Ha ocurrido un error al generar la contraseña del usuario"})
	}

	userForm.Password = string(encriptedPassByte)

	err3 := db.Create(&userForm).Error
	if err3 != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"Error:": "El correo ya existe."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"El usuario se ha registrado correctamente:": userForm})

}

func LoginUser(ctx *gin.Context) {
	var userForm, userAux models.User
	// Estas variables de modelo user, 1 sera usada para ser bindeadas
	// directametne del formulario y la 2 para bindear la info de la bd

	// Verificamos el correo enviado del JSON
	ctx.BindJSON(&userForm)
	if userForm.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Campos faltantes"})
		return
	}

	// Instanciamos la base de datos.
	db, err1 := config.DB_Instance()
	if err1 != nil {
		panic("Error al conectar la base de datos")
	}

	// Verificacion #1: Verificamos en la BD si el correo Existe.
	db.Where("email = ?", userForm.Email).First(&userAux)
	if userAux.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"Error:": "El correo ingresado no se encuentra registrado."})
		return
	}

	// Empleando la desencriptacion de tipo Bcrypt, verificamos si
	// coinciden o no las contraseñas enviada desde el formulario con
	// respecto a la almacenada en la base de datos.
	err2 := bcrypt.CompareHashAndPassword([]byte(userAux.Password), []byte(userForm.Password))
	if err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"ERRORX:": err2.Error()})
		return
	}

	// Crearemos un Token JWT con los parametros usr (ID del usuario), exp
	// (fecha de expiracion del mismo.)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usr": userAux.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Se firma y obtiene el token completamente codificado, de acuerdo a
	// una cadena que contiene informacion confidencial complete encoded token as a string using the secret
	tokenString, err3 := token.SignedString([]byte(os.Getenv("SECRET")))
	if err3 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"ERROR:": err3.Error()})
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
