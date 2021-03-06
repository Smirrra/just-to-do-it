package auth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func CreateAccessToken(userId int) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Minute * 20).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	JwtSecret, _ := os.LookupEnv("secret")
	JwtKey := []byte(JwtSecret)
	return token.SignedString(JwtKey)
}

func CreateRefreshToken(userId int) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	JwtSecret, _ := os.LookupEnv("secret")
	JwtKey := []byte(JwtSecret)
	return token.SignedString(JwtKey)
}

func ExtractAccessToken(r *http.Request) string {
	c, err := r.Cookie("access_token")
	if err != nil {
		return ""
	}
	token := c.Value
	return token
}

func ExtractRefreshToken(r *http.Request) string {
	c, err := r.Cookie("refresh_token")
	if err != nil {
		return ""
	}
	token := c.Value
	return token
}

func TokenValid(w http.ResponseWriter, r *http.Request, userId int) error {
	//check acc token
	tokenString := ExtractAccessToken(r)
	// if acc token has expired
	if tokenString == "" {
		err := RefTokenValid(r, userId)
		if err != nil {
			return err
		}
		err = RefreshTokens(w, r)
		if err != nil {
			return err
		}
		return nil
	}

	JwtSecret, _ := os.LookupEnv("secret")
	JwtKey := []byte(JwtSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v ", token.Header["alg"])
		}
		return JwtKey, nil
	})
	if err != nil {
		return err
	}
	claims := token.Claims.(jwt.MapClaims)
	userIdFromToken := int(claims["user_id"].(float64))
	if userId != userIdFromToken {
		return fmt.Errorf("Id don't match ")
	}
	return nil
}

func RefTokenValid(r *http.Request, userId int) error {
	//check refresh token
	tokenString := ExtractRefreshToken(r)
	JwtSecret, _ := os.LookupEnv("secret")
	JwtKey := []byte(JwtSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v ", token.Header["alg"])
		}
		return JwtKey, nil
	})
	if err != nil {
		return err
	}
	claims := token.Claims.(jwt.MapClaims)
	userIdFromToken := int(claims["user_id"].(float64))
	if userId != userIdFromToken {
		return fmt.Errorf("Id don't match ")
	}
	return nil
}

func RefreshTokens (w http.ResponseWriter, r *http.Request) error {
	refTokenStr := ExtractRefreshToken(r)
	JwtSecret, _ := os.LookupEnv("secret")
	JwtKey := []byte(JwtSecret)
	token, err := jwt.Parse(refTokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v ", token.Header["alg"])
		}
		return JwtKey, nil
	})
	if err != nil {
		return err
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))
	accToken, err := CreateAccessToken(userId)
	if err != nil {
		return err
	}
	SetCookieForAccToken(w, accToken)
	refTokenStr, err = CreateRefreshToken(userId)
	if err != nil {
		return err
	}
	SetCookieForRefToken(w, refTokenStr)
	return nil
}

func SetCookieForAccToken (w http.ResponseWriter, token string) {
	expirationTime := time.Now().Add(1 * time.Minute)
	http.SetCookie(w, &http.Cookie{
		Name: "access_token",
		Value: token,
		Expires: expirationTime,
		Path: "/",
	})
}

func SetCookieForRefToken (w http.ResponseWriter, token string) {
	expirationTime := time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: token,
		Expires: expirationTime,
		Path: "/",
	})
}

func CheckUser(w http.ResponseWriter, r *http.Request) (int, error) {
	paramFromURL := mux.Vars(r)
	userId, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		return 0, err
	}
	//???????????????? ?? ?? ???????????? ?????????????? ???????????? ????????????
	err = TokenValid(w, r, userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

var JwtCheck = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		notAuth := []string{"/register", "/login"} // List of endpoints for which authorization is not required
		requestPath := r.URL.Path // Current request path

		// We check if the authentication request does not require, we serve the request if it is not needed
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		//
		paramFromURL := mux.Vars(r)
		userId, err := strconv.Atoi(paramFromURL["id"])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Check and in case of timeout refresh token
		err = TokenValid(w, r, userId)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
