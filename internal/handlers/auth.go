package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
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

//repository
func checkUserExists(db *sql.DB, email string) error {
	var id string
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
	if err == nil {
		return errors.New("user already exists")
	}
	return nil
}

//handler functions
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
	if len(req.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}
	// ctx := r.Context()
	// userExist := ah.db.QueryRowContext(ctx, `
	//   SELECT id , email , password 
	//   FROM users 
	//   WHERE email = $1
	// `, req.Email)
	// if userExist == nil {
	// 	http.Error(w, "User already exists", http.StatusBadRequest)
	// 	return
	// }
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
}

func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

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
