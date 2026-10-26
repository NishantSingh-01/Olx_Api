package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"olx_api/internal/httpx"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewAuthHandler(db *sql.DB, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		db:     db,
		logger: logger,
	}
}

// user models
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthResponse struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
}

// utils
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// repository
func checkUserExists(db *sql.DB, email string) error {
	var id string
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
	if err == nil {
		return errors.New("user already exists")
	}
	return nil
}

// response fun //>> responseHandler >> we can do it by httpx package 
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// handler functions
func (ah *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ah.logger.Error("Register : Invalid request body", "register", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 || len(req.Password) > 72 {
		http.Error(w, "Password must be between 8 and 72 characters", http.StatusBadRequest)
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = ah.db.ExecContext(r.Context(), `
		INSERT INTO users (email, name, hashed_password)
		VALUES ($1, $2, $3)
	`, req.Email, req.Name, hashedPassword)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		ah.logger.Error("register: failed to create user", "error", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"message": "User created successfully",
	// })
	httpx.WriteJSON(w, http.StatusCreated, map[string]any{
		"email":   req.Email,
		"name":    req.Name,
		"message": "user created successfully",
	})
}

func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		ah.logger.Error("Login : Invalid request body", "login", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if input.Email == "" || input.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if len(input.Password) < 8 || len(input.Password) > 72 {
		http.Error(w, "Password must be between 8 and 72 characters", http.StatusBadRequest)
		return
	}
	// var u users
	// row := ah.db.QueryRowContext(r.Context(), `
	// 	SELECT id,name,hashed_password FROM users WHERE email = $1
	// 	`, input.Email)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
	// 		return
	// 	}
	// 	ah.logger.Error("login: failed to find user", "error", err)
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}

func (ah *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}
func (ah *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}
func (ah *AuthHandler) GetMyProfileHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}
