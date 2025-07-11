package scripts

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Este breve script se usa en la lógica del middleware.
// Basicamente recibe por parametro un string del token de la cookie
// con la información cifrada, que luego sera procesado y devuelto
// como tipo de dato jwt.Token en conjunto con un error.
func DecryptToken(tokenString string) (*jwt.Token, error) {
	tokenProcessed, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	return tokenProcessed, err
}
