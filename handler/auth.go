package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mr-Evgeny/go_final_project/config"
	"net/http"
	"time"
)

func createToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{"exp": time.Now().Add(time.Hour * 24).Unix()})

	tokenString, err := token.SignedString([]byte(config.SECRETKEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SECRETKEY), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(config.PASSWORD) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if err := verifyToken(cookie.Value); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func Sign(w http.ResponseWriter, r *http.Request) {
	if len(config.PASSWORD) == 0 {
		ErrorJson(w, "no need to auth")
		return
	}
	if r.Method != "POST" {
		http.Error(w, "wrong method", http.StatusBadRequest)
		return
	}
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		ErrorJson(w, err.Error())
		return
	}
	jsonPass := struct {
		V string `json:"password"`
	}{}

	if err = json.Unmarshal(buf.Bytes(), &jsonPass); err != nil {
		ErrorJson(w, err.Error())
		return
	}
	if len(jsonPass.V) == 0 {
		ErrorJson(w, "No pass specified")
		return
	}
	if jsonPass.V != config.PASSWORD {
		ErrorJson(w, "Wrong pass specified")
		return
	}
	token, err := createToken()
	if err != nil {
		ErrorJson(w, err.Error())
		return
	}

	json_resp, err := json.Marshal(struct {
		T string `json:"token"`
	}{token})
	if err != nil {
		ErrorJson(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(json_resp)
}
