package domain

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

// getJWTKey retorna la clave secreta para la firma de JWT de forma perezosa (lazy loading),
// leyendola de variables de entorno o recurriendo a un valor por defecto seguro.
func getJWTKey() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("la variable de entorno JWT_SECRET no esta configurada")
	}
	return []byte(secret), nil
}

// GenerateToken crea un nuevo token firmado digitalmente con JWT que contiene el username
// del usuario autenticado y una duracion de validez de 24 horas.
func GenerateToken(username string) (string, error) {
	key, err := getJWTKey()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken analiza, verifica y valida la integridad de un token JWT recibido,
// devolviendo el nombre de usuario si el token es valido y no ha expirado.
func ValidateToken(tokenString string) (string, error) {
	key, err := getJWTKey()
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metodo de firma inesperado")
		}
		return key, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("claims de token invalidos")
		}
		return username, nil
	}

	return "", errors.New("token invalido")
}
