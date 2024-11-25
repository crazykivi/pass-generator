package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
)

type PasswordRequest struct {
	Length      int  `json:"length"`
	WithSymbols bool `json:"with_symbols"`
	WithDigits  bool `json:"with_digits"`
	Count       int  `json:"count"`
}

// Функция для генерации пароля
func generatePassword(length int, withSymbols bool, withDigits bool) string {
	var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var digits = "0123456789"
	var symbols = "!@#$%^&*()-_=+[]{}|;:,.<>?/`~"
	var charSet string

	if withSymbols {
		charSet += symbols
	}
	if withDigits {
		charSet += digits
	}

	if charSet == "" {
		charSet = letters
	}

	var password strings.Builder
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			log.Fatal("Ошибка при генерации случайного числа: ", err)
		}
		password.WriteByte(charSet[index.Int64()])
	}

	return password.String()
}

// Функция обработчик POST запроса для генерации паролей
func generatePasswords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request PasswordRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.Length < 6 {
		http.Error(w, "Пароль должен быть не короче 6 символов", http.StatusBadRequest)
		return
	}

	var passwords []string
	for i := 0; i < request.Count; i++ {
		password := generatePassword(request.Length, request.WithSymbols, request.WithDigits)
		passwords = append(passwords, password)
	}

	// Ответ в формате JSON
	response := map[string]interface{}{
		"passwords": passwords,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/generate-passwords", generatePasswords)

	port := ":8080"
	fmt.Println("Сервер запущен на порту ", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
