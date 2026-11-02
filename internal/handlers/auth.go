package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"olx_api/internal/config"
	"olx_api/internal/httpx"
	"olx_api/internal/utils"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
	Token        string `json:"token"`
}
type AuthResponse struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
}
type ProfileResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
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

	// var id,email,name string
	type USER struct{
		id string
		email string
		name string
	}
	var user USER
	err = ah.db.QueryRowContext(r.Context(), `
		INSERT INTO users (email, name, hashed_password)
		VALUES ($1, $2, $3)
		RETURNING id,email,name
	`, req.Email, req.Name, hashedPassword).Scan(&user.id,&user.email,&user.name)
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
	
	cfg := config.MustLoad()
	token,err := utils.GenerateToken(cfg.JWTSecret,user.id,time.Hour*24*7)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, map[string]any{
		"token":   token,
		"email":   user.email,
		"name":    user.name,
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
	var input RefreshInput
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&input)
	}

	refreshTokenString := strings.TrimSpace(input.RefreshToken)
	if refreshTokenString == "" {
		refreshTokenString = strings.TrimSpace(input.Token)
	}
	if refreshTokenString == "" {
		authHeader := r.Header.Get("Authorization")
		parts := strings.Fields(authHeader)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			refreshTokenString = parts[1]
		}
	}

	if refreshTokenString == "" {
		httpx.Error(w, "Refresh token is required", httpx.Errvalidation_failed, http.StatusBadRequest)
		return
	}

	cfg := config.MustLoad()
	token, err := utils.VerifyToken(refreshTokenString, cfg.JWTSecret)
	if err != nil || token == nil || !token.Valid {
		ah.logger.Error("Refresh : invalid or expired token", "refresh", err)
		httpx.Error(w, "Invalid or expired refresh token", httpx.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims == nil {
		httpx.Error(w, "Invalid token claims", httpx.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var userID string
	switch v := claims["user_id"].(type) {
	case string:
		userID = v
	case float64:
		userID = fmt.Sprintf("%.0f", v)
	}

	if strings.TrimSpace(userID) == "" {
		httpx.Error(w, "Invalid user in token claims", httpx.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var userExists bool
	err = ah.db.QueryRowContext(r.Context(), "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&userExists)
	if err != nil || !userExists {
		ah.logger.Error("Refresh : user not found or database error", "error", err, "user_id", userID)
		httpx.Error(w, "User not found", httpx.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	newToken, err := utils.GenerateToken(cfg.JWTSecret, userID, time.Hour*24*7)
	if err != nil {
		ah.logger.Error("Refresh : failed to generate token", "error", err)
		httpx.Error(w, "Failed to generate token", httpx.Errinternal_error, http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := utils.GenerateToken(cfg.JWTSecret, userID, time.Hour*24*30)
	if err != nil {
		ah.logger.Error("Refresh : failed to generate refresh token", "error", err)
		httpx.Error(w, "Failed to generate refresh token", httpx.Errinternal_error, http.StatusInternalServerError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"token":         newToken,
		"refresh_token": newRefreshToken,
		"message":       "Token refreshed successfully",
	})
}
func (ah *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear any auth cookies that may have been set
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	ah.logger.Info("LogoutHandler : user logged out successfully")
	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
func (ah *AuthHandler) GetMyProfileHandler(w http.ResponseWriter, r *http.Request) {
	var userID string
	if v := r.Context().Value("user_id"); v != nil {
		switch id := v.(type) {
		case string:
			userID = id
		case uint:
			userID = fmt.Sprintf("%d", id)
		case int:
			userID = fmt.Sprintf("%d", id)
		case float64:
			userID = fmt.Sprintf("%.0f", id)
		}
	}

	if userID == "" {
		authHeader := r.Header.Get("Authorization")
		parts := strings.Fields(authHeader)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			cfg := config.MustLoad()
			token, err := utils.VerifyToken(parts[1], cfg.JWTSecret)
			if err == nil && token != nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && claims != nil {
					switch id := claims["user_id"].(type) {
					case string:
						userID = id
					case float64:
						userID = fmt.Sprintf("%.0f", id)
					}
				}
			}
		}
	}

	if strings.TrimSpace(userID) == "" {
		ah.logger.Warn("GetMyProfile : unauthorized access attempt")
		httpx.Error(w, "Unauthorized: authentication required", httpx.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var profile ProfileResponse
	err := ah.db.QueryRowContext(r.Context(), `
		SELECT id, name, email
		FROM users
		WHERE id = $1
	`, userID).Scan(&profile.ID, &profile.Name, &profile.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ah.logger.Warn("GetMyProfile : user not found", "user_id", userID)
			httpx.Error(w, "User not found", httpx.Errnot_found, http.StatusNotFound)
			return
		}
		ah.logger.Error("GetMyProfile : database query failed", "error", err, "user_id", userID)
		httpx.Error(w, "Database query failed", httpx.Errinternal_error, http.StatusInternalServerError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, profile)
}
