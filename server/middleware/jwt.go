package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwtgo "github.com/auth0/go-jwt-middleware/validate/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const signatureAlgorithm = "RS256"

var _ jwtgo.CustomClaims = &CustomClaims{}

type User struct {
	Id string `json:"user_id"`
}

type CustomClaims struct {
	Scope string `json:"scope"`
	User  User   `json:"https://tourex.wie.gg/user"`
	jwt.StandardClaims
}

func (c CustomClaims) Validate(_ context.Context) error {
	expectedAudience := os.Getenv("AUTH0_AUDIENCE")
	if c.Audience != expectedAudience {
		return fmt.Errorf("token claim validation failed: unexpected audience %q", c.Audience)
	}

	expectedIssuer := "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
	if c.Issuer != expectedIssuer {
		return fmt.Errorf("token claim validation failed: unexpected issuer %q", c.Issuer)
	}

	return nil
}

func EnsureValidToken() gin.HandlerFunc {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		certificate, err := getPEMCertificate(token)
		if err != nil {
			return token, err
		}

		return jwt.ParseRSAPublicKeyFromPEM([]byte(certificate))
	}

	customClaims := func() jwtgo.CustomClaims {
		return &CustomClaims{}
	}

	validator, err := jwtgo.New(keyFunc, signatureAlgorithm, jwtgo.WithCustomClaims(customClaims))
	if err != nil {
		log.Fatal("Failed to set up JWT validator")
	}

	m := jwtmiddleware.New(validator.ValidateToken)

	return func(c *gin.Context) {
		var encounteredError = true
		var handler http.HandlerFunc = func(rw http.ResponseWriter, r *http.Request) {
			encounteredError = false
			c.Request = r
			c.Next()
		}

		m.CheckJWT(handler).ServeHTTP(c.Writer, c.Request)

		if encounteredError {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "Failed to validate JWT"})
		}
	}
}

type (
	jwks struct {
		Keys []jsonWebKeys `json:"keys"`
	}

	jsonWebKeys struct {
		Kty string   `json:"kty"`
		Kid string   `json:"kid"`
		Use string   `json:"use"`
		N   string   `json:"n"`
		E   string   `json:"e"`
		X5c []string `json:"x5c"`
	}
)

func getPEMCertificate(token *jwt.Token) (string, error) {
	response, err := http.Get("https://" + os.Getenv("AUTH0_DOMAIN") + "/.well-known/jwks.json")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var jwks jwks
	if err = json.NewDecoder(response.Body).Decode(&jwks); err != nil {
		return "", err
	}

	var cert string
	for _, key := range jwks.Keys {
		if token.Header["kid"] == key.Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + key.X5c[0] + "\n-----END CERTIFICATE-----"
			break
		}
	}

	if cert == "" {
		return cert, errors.New("unable to find appropriate key")
	}

	return cert, nil
}

func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}
