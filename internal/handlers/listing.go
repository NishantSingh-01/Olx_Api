package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"olx_api/internal/httpx"
	"olx_api/internal/middleware"
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

type Listing_Handler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListing_Handler(db *sql.DB, logger *slog.Logger) *Listing_Handler {
	return &Listing_Handler{
		db:     db,
		logger: logger,
	}
}

func (lh *Listing_Handler) ListingHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancle := context.WithTimeout(r.Context(), 3*time.Second)//this is by gpt to stop cancelled query ohter method is direct to the cts pass there 
	defer cancle()
	rows, err := lh.db.QueryContext(
		ctx,
		`SELECT id , title , description , price , city , status , created_at  
		 FROM listings
		 ORDER BY created_at DESC
		 LIMIT 50
		 `)
	if err != nil {
		// log.Printf("List : Query Error : %v", err)
		lh.logger.Error("Listing Query Error", "listing_handler", err)
		httpx.Error(w, "Database query failed", "test", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var listings []Listing
	for rows.Next() {
		var l Listing
		err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.Status, &l.CreatedAt)
		if err != nil {
			lh.logger.Error("Listing Scan Error", "listing_handler", err)
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}
		listings = append(listings, l)
	}
	if err := rows.Err(); err != nil {
		lh.logger.Error("Listing Rows Error", "listing_handler", err)
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	lh.logger.Info("ListingHandler Success", "listing_handler", len(listings))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(listings)
}
func (lh *Listing_Handler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx := r.Context()
	requestId := middleware.GetRequestID(ctx)
	// slog.Debug("Delete Debug : id","listing_delete", id) //1
	// slog.Info("Delete Info : id","listing_delete", id) //2
	// slog.Warn("Delete warn: id","listing_delete", id) //3

	result, err := lh.db.ExecContext(ctx, `
	   DELETE 
	   FROM listings
	   WHERE id = $1 
	   `, id)
	if err != nil {
		lh.logger.Error("Delete : Exec Error", "listing_delete", id,"requestId",requestId ,"Error", err)
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		lh.logger.Error("Delete : RowsAffected Error", "listing_delete", id,"requestId",requestId , "Error", err)
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	lh.logger.Info("Delete Exec Success : RowsAffected :", "listing_delete", id, "rows_affected", rowsAffected)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]int64{"rows_affected": rowsAffected})
}















// listing handler function
func (lh *Listing_Handler) CreateListing(w http.ResponseWriter, r *http.Request) {
	var input ListingInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		// log.Printf("Create : Decode Error : %v", err)
		lh.logger.Error("Create : Invalid request body", "listing_create", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.City) == "" || strings.TrimSpace(input.Description) == "" {
		lh.logger.Error("Create : All fields are required", "listing_create", err)
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}
	if input.Price < 0 {
		http.Error(w, "Price cannot be negative", http.StatusBadRequest)
		return
	}
	//status logic there ? // todo

	var listing Listing
	ctx := r.Context()
	rows, err := lh.db.QueryContext(ctx, `
		INSERT INTO listings
				(title, description, price, city, status)
			VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, ''), 'active'))
			RETURNING id, title, description, price, city, status, created_at
		`, input.Title, input.Description, input.Price, input.City, input.Status)
	if err != nil {
		// log.Printf("Create : Query Error : %v", err)
		lh.logger.Error("Create : Database query failed", "listing_create", err)
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		err := rows.Scan(&listing.ID, &listing.Title, &listing.Description, &listing.Price, &listing.City, &listing.Status, &listing.CreatedAt)
		if err != nil {
			// log.Printf("Create : Scan Error : %v", err)
			lh.logger.Error("Create : Database query failed", "listing_create", err)
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}
	}
	if err := rows.Err(); err != nil {
		// log.Printf("Create : Rows Error : %v", err)
		lh.logger.Error("Create : Database query failed", "listing_create", err)
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(listing)

}
func (lh *Listing_Handler) GetListingByID(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			lh.logger.Error("GetListingByID : Invalid request body","listing_getbyid",nil)
			http.Error(w, "Listing ID is required", http.StatusBadRequest)
			return
		}
		var listing Listing
		ctx := r.Context()
		err := lh.db.QueryRowContext(ctx,`
		SELECT id , title , description , price , city , status , created_at 
		FROM listings
		WHERE id = $1
		`, id).Scan(
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
				lh.logger.Error("GetListingByID : Listing not found","",err)
				http.Error(w, "Listing not found", http.StatusNotFound)
				return
			}
			lh.logger.Error("GetListingByID : Database query failed","listing_getbyid",err)
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}
		lh.logger.Info("GetListingByID : Success","listing_getbyid","listing",listing)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]Listing{"listing found successfully": listing})
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

