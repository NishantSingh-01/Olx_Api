package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// listing handeler
type Listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
type ListingInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	City        string  `json:"city"`
	Status      string  `json:"status"`
}

func ListingHandler(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query(
			`SELECT id , title , description , price , city , status , created_at  
		 FROM listings
		 ORDER BY created_at DESC
		 LIMIT 50
		 `)
		if err != nil {
			log.Printf("List : Query Error : %v", err)
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var listings []Listing
		for rows.Next() {
			var l Listing
			err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.Status, &l.CreatedAt)
			if err != nil {
				log.Printf("List : Scan Error : %v", err)
				http.Error(w, "Database query failed", http.StatusInternalServerError)
				return
			}
			listings = append(listings, l)
		}
		if err := rows.Err(); err != nil {
			log.Printf("List : Rows Error : %v", err)
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(listings)
	}

}
func DeleteListing(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		result, err := db.Exec(`
	   DELETE 
	   FROM listings
	   WHERE id = $1 
	   `, id)
		if err != nil {
			log.Printf("Delete : Exec Error : %v", err)
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Delete : RowsAffected Error : %v", err)
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]int64{"rows_affected": rowsAffected})
	}
}

func CreateListing(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input ListingInput
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Printf("Create : Decode Error : %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.City) == "" || strings.TrimSpace(input.Description) == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}
		if input.Price < 0 {
			http.Error(w, "Price cannot be negative", http.StatusBadRequest)
			return
		}
		//status logic there ? // todo

		var listing Listing
		rows, err := db.Query(`
		INSERT INTO listings
				(title, description, price, city, status)
			VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, ''), 'active'))
			RETURNING id, title, description, price, city, status, created_at
		`, input.Title, input.Description, input.Price, input.City, input.Status)
		if err != nil {
			log.Printf("Create : Query Error : %v", err)
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			err := rows.Scan(&listing.ID, &listing.Title, &listing.Description, &listing.Price, &listing.City, &listing.Status, &listing.CreatedAt)
			if err != nil {
				log.Printf("Create : Scan Error : %v", err)
				http.Error(w, "Database query failed", http.StatusInternalServerError)
				return
			}
		}
		if err := rows.Err(); err != nil {
			log.Printf("Create : Rows Error : %v", err)
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(listing)

	}
}
func UpdateListing(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if strings.TrimSpace(id) == "" {
			http.Error(w, "Listing ID is required", http.StatusBadRequest)
			return
		}
		var input ListingInput
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Printf("Update : Decode Error : %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.City) == "" || strings.TrimSpace(input.Description) == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}
		if input.Price < 0 {
			http.Error(w, "Price cannot be negative", http.StatusBadRequest)
			return
		}
		status := input.Status
		if strings.TrimSpace(status) == "" {
			status = "active"
		}
		var listing Listing
		err = db.QueryRow(`
		UPDATE listings
		SET title = $1,
			description = $2,
			price = $3,
			city = $4,
			status = $5
		WHERE id = $6
		RETURNING id, title, description, price, city, status, created_at
	`,
			strings.TrimSpace(input.Title),
			strings.TrimSpace(input.Description),
			input.Price,
			strings.TrimSpace(input.City),
			status,
			id,
		).Scan(
			&listing.ID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.City,
			&listing.Status,
			&listing.CreatedAt,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Listing not found", http.StatusNotFound)
				return
			}

			log.Printf("Update: Database Error: %v", err)
			http.Error(w, "Database update failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(listing)
	}
}
