package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(id uint32) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	// Create a new token with the HS256 signing method and the claims map
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the API_SECRET environment variable and return the signed token as a string
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func TokenValid(r *http.Request) error {
	// extract token from request
	tokenString := ExtractToken(r)

	// parse token and verify signature with HMAC method
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	// if token is valid, extract and print the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}
	// return nil if no errors occurred
	return nil
}

func ExtractToken(r *http.Request) string {
	// Get token from query parameters
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	// Get token from Authorization header
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	// If token is not found in query or header, return empty string
	return ""
}

// ExtractTokenID extracts the ID from the JWT token in the request header or URL query parameter
func ExtractTokenID(r *http.Request) (uint32, error) {
	// Extract the token string from the request
	tokenString := ExtractToken(r)
	// Parse the token using the API_SECRET as the key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Return the API_SECRET as the key
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		// Return any errors encountered during parsing
		return 0, err
	}
	// Extract the claims from the token and convert the "id" claim to a uint32
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		// Return the ID as a uint32
		return uint32(uid), nil
	}
	// Return 0 and nil if the token is invalid or does not contain an ID claim
	return 0, nil
}

// Pretty display the claims licely in the terminal
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}
